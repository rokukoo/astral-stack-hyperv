package compute

func PreCheckCreateInstance(instance *Instance) (bool, error) {
	// 1. 检测是否存在同名实例
	if duplicated, err := checkNameDuplicated(); err != nil || duplicated {
		return false, err
	}
	// 2. 检测是否存在可用空间
	if engh, err := checkDiskEnough(); err != nil || !engh {
		return false, err
	}
	// 3. 检测内存是否可用
	if engh, err := checkMemEnough(); err != nil || !engh {
		return false, err
	}
	return true, nil
}

func checkNameDuplicated() (bool, error) {

	return false, nil
}

func checkDiskEnough() (bool, error) {

	return true, nil
}

func checkMemEnough() (bool, error) {

	return true, nil
}
