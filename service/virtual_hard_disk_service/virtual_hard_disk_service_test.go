package virtual_hard_disk_service

import (
	"astralstack-hyperv/boostrap"
	"testing"
)

func TestCreate(t *testing.T) {
	err := boostrap.Init()
	if err != nil {
		panic(err)
	}
	_, err = add("cn_windows_server_2019_datacenter_x64", `D:\Hyper-V\cn_windows_server_2019_datacenter_x64\Virtual Hard Disks\System.vhdx`, 20*1024*1024*1024, true)
	if err != nil {
		panic(err)
	}
}
