package dtos

import (
	"time"
	"vision/types"
)

type BaseURI struct {
	ID string `uri:"id" binding:"required,uuid"`
}

type Base struct {
	ID        *types.ID  `json:"id,omitempty"`
	IsActive  bool       `json:"is_active,omitempty"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}
