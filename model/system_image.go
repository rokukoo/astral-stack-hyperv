package model

import (
	"gorm.io/gorm"
)

type SystemImage struct {
	Tag   string // 系统标签, windows server, ubuntu, ...
	Alias string // 系统镜像别名
	Name  string // 系统镜像名称
	Path  string // 系统镜像路径
	Size  int    // 系统路径大小
	gorm.Model
}
