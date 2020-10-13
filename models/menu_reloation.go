package models

type MenuRelation struct {
	BaseModel
	RoleID   string `gorm:"column:role_id"`
	UserID   string `gorm:"column:user_id"`
	MenuID   string `gorm:"column:menu_id"`
	ActionID string `gorm:"column:action_id"`
}
