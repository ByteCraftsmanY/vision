package entities

import (
	"time"
	"vision/types"
)

type Base struct {
	ID        *types.ID          `json:"id,omitempty" db:"id,part"`
	IsActive  bool               `json:"is_active,omitempty" db:"is_active"`
	CreatedAt *time.Time         `json:"created_at,omitempty" db:"created_at"`
	UpdatedAt *time.Time         `json:"updated_at,omitempty" db:"updated_at"`
	DeletedAt *time.Time         `json:"deleted_at,omitempty" db:"deleted_at"`
	Extra     *map[string]string `json:"extra,omitempty" db:"extra"`
}

func (b *Base) Initialize() {
	currentTime := time.Now()
	b.ID = types.CreateNewID()
	b.IsActive = true
	b.CreatedAt = &currentTime
	b.UpdatedAt = &currentTime
}

func (b *Base) Update() {
	currentTime := time.Now()
	b.UpdatedAt = &currentTime
}

func (b *Base) Delete() {
	currentTime := time.Now()
	b.DeletedAt = &currentTime
	b.IsActive = false
}
