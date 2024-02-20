package hyperv

import (
	"astralstack-hyperv/system/utils/wmi"
	"fmt"
	"github.com/go-ole/go-ole"
	"github.com/pkg/errors"
)

const DiskResType = 17

//var Con, _ = wmi.NewVirtualizationConnection()

type SyntheticDiskDrive struct {
	pth  string
	vm   *VirtualMachine
	ctrl *wmi.Result
	slot IDE_CONTROLLER_SLOT // AddressOnParent
}

func GetIDEController(vm *VirtualMachine, address IDE_CONTROLLER_NUM) (*wmi.Result, error) {
	qParams := []wmi.Query{
		&wmi.AndQuery{
			wmi.QueryFields{Key: "ResourceType", Value: 5, Type: wmi.Equals},
		},
		&wmi.AndQuery{
			wmi.QueryFields{Key: "InstanceID", Value: fmt.Sprintf(`Microsoft:%s%%`, vm.instanceId), Type: wmi.Like},
		},
		&wmi.AndQuery{
			wmi.QueryFields{Key: "Address", Value: address, Type: wmi.Equals},
		},
	}
	con := vm.mgr.con
	result, err := con.GetOne(ResourceAllocSettingDataClass, nil, qParams)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func ExistsIDEController(vm *VirtualMachine, slot IDE_CONTROLLER_SLOT) (bool, error) {
	ideController, err := GetIDEController(vm, slot)
	return ideController != nil, err
}

func NewSyntheticDiskDrive(vm *VirtualMachine, ctrl *wmi.Result, slot string) (*SyntheticDiskDrive, error) {
	if _, err := GetIDEController(vm, slot); err != nil {
		return nil, errors.New("Not found IDE0 controller!")
	}
	return &SyntheticDiskDrive{
		vm:   vm,
		slot: slot,
		ctrl: ctrl,
	}, nil
}

func (s *SyntheticDiskDrive) Create() (bool, error) {
	drv, err := NewResourceSettingData(s.vm.mgr.con)
	if err != nil {
		return false, err
	}
	if err := drv.Set("AddressOnParent", s.slot); err != nil {
		return false, errors.Wrap(err, "Set AddressOnParent")
	}
	if err := drv.Set("ResourceType", DiskResType); err != nil {
		return false, errors.Wrap(err, "Set ResourceType")
	}
	if err := drv.Set("ResourceSubType", DiskResSubtype); err != nil {
		return false, errors.Wrap(err, "Set ResourceSubType")
	}
	parentPath, err := s.ctrl.Path()
	if err != nil {
		return false, err
	}
	if err := drv.Set("Parent", parentPath); err != nil {
		return false, errors.Wrap(err, "Set Parent")
	}
	affectedConfiguration := s.vm.path
	drvText, err := drv.GetText(1)
	if err != nil {
		return false, err
	}
	resourceSettings := []string{drvText}
	resultingResourceSettings, err := AddResourceSetting(s.vm.mgr.svc, affectedConfiguration, resourceSettings)
	if err != nil {
		return false, err
	}
	s.pth = resultingResourceSettings[0]
	return true, nil
}

func AttachSyntheticDiskDrive() {

}

func CreateSyntheticDiskDrive(address int) {

}

func CreateHardDisk(name string, id string, path string) (*wmi.Result, error) {
	imm, err := NewImageManager()
	if err != nil {
		return nil, err
	}
	virtualHardDiskInstance, err := imm.con.Get(VirtualHardDiskSettingDataClass)
	if err != nil {
		return nil, err
	}
	newVirtualHardDisk, err := virtualHardDiskInstance.Get("SpawnInstance_")
	if err != nil {
		return nil, errors.Wrap(err, "calling SpawnInstance_")
	}
	if err := newVirtualHardDisk.Set("Format", 3); err != nil {
		return nil, errors.Wrap(err, "Set Format")
	}
	if err := newVirtualHardDisk.Set("Type", 3); err != nil {
		return nil, errors.Wrap(err, "Set Type")
	}
	if err := newVirtualHardDisk.Set("Caption", name); err != nil {
		return nil, errors.Wrap(err, "Set Caption")
	}
	if err := newVirtualHardDisk.Set("VirtualDiskId", id); err != nil {
		return nil, errors.Wrap(err, "Set VirtualDiskId")
	}
	if err := newVirtualHardDisk.Set("Path", fmt.Sprintf(`%s\%s.vhdx`, path, name)); err != nil {
		return nil, errors.Wrap(err, "Set Path")
	}
	if err := newVirtualHardDisk.Set("MaxInternalSize", 1*1024*1024*1024); err != nil {
		return nil, errors.Wrap(err, "Set BlockSize")
	}
	vhdText, err := newVirtualHardDisk.GetText(1)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to get Virtual Hard Disk instance XML")
	}
	jobPath := ole.VARIANT{}
	jobState, err := imm.svc.Get("CreateVirtualHardDisk", vhdText, &jobPath)
	if err != nil {
		return nil, errors.Wrap(err, "calling DefineSystem")
	}
	if jobState.Value().(int32) == wmi.JobStatusStarted {
		err := wmi.WaitForJob(jobPath.Value().(string))
		if err != nil {
			return nil, errors.Wrap(err, "waiting for job")
		}
	}
	return nil, nil
}
