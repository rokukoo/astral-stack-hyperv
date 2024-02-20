package boostrap

import (
	"astralstack-hyperv/database"
)

func Init() error {
	// 1. 初始化配置文件
	if err := initConfig(); err != nil {
		return err
	}
	// 2. 初始化数据库
	if err := initDb(); err != nil {
		return err
	}
	return nil
}

func initConfig() error {

	return nil
}

func initDb() error {
	if err := database.Init(); err != nil {
		return err
	}
	return nil
}
