package models

import "goblog/types"

type BaseModel struct {
	ID int64
}

func (base BaseModel) GetStringID() string {
	return types.Int64ToString(base.ID)
}
