package model

import (
	"gorm.io/gorm"
	"time"

	"github.com/gofrs/uuid"
)

type Model interface {
	GetID() string
	BeforeCreate(tx *gorm.DB) error
	BeforeUpdate(tx *gorm.DB) error
	RevisionKey() string
}

type BaseModel struct {
	Id        string    `gorm:"type:uuid;primary_key" json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (m *BaseModel) GetID() string {
	return m.Id
}

func (m *BaseModel) BeforeCreate(tx *gorm.DB) error {
	uuid4, err := uuid.NewV4()
	if err != nil {
		return err
	}

	if m.Id == "" {
		tx.Statement.SetColumn("Id", uuid4.String())
	}

	now := time.Now()
	now = now.Round(time.Second)
	if m.CreatedAt.IsZero() {
		tx.Statement.SetColumn("CreatedAt", now)
	}
	if m.UpdatedAt.IsZero() {
		tx.Statement.SetColumn("UpdatedAt", now)
	}
	return nil
}

func (m *BaseModel) BeforeUpdate(tx *gorm.DB) error {
	now := time.Now()
	now = now.Round(time.Second)
	tx.Statement.SetColumn("UpdatedAt", now)
	return nil
}

func (m *BaseModel) RevisionKey() string {
	return m.Id
}
