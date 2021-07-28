package models

type CasbinModel struct {
	ID       int    `gorm:"primaryKey;autoIncrement" json:"id"`
	RoleName string `gorm:"role_name" json:"role_name"`
	Path     string `gorm:"path" json:"path"`
	Method   string `gorm:"method" json:"method"`
}
