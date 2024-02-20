package virtual_hard_disk_service

import (
	"astralstack-hyperv/database"
	"astralstack-hyperv/model"
	"fmt"
	"github.com/google/uuid"
	"log"
	"os"
)

func init() {
	path := `D:\Hyper-V\virtual hard disks\`
	_, err := os.Stat(path)
	if err != nil {
		log.Panicln(err)
	}
	if os.IsNotExist(err) {
		// 创建文件夹
		err := os.MkdirAll(path, os.ModePerm)
		if err != nil {
			fmt.Printf("mkdir failed![%v]\n", err)
		}
	}
}

func List() ([]model.VirtualHardDisk, error) {
	var vhdList []model.VirtualHardDisk
	err := database.DB.Find(&vhdList).Error
	if err != nil {
		return nil, err
	}
	return vhdList, nil
}

func DeleteById(id uint) {

}

func Create(name string, path string, size int32, isSystem bool) error {

	return nil
}

func Import(name string, path string, size int32, isSystem bool) error {

	return nil
}

func add(alias string, path string, total uint64, isSystem bool) (*model.VirtualHardDisk, error) {

	uid := uuid.New().String()
	vhd := &model.VirtualHardDisk{
		Uuid:     uid,
		Path:     path,
		Alias:    alias,
		IsSystem: isSystem,
		Total:    total,
	}
	db := database.DB
	tx := db.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	err := tx.Create(&vhd).Error

	if err != nil {
		log.Panicln("虚拟硬盘创建失败!")
	}

	tx.Commit()
	return vhd, nil
}
