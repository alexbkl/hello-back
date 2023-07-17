package handlers

import (
	"bytes"
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"regexp"
	"sync"
	"time"

	"github.com/Hello-Storage/hello-back/config"
	"github.com/Hello-Storage/hello-back/entities"
	"github.com/Hello-Storage/hello-back/s3client"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/ipfs/go-cid"
	mc "github.com/multiformats/go-multicodec"
	mh "github.com/multiformats/go-multihash"
)

var (
	ErrUserNotExists  = errors.New("user does not exist")
	ErrUserExists     = errors.New("user already exists")
	ErrInvalidAddress = errors.New("invalid address")
	ErrInvalidNonce   = errors.New("invalid nonce")
	ErrMissingSig     = errors.New("signature is missing")
	ErrAuthError      = errors.New("authentication error")
)

var jwtProvider = NewJwtHmacProvider(
	"env secret",
	"hello-storage",
	time.Minute*60,
)

/*
	func printFileSize(weight int) {
		const (
			KB = 1024
			MB = KB * 1024
			GB = MB * 1024
			TB = GB * 1024
		)

		var size string
		switch {
		case weight >= TB:
			size = fmt.Sprintf("%.2fTB", float64(weight)/TB)
		case weight >= GB:
			size = fmt.Sprintf("%.2fGB", float64(weight)/GB)
		case weight >= MB:
			size = fmt.Sprintf("%.2fMB", float64(weight)/MB)
		case weight >= KB:
			size = fmt.Sprintf("%.2fKB", float64(weight)/KB)
		default:
			size = fmt.Sprintf("%dB", weight)
		}

		fmt.Println("File Size:", size)
	}
*/
type JwtHmacProvider struct {
	hmacSecret []byte
	issuer     string
	duration   time.Duration
}

func NewJwtHmacProvider(hmacSecret string, issuer string, duration time.Duration) *JwtHmacProvider {
	ans := JwtHmacProvider{
		hmacSecret: []byte(hmacSecret),
		issuer:     issuer,
		duration:   duration,
	}
	return &ans
}

func (j *JwtHmacProvider) CreateStandard(subject string) (string, error) {
	now := time.Now()
	claims := jwt.RegisteredClaims{
		Issuer:    j.issuer,
		Subject:   subject,
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(j.duration)),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.hmacSecret)
}

func (j *JwtHmacProvider) Verify(tokenString string) (*jwt.RegisteredClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return j.hmacSecret, nil
	})
	if err != nil {
		return nil, ErrAuthError
	}
	if claims, ok := token.Claims.(*jwt.RegisteredClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, ErrAuthError
}

type MemStorage struct {
	lock  sync.RWMutex
	users map[string]entities.User
}

func (m *MemStorage) CreateIfNotExists(u entities.User) error {
	m.lock.Lock()
	defer m.lock.Unlock()
	if _, exists := m.users[u.Address]; exists {
		return ErrUserExists
	}
	m.users[u.Address] = u
	return nil
}

func (m *MemStorage) Get(address string) (entities.User, error) {
	m.lock.RLock()
	defer m.lock.RUnlock()
	u, exists := m.users[address]
	if !exists {
		return u, ErrUserNotExists
	}
	return u, nil
}

func (m *MemStorage) Update(user entities.User) error {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.users[user.Address] = user
	return nil
}

func NewMemStorage() *MemStorage {
	ans := MemStorage{
		users: make(map[string]entities.User),
	}
	return &ans
}

// ============================================================================

var (
	hexRegex   *regexp.Regexp = regexp.MustCompile(`^0x[a-fA-F0-9]{40}$`)
	nonceRegex *regexp.Regexp = regexp.MustCompile(`^[0-9]+$`)
)

type RegisterPayload struct {
	Address string `json:"address"`
}

func (p RegisterPayload) Validate() error {
	if !hexRegex.MatchString(p.Address) {
		return ErrInvalidAddress
	}
	return nil
}

func RegisterHandler(c *fiber.Ctx) error {
	var p RegisterPayload

	if err := c.BodyParser(&p); err != nil {
		return c.Status(503).SendString(err.Error())
	}

	if err := p.Validate(); err != nil {
		return c.Status(503).SendString(err.Error())
	}

	// check if user already exists
	var user entities.User

	//where user.address = address
	result := config.Database.Find(&user, "address = ?", p.Address)

	if result.RowsAffected != 0 {
		return c.Status(203).SendString("User already exists")
	}

	nonce, err := GetNonce()
	if err != nil {
		return c.Status(503).SendString(err.Error())
	}

	u := entities.User{
		Address: p.Address,
		Nonce:   nonce,
		DataCap: 100,
	}

	config.Database.Create(&u)

	return c.Status(201).JSON(u)
}

func UserNonceHandler(c *fiber.Ctx) error {

	/*
		return func(w http.ResponseWriter, r *http.Request) {
			address := chi.URLParam(r, "address")
			if !hexRegex.MatchString(address) {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			user, err := storage.Get(strings.ToLower(address))
			if err != nil {
				switch errors.Is(err, ErrUserNotExists) {
				case true:
					w.WriteHeader(http.StatusNotFound)
				default:
					w.WriteHeader(http.StatusInternalServerError)
				}
				return
			}
			resp := struct {
				Nonce string
			}{
				Nonce: user.Nonce,
			}
			renderJson(r, w, http.StatusOK, resp)
		}
	*/

	//Refactored to use Fiber:

	address := c.Params("address")

	if !hexRegex.MatchString(address) {
		return c.Status(400).SendString("Invalid address")
	}

	var user entities.User

	//where user.address = address
	result := config.Database.Find(&user, "address = ?", address)

	if result.RowsAffected == 0 {
		return c.Status(404).SendString("Address not found")
	}

	resp := struct {
		Nonce string `json:"nonce"`
	}{
		Nonce: user.Nonce,
	}

	return c.Status(200).JSON(resp)
}

type SigninPayload struct {
	Address string `json:"address"`
	Nonce   string `json:"nonce"`
	Sig     string `json:"sig"`
}

func (s SigninPayload) Validate() error {
	if !hexRegex.MatchString(s.Address) {
		return ErrInvalidAddress
	}
	if !nonceRegex.MatchString(s.Nonce) {
		return ErrInvalidNonce
	}
	if len(s.Sig) == 0 {
		return ErrMissingSig
	}
	return nil
}

func SigninHandler(c *fiber.Ctx) error {

	var p SigninPayload

	if err := c.BodyParser(&p); err != nil {
		return c.Status(503).SendString(err.Error())
	}

	if err := p.Validate(); err != nil {
		return c.Status(503).SendString(err.Error())
	}

	address := p.Address

	user, err := Authenticate(address, p.Nonce, p.Sig)

	switch err {
	case nil:
	case ErrAuthError:
		return c.Status(401).SendString("Authentication error")
	default:
		return c.Status(500).SendString("Internal server error")
	}

	signedToken, err := jwtProvider.CreateStandard(user.Address)
	if err != nil {
		return c.Status(500).SendString("Internal server error")
	}

	resp := struct {
		AccessToken string `json:"access"`
	}{
		AccessToken: signedToken,
	}

	return c.Status(200).JSON(resp)
}

func WelcomeHandler(c *fiber.Ctx) error {
	var user entities.User

	//get user from context
	user = c.Locals("user").(entities.User)

	resp := struct {
		Msg string `json:"msg"`
	}{
		Msg: "Welcome, address: " + user.Address + "\n",
	}

	return c.Status(200).JSON(resp)
}

//upload handler s3 client

func UploadHandler(c *fiber.Ctx) error {
	//get user and save it to database
	var user entities.User

	//get user from context
	user = c.Locals("user").(entities.User)

	//get encryptedFileBlob from request
	file, err := c.FormFile("encryptedFileBlob")

	//print formatted the file name

	//get encryptedMetadataStr from request
	encryptedMetadataStr := c.FormValue("encryptedMetadataStr")

	//get encryptedCidStr
	cidOfEncryptedBufferStr := c.FormValue("cidOfEncryptedBufferStr")

	//get cidEncryptedOriginalStr
	cidEncryptedOriginalStr := c.FormValue("cidEncryptedOriginalStr")

	//get ivString
	ivString := c.FormValue("ivString")

	//open file
	src, err := file.Open()
	if err != nil {

		return c.Status(500).SendString("Internal server error: " + err.Error())
	}
	defer src.Close()

	//create a new file in s3 bucket

	srcBytes, error := s3client.UploadFile(cidOfEncryptedBufferStr, src)

	//get the weight of the file
	weight := len(srcBytes)

	//print the weight and process is accordingly (TB, GB, MB, KB, B,...)

	//printFileSize(weight)

	if error != nil {
		//if 		error is fmt.Errorf("File already exists"), return 409
		if error.Error() == "File already exists" {
			fmt.Println(error, "UploadFile error: File already exists")

			//return c.Status(409).SendString("File already exists")
		}
		fmt.Println(error, "Error UploadFile")

		return c.Status(500).SendString("Internal server error: " + error.Error())
	}
	// Create a cid manually by specifying the 'prefix' parameters
	pref := cid.Prefix{
		Version:  1,
		Codec:    uint64(mc.Raw),
		MhType:   mh.SHA2_256,
		MhLength: -1, // default length
	}

	// And then feed it some data
	cid, err := pref.Sum(srcBytes)
	if err != nil {
		fmt.Println("Error creating CID: ", err)
		return err
	}

	//print the cid and the cidOfEncryptedBufferStr
	//fmt.Println("CID: ", cid)
	//fmt.Println("cidOfEncryptedBufferStr: ", cidOfEncryptedBufferStr)

	//compare cid and cidOfEncryptedBufferStr, if not equal, return 400
	if cid.String() != cidOfEncryptedBufferStr {
		fmt.Println("cid and cidOfEncryptedBufferStr are not equal")
		return c.Status(400).SendString("Backend CID and Frontend CID of encrypted buffer are not equal")
	}

	// close file
	src.Close()
	if error != nil {
		fmt.Println("Error creating CID: ", err)
		return err
	}
	//save file to database
	fileToUpload := entities.File{
		EncryptedMetadata:       encryptedMetadataStr,
		UserAddress:             user.Address,
		CIDOfEncryptedBuffer:    cid.String(),
		CIDEncryptedOriginalStr: cidEncryptedOriginalStr,
		IV:                      ivString,
		BytesLength:             weight,
	}

	config.Database.Create(&fileToUpload)

	//add weight of the file to user's total UsedStorage (be aware that UsedStorage is an int64)
	user.UsedStorage += int64(weight)
	user.TotalUploadedFiles += 1
	config.Database.Save(&user)

	resp := struct {
		File entities.File `json:"file"`
	}{
		File: fileToUpload,
	}
	return c.Status(200).JSON(resp)

}

func DeleteFileHandler(c *fiber.Ctx) error {
	//get user
	var user entities.User

	//get user from context
	user = c.Locals("user").(entities.User)

	//get address from user
	address := user.Address
	fileId := c.Params("fileId")
	fmt.Println(fileId)

	//check if file from entities.File belongs to user
	var file entities.File

	config.Database.Where("id = ? AND user_address = ?", fileId, address).First(&file)

	if file.ID == 0 {
		return c.Status(403).SendString("File not found")
	}

	//get weight of the file
	weight := file.BytesLength

	//subtract weight of the file from user's total UsedStorage (be aware that UsedStorage is an int64)
	user.UsedStorage -= int64(weight)

	//delete file from database
	config.Database.Delete(&file)

	//check if there are other files with the same cid
	var otherFiles []entities.File
	config.Database.Where("c_id_of_encrypted_buffer = ?", file.CIDOfEncryptedBuffer).Find(&otherFiles)

	//delete file from s3 bucket only if no other files with the same CID exist in the database
	if len(otherFiles) == 0 {
		err := s3client.DeleteFile(file.CIDOfEncryptedBuffer)

		if err != nil {
			return c.Status(500).SendString("Internal server error: " + err.Error())
		}
	}

	resp := struct {
		Msg string `json:"msg"`
	}{
		Msg: "File deleted successfully",
	}

	return c.Status(200).JSON(resp)

}

func DownloadFileHandler(c *fiber.Ctx) error {
	//get user
	var user entities.User

	//get user from context
	user = c.Locals("user").(entities.User)

	//get address from user
	address := user.Address

	//get cid from params
	cid := c.Params("cid")

	//check if file from entities.File belongs to user
	var file entities.File

	config.Database.Where("c_id_of_encrypted_buffer = ? AND user_address = ?", cid, address).First(&file)

	if file.ID == 0 {
		return c.Status(403).SendString("File not found")
	}

	//download file from s3 bucket
	//result is *s3.GetObjectOutput type
	result, err := s3client.DownloadFile(cid)

	if err != nil {
		return c.Status(500).SendString("Internal server error: " + err.Error())
	}

	// Create a buffer to hold the file contents
	buf := new(bytes.Buffer)
	buf.ReadFrom(result.Body)

	// You can also add a Content-Disposition header to the response to suggest a filename to the client
	filename, ok := result.Metadata["filename"]
	if !ok {
		filename = &file.CIDOfEncryptedBuffer
	}

	c.Response().Header.Set("Access-Control-Expose-Headers", "Original-Filename")

	c.Response().Header.Set("Content-Disposition", "attachment; filename="+*filename)
	//add all the metadata to the response
	for key, value := range result.Metadata {
		fmt.Println(key, *value)
		c.Response().Header.Set(key, *value)
	}

	// Write the file contents into the response body
	return c.Status(200).Send(buf.Bytes())

}

func GetFilesHandler(c *fiber.Ctx) error {
	//get user
	var user entities.User

	//get user from context
	user = c.Locals("user").(entities.User)

	//get files from database
	var files []entities.File

	config.Database.Find(&files, "user_address = ?", user.Address)

	resp := struct {
		Files []entities.File `json:"files"`
	}{
		Files: files,
	}

	return c.Status(200).JSON(resp)
}

// ============================================================================

func getUserFromReqContext(r *http.Request) entities.User {
	ctx := r.Context()
	key := ctx.Value("user").(entities.User)
	return key
}

func AuthMiddleware(c *fiber.Ctx) error {

	//get header value authorization
	headerValue := c.Get("Authorization")
	const prefix = "Bearer "
	if len(headerValue) < len(prefix) {
		return c.Status(401).SendString("Unauthorized")
	}

	tokenString := headerValue[len(prefix):]
	if len(tokenString) == 0 {
		return c.Status(401).SendString("Unauthorized")
	}

	claims, err := jwtProvider.Verify(tokenString)
	if err != nil {
		return c.Status(401).SendString("Unauthorized")
	}
	//claims.Subject is the address of the user
	var user entities.User

	//where user.address = address
	result := config.Database.Find(&user, "address = ?", claims.Subject)

	if result.RowsAffected == 0 {
		return c.Status(401).SendString("Unauthorized")
	}

	//set the user in the context
	c.Locals("user", user)

	return c.Next()

}

func Authenticate(address string, nonce string, sigHex string) (entities.User, error) {
	var user entities.User
	result := config.Database.Find(&user, "address = ?", address)

	if result.RowsAffected == 0 {
		return user, ErrUserNotExists
	}

	if user.Nonce != nonce {
		return user, ErrAuthError
	}

	sig := hexutil.MustDecode(sigHex)
	// https://github.com/ethereum/go-ethereum/blob/master/internal/ethapi/api.go#L516
	// check here why I am subtracting 27 from the last byte
	sig[crypto.RecoveryIDOffset] -= 27
	msg := accounts.TextHash([]byte(nonce))
	recovered, err := crypto.SigToPub(msg, sig)
	if err != nil {
		return user, err
	}
	recoveredAddr := crypto.PubkeyToAddress(*recovered)

	if user.Address != recoveredAddr.Hex() {
		return user, ErrAuthError
	}

	// update the nonce here so that the signature cannot be reSUSed
	nonce, err = GetNonce()
	if err != nil {
		return user, err
	}
	user.Nonce = nonce
	result = config.Database.Where("address = ?", address).Updates(&user)

	if result.RowsAffected == 0 {
		return user, ErrAuthError
	}

	return user, nil
}

var (
	max  *big.Int
	once sync.Once
)

func GetNonce() (string, error) {
	once.Do(func() {
		max = new(big.Int)
		max.Exp(big.NewInt(2), big.NewInt(130), nil).Sub(max, big.NewInt(1))
	})
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		return "", err
	}
	return n.Text(10), nil
}

func bindReqBody(r *http.Request, obj any) error {
	return json.NewDecoder(r.Body).Decode(obj)
}

func renderJson(r *http.Request, w http.ResponseWriter, statusCode int, res interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8 ")
	var body []byte
	if res != nil {
		var err error
		body, err = json.Marshal(res)
		if err != nil { // TODO handle me better
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
	w.WriteHeader(statusCode)
	if len(body) > 0 {
		w.Write(body)
	}
}

// ============================================================================

func Run() error {
	// initialization of storage

	// setup the endpoints
	r := chi.NewRouter()

	//  Just allow all for the reference implementation
	r.Use(cors.AllowAll().Handler)

	/*
		r.Group(func(r chi.Router) {
			r.Use(AuthMiddleware(jwtProvider))
			r.Get("/welcome", WelcomeHandler())
		})
	*/
	// start the server on port 8001
	err := http.ListenAndServe("185.166.212.43:8001", r)
	return err
}
