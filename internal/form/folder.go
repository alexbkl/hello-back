package form

import (
		"github.com/Hello-Storage/hello-back/internal/entity"
)


type CreateFolder struct {
	Title string `json:"title"`
	Root  string `json:"root"`
	Status    entity.EncryptionStatus `json:"status"`
}

type UpdateFolder struct {
	Id   string `json:"id"`
	Uid  string `json:"uid"`
	Root string `json:"root"`
}