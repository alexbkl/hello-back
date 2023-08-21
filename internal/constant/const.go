package constant

import (
	"fmt"
)

const (
	AuthorizationHeaderKey  = "authorization"
	AuthorizationTypeBearer = "bearer"
	AuthorizationPayloadKey = "authorization_payload"
)

const LoginMessage = "Greetings from joinhello\nSign this message to log into joinhello\nnonce: "

func BuildLoginMessage(nonce string) []byte {
	return []byte(fmt.Sprintf("%s%s", LoginMessage, nonce))
}
