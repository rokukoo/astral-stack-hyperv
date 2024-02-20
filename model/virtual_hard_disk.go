package model

import (
	"gorm.io/gorm"
)

type VirtualHardDisk struct {
	Uuid       string  `json:"uuid"`          // uid
	IsSystem   bool    `json:"is_system"`     // 是否为系统盘
	Alias      string  `json:"alias"`         // 名称
	Path       string  `json:"path"`          // 路径
	Used       uint64  `json:"used" gorm:"-"` // 外部字段
	Total      uint64  `json:"total"`         // 大小
	InstanceId *uint   `json:"instance_id"`
	Comment    *string `json:"comment"`
	gorm.Model
}
