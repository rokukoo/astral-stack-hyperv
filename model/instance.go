package model

import (
	"gorm.io/gorm"
)

// VCpuCount 虚拟cpu核心数
// MemorySizeMB 内存 (MB)
// MaxBandwidthMbps 带宽
// System Name 系统名称
//

type Status = string

const (
	Creating Status = "创建中"
	Starting        = "启动中"
	Running         = "运行中"
	Stopping        = "关闭中"
	Stopped         = "已停止"
)

type Instance struct {
	Name             string
	Alias            string
	Ipv4Address      string
	VCpuCount        int
	MemorySizeMB     uint64
	MaxBandwidthMbps int64
	SystemName       string
	Username         string
	Password         string
	Status           Status
	VirtualHardDisks *[]VirtualHardDisk
	gorm.Model
}
