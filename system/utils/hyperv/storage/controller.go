package vm

import "astralstack-hyperv/system/utils/hyperv"

type StorageController struct {
	mgr  *hyperv.Manager
	path string
}

type IDEController struct {
	*StorageController
}

type SCSIController struct {
	*StorageController
}

// Path returns the WMI path of this controller
func (s *StorageController) Path() string {
	return s.path
}
