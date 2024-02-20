package hyperv

import (
	"testing"
)

func TestCreateHardDisk(t *testing.T) {
	//name := "1111"
	//uuid := strings.ToUpper(uuid2.New().String())
	//diskName := "Data_" + uuid
	//path := fmt.Sprintf(`D:\Hyper-V\%s\Virtual Hard Disks`, name)
	//_, err := CreateHardDisk(diskName, uuid, path)
	//if err != nil {
	//	return
	//}
	//drive, err := NewSyntheticDiskDrive()
	//_, err = drive.Create()
	//if err != nil {
	//	println(err)
	//}
}

func TestExistsIDEController(t *testing.T) {
	name := "1111"
	vmm, err := NewVMManager()
	if err != nil {
		panic(err)
	}
	vm, err := vmm.GetVmByName(name)
	if err != nil {
		panic(err)
	}
	ctrl, err := vm.GetIDEController0()
	if err != nil {
		panic(err)
	}
	drive, err := NewSyntheticDiskDrive(vm, ctrl, "1")
	if err != nil {
		panic(err)
	}
	_, err = drive.Create()
	if err != nil {
		panic(err)
	}
	vhdPath := `D:\Hyper-V\1111\Virtual Hard Disks\Data_37EBE251-3FF0-4EAB-80FF-104BC1FF0CFC.vhdx`
	sasd, err := SpawnInstance("Msvm_StorageAllocationSettingData", vmm.con)
	sasd.Set("ResourceType", 31)
	sasd.Set("ResourceSubType", IDEDiskResSubType)
	sasd.Set("HostResource", []string{vhdPath})
	sasd.Set("Parent", drive.pth)
	text, _ := sasd.GetText(1)
	AddResourceSetting(vmm.svc, vm.path, []string{text})
	//imm, err := NewImageManager()
	//if err != nil {
	//	return
	//}
	//imm
}
