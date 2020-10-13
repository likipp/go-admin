package models

type BaseMenu struct {
	BaseModel
	Icon       string `json:"icon" gorm:"column:icon"`
	MenuLevel  uint   `json:"-"`
	Sequence   int    `json: "sequence" gorm:"column:sequence;index;default:0;not null;"`
	ParentId   string `json:"parent_id" gorm:"comment:父菜单ID;index"`
	ParentPath string `json:"parent_path" gorm:"comment:路由path;index"`
	Name       string `json:"name" gorm:"comment:路由name;index"`
	Hidden     bool   `json:"hidden" gorm:"comment:是否在列表隐藏;index" binding:"required,max=2,min=1"`
	Status     int    `json: "status" gorm:"column:status;index;default:0;not null;" "binding:"required,max=2,min=1"`
	Component  string `json:"component" gorm:"comment:对应前端文件路径;index"`
	//SysRoles     []SysRole     `json:"roles" gorm:"many2many:sys_authority_menus;"`
	//Children     []BaseMenu    `json:"children" gorm:"-"`
}
