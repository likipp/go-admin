package global

import (
	"go-admin/config"
	"gorm.io/gorm"
)

var (
	GLA_DB     *gorm.DB
	GLA_CONFIG config.Server
)
