package models

import "time"

type BaseModel struct {
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	CreateBy  string    `json:"createBy"`
	UpdateBy  string    `json:"updateBy"`
	ID        int       `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
}
