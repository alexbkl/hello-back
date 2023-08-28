package rds

import (
	"fmt"

	"github.com/redis/rueidis"
)

type InitState struct {
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
		Value(rueidis.JSON(InitState{})).
		Build()
	err := client.Do(ctx, cmd).Error()

	if err != nil {
		log.Errorf("failed to set upload progress at redis \n error: %v", err)
	}
}

func SetUploadProgress(key string, val int) {
	if rdsConn.jsonRds == nil {
		log.Errorf("json redis client nil")
	}

	client := rdsConn.jsonRds
	ctx := rdsConn.ctx

	// SET key val NX
	cmd := client.B().
		JsonSet().
		Key(UPLOAD_PROGRESS_LABEL).
		Path(fmt.Sprintf("$.%s", key)).
		Value(rueidis.JSON(val)).
		Build()
	err := client.Do(ctx, cmd).Error()

	if err != nil {
		log.Errorf("failed to set upload progress at redis \n error: %v", err)
	}
}

func DelUploadProgress(key string) {
	if rdsConn.jsonRds == nil {
		log.Errorf("json redis client nil")
	}

	client := rdsConn.jsonRds
	ctx := rdsConn.ctx

	// SET key val NX
	cmd := client.B().
		JsonDel().
		Key(UPLOAD_PROGRESS_LABEL).
		Path(fmt.Sprintf("$.%s", key)).
		Build()
	err := client.Do(ctx, cmd).Error()

	if err != nil {
		log.Errorf("failed to delete upload progress at redis \n error: %v", err)
	}
}
