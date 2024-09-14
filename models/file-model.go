package models

import (
	"time"

	"gorm.io/gorm"
)

type File struct {
	gorm.Model

	Createdby string `json:"createdby"`
	Key       string `json:"key"`
	Url       string `json:"url"`
	ExpirationDate time.Time `json:"expirationDate"`
}
