package instance_service

import (
	"astralstack-hyperv/database"
	"astralstack-hyperv/model"
	"fmt"
	"github.com/pkg/errors"
	"github.com/shirou/gopsutil/mem"
)

const SafeRateThreshold = 0.9

type CreateInstanceDTO struct {
	*model.Instance
}

func PreCheckCreateInstance(dto CreateInstanceDTO) error {
	// 1.检测是否存在同名云主机实例
	if ExistsByAlias(dto.Alias) {
		return errors.New(fmt.Sprintf("Duplicated instance %s!", dto.Alias))
	}
	// 2. 检测资源是否安全
	// 2.1 检测内存是否安全
	v, err := mem.VirtualMemory()
	if err != nil {
		return err
	}
	need := dto.MemorySizeMB * 1024 * 1024 * 1024
	total := v.Total
	used := v.Used
	threshold := uint64(float64(total) * 0.9)
	if need+used > threshold {
		return errors.New("Resource memory not enough!")
	}
	// 2.2 检测硬盘是否安全

	return nil
}

func CreateInstance(dto CreateInstanceDTO) (*model.Instance, error) {
	if err := PreCheckCreateInstance(dto); err != nil {
		return nil, err
	}
	instance := &model.Instance{
		Name:             dto.Name,
		Alias:            dto.Alias,
		VCpuCount:        dto.VCpuCount,
		MemorySizeMB:     dto.MemorySizeMB,
		MaxBandwidthMbps: dto.MaxBandwidthMbps,
	}
	return instance, nil
}

func ExistsByAlias(alias string) bool {
	var count int64
	db := database.DB
	db.Model(&model.Instance{}).Where("alias = ?", alias).Count(&count)
	return count > 0
}
