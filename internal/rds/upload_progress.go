package rds

import (
	"fmt"
	"strings"

	"github.com/redis/rueidis"
)

type UploadProgressValue struct {
	Name string `json:"name"`
	Size int64  `json:"size"`
	Read int64  `json:"read"`
}

func (g *RdsConn) InitUploadProgress() {
	if g.jsonRds == nil {
		log.Errorf("json redis client nil")
	}

	client := g.jsonRds
	ctx := g.ctx

	// SET key val NX
	cmd := client.B().
		JsonSet().
		Key(UPLOAD_PROGRESS_LABEL).
		Path("$").
		Value(rueidis.JSON(struct{}{})).
		Build()
	err := client.Do(ctx, cmd).Error()

	if err != nil {
		log.Errorf("failed to set upload progress at redis \n error: %v", err)
	}
}

func SetUploadProgress(key string, val UploadProgressValue) {
	if rdsConn.jsonRds == nil {
		log.Errorf("json redis client nil")
	}

	client := rdsConn.jsonRds
	ctx := rdsConn.ctx

	str_arr := strings.Split(key, "/")
	user_uid := str_arr[0]
	file_uid := str_arr[1]

	// GET val
	get_cmd := client.B().
		JsonGet().
		Key(UPLOAD_PROGRESS_LABEL).
		Path(fmt.Sprintf("$.%s", user_uid)).
		Build()
	result, err := client.Do(ctx, get_cmd).ToString()

	// if redis don't have user_uid property
	if result == "[]" {
		// SET key val NX
		cmd := client.B().
			JsonSet().
			Key(UPLOAD_PROGRESS_LABEL).
			Path(fmt.Sprintf("$.%s", user_uid)).
			Value(rueidis.JSON(struct{}{})).
			Build()
		err = client.Do(ctx, cmd).Error()
	}
	// SET key val NX
	cmd := client.B().
		JsonSet().
		Key(UPLOAD_PROGRESS_LABEL).
		Path(fmt.Sprintf("$.%s.%s", user_uid, file_uid)).
		Value(rueidis.JSON(val)).
		Build()
	err = client.Do(ctx, cmd).Error()

	if err != nil {
		log.Errorf("failed to set upload progress at redis \n error: %v", err)
	}
}

func DelUploadProgress(key string) {
	if rdsConn.jsonRds == nil {
		log.Errorf("json redis client nil")
	}

	str_arr := strings.Split(key, "/")
	user_uid := str_arr[0]
	file_uid := str_arr[1]

	client := rdsConn.jsonRds
	ctx := rdsConn.ctx

	// SET key val NX
	cmd := client.B().
		JsonDel().
		Key(UPLOAD_PROGRESS_LABEL).
		Path(fmt.Sprintf("$.%s.%s", user_uid, file_uid)).
		Build()
	err := client.Do(ctx, cmd).Error()

	if err != nil {
		log.Errorf("failed to delete upload progress at redis \n error: %v", err)
	}
}

func GetUploadProgress(user_uid string) (string, error) {
	if rdsConn.jsonRds == nil {
		log.Errorf("json redis client nil")
	}
	client := rdsConn.jsonRds
	ctx := rdsConn.ctx

	cmd := client.B().
		JsonGet().
		Key(UPLOAD_PROGRESS_LABEL).
		Path(fmt.Sprintf("$.%s", user_uid)).
		Build()

	result, err := client.Do(ctx, cmd).ToString()

	return result[1 : len(result)-1], err
}
