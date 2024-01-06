package models

import (
	"github.com/gocql/gocql"
	"time"
)

type Base struct {
	UUID      gocql.UUID `json:"uuid,omitempty" db:"uuid"`
	IsActive  bool       `json:"is_active,omitempty" db:"is_active"`
	CreatedAt time.Time  `json:"created_at,omitempty" db:"created_at"`
	UpdatedAt time.Time  `json:"updated_at,omitempty" db:"updated_at"`
	DeletedAT time.Time  `json:"deleted_at,omitempty" db:"deleted_at"`
}

func (b *Base) Initialize() {
	b.UUID = gocql.UUIDFromTime(time.Now())
	b.IsActive = true
	b.CreatedAt = time.Now()
	b.UpdatedAt = time.Now()
}

func (b *Base) Update() {
	b.UpdatedAt = time.Now()
}

func (b *Base) Delete() {
	b.DeletedAT = time.Now()
	b.IsActive = false
}
