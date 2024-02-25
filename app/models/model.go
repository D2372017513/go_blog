package models

import (
	"goblog/pkg/types"
	"time"
)

type BaseModel struct {
	ID       int64     `gorm:"column:id;primaryKey;autoIncrement;not null"`
	CreateAt time.Time `gorm:"column:create_at;index;autoCreateTime:milli"`
	UpdateAt time.Time `gorm:"column:update_at;index;autoUpdateTime:milli"`
}

func (base BaseModel) GetStringID() string {
	return types.Int64ToString(base.ID)
}
