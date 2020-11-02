package model

import (
	"time"

	"github.com/gofrs/uuid"
	"github.com/jinzhu/gorm"
)

type Model interface {
	GetId() string
	BeforeCreate(scope *gorm.Scope) error
	BeforeUpdate(scope *gorm.Scope) error
	RevisionKey() string
}

type BaseModel struct {
	Id        string    `gorm:"type:uuid;primary_key" json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (m *BaseModel) GetId() string {
	return m.Id
}

func (m *BaseModel) BeforeCreate(scope *gorm.Scope) error {
	uuid4, err := uuid.NewV4()
	if err != nil {
		return err
	}

	if m.Id == "" {
		err := scope.SetColumn("Id", uuid4.String())
		if err != nil {
			return err
		}
	}

	now := time.Now()
	now = now.Round(time.Second)
	if m.CreatedAt.IsZero() {
		err := scope.SetColumn("CreatedAt", now)
		if err != nil {
			return err
		}
	}
	if m.UpdatedAt.IsZero() {
		err := scope.SetColumn("UpdatedAt", now)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *BaseModel) BeforeUpdate(scope *gorm.Scope) error {
	now := time.Now()
	now = now.Round(time.Second)
	err := scope.SetColumn("UpdatedAt", now)
	if err != nil {
		return err
	}
	return nil
}

func (m *BaseModel) RevisionKey() string {
	return m.Id
}
