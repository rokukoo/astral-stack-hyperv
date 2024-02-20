package hyperv

import (
	"astralstack-hyperv/system/utils/wmi"
	"fmt"
	"github.com/go-ole/go-ole"
	"github.com/pkg/errors"
)

var Con, _ = wmi.NewVirtualizationConnection()

// VirtualMachine represents a single virtual machine
type VirtualMachine struct {
	mgr *VMManager

	instanceId         string
	activeSettingsData *wmi.Result
	computerSystem     *wmi.Result
	path               string
}

func (vm *VirtualMachine) GetIDEController0() (*wmi.Result, error) {
	return GetIDEController(vm, IDE_CTRL_0)
}

func (m *VMManager) GetVmByName(name string) (*VirtualMachine, error) {
	var qParams []wmi.Query
	vmRefByName, err := GetVmRefByName(name)
	if err != nil {
		return nil, err
	}
	instanceID, err := vmRefByName.GetProperty("Name")
	instanceId := instanceID.Value().(string)
	if err != nil {
		return nil, err
	}
	qParams = []wmi.Query{
		&wmi.AndQuery{
			wmi.QueryFields{
				Key:   "VirtualSystemIdentifier",
				Value: instanceId,
				Type:  wmi.Equals},
		},
	}
	result, err := m.con.Gwmi(VirtualSystemSettingDataClass, nil, qParams)
	if err != nil {
		return nil, errors.Wrap(err, "VirtualSystemSettingDataClass")
	}

	vssd, err := result.ItemAtIndex(0)
	if err != nil {
		return nil, errors.Wrap(err, "fetching element")
	}
	pth, err := vmRefByName.Path()
	if err != nil {
		return nil, errors.Wrap(err, "VM path")
	}
	return &VirtualMachine{
		instanceId:         instanceId,
		mgr:                m,
		activeSettingsData: vssd,
		computerSystem:     vmRefByName,
		path:               pth,
	}, nil
}

// GetVM returns the virtual machine identified by instanceID
func (m *VMManager) GetVM(instanceID string) (*VirtualMachine, error) {
	qParams := []wmi.Query{
		&wmi.AndQuery{
			wmi.QueryFields{
				Key:   "VirtualSystemIdentifier",
				Value: instanceID,
				Type:  wmi.Equals},
		},
	}

	result, err := m.con.Gwmi(VirtualSystemSettingDataClass, nil, qParams)
	if err != nil {
		return nil, errors.Wrap(err, "VirtualSystemSettingDataClass")
	}

	vssd, err := result.ItemAtIndex(0)
	if err != nil {
		return nil, errors.Wrap(err, "fetching element")
	}
	cs, err := vssd.Get("associators_", nil, ComputerSystemClass)
	if err != nil {
		return nil, errors.Wrap(err, "getting ComputerSystemClass")
	}
	elem, err := cs.Elements()
	if err != nil || len(elem) == 0 {
		return nil, errors.Wrap(err, "getting elements")
	}
	pth, err := elem[0].Path()
	if err != nil {
		return nil, errors.Wrap(err, "VM path")
	}
	return &VirtualMachine{
		instanceId:         instanceID,
		mgr:                m,
		activeSettingsData: vssd,
		computerSystem:     elem[0],
		path:               pth,
	}, nil
}

func DeleteVM(name string) {

}

// SetMemory sets the virtual machine memory allocation
func (v *VirtualMachine) SetMemory(memoryMB int64) error {
	memorySettingsResults, err := v.activeSettingsData.Get("associators_", nil, MemorySettingDataClass)
	if err != nil {
		return errors.Wrap(err, "getting MemorySettingDataClass")
	}

	memorySettings, err := memorySettingsResults.ItemAtIndex(0)
	if err != nil {
		return errors.Wrap(err, "ItemAtIndex")
	}

	if err := memorySettings.Set("Limit", memoryMB); err != nil {
		return errors.Wrap(err, "Limit")
	}

	if err := memorySettings.Set("Reservation", memoryMB); err != nil {
		return errors.Wrap(err, "Reservation")
	}

	if err := memorySettings.Set("VirtualQuantity", memoryMB); err != nil {
		return errors.Wrap(err, "VirtualQuantity")
	}

	if err := memorySettings.Set("DynamicMemoryEnabled", false); err != nil {
		return errors.Wrap(err, "DynamicMemoryEnabled")
	}

	memText, err := memorySettings.GetText(1)
	if err != nil {
		return errors.Wrap(err, "Failed to get VM instance XML")
	}

	return v.modifyResourceSettings([]string{memText})
}

func (v *VirtualMachine) modifyResourceSettings(settings []string) error {
	jobPath := ole.VARIANT{}
	resultingSystem := ole.VARIANT{}
	jobState, err := v.mgr.svc.Get("ModifyResourceSettings", settings, &resultingSystem, &jobPath)
	if err != nil {
		return errors.Wrap(err, "calling ModifyResourceSettings")
	}
	if jobState.Value().(int32) == wmi.JobStatusStarted {
		err := wmi.WaitForJob(jobPath.Value().(string))
		if err != nil {
			return errors.Wrap(err, "waiting for job")
		}
	}
	return nil
}

// SetCPUs sets the number of CPU cores on the VM
func (v *VirtualMachine) SetCPUs(cpus int, limitCPUFeatures bool) error {
	//hostCpus := runtime.NumCPU()
	//if hostCpus < cpus {
	//	return fmt.Errorf("Number of cpus exceeded available host resources")
	//}

	procSettingsResults, err := v.activeSettingsData.Get("associators_", nil, ProcessorSettingDataClass)
	if err != nil {
		return errors.Wrap(err, "getting ProcessorSettingDataClass")
	}

	procSettings, err := procSettingsResults.ItemAtIndex(0)
	if err != nil {
		return errors.Wrap(err, "ItemAtIndex")
	}

	if err := procSettings.Set("VirtualQuantity", uint64(cpus)); err != nil {
		return errors.Wrap(err, "VirtualQuantity")
	}

	if err := procSettings.Set("Reservation", cpus); err != nil {
		return errors.Wrap(err, "Reservation")
	}

	// Use 100% of CPU core
	if err := procSettings.Set("Limit", 100000); err != nil {
		return errors.Wrap(err, "Limit")
	}

	if err := procSettings.Set("LimitProcessorFeatures", limitCPUFeatures); err != nil {
		return errors.Wrap(err, "LimitProcessorFeatures")
	}

	procText, err := procSettings.GetText(1)
	if err != nil {
		return errors.Wrap(err, "Failed to get VM instance XML")
	}
	return v.modifyResourceSettings([]string{procText})
}

func CreateVM(name string, cpu int, memMB int64, generation GenerationType) (*wmi.Result, error) {
	ref, err := GetVmRefByName(name)
	if ref != nil {
		return nil, errors.New("Error occurred when create the hyperv! (Reason: Name duplicated.)")
	}
	vmm, err := NewVMManager()
	if err != nil {
		return nil, err
	}
	vmSettingsDataInstance, err := vmm.con.Get(VirtualSystemSettingDataClass)
	if err != nil {
		return nil, err
	}
	newVMInstance, err := vmSettingsDataInstance.Get("SpawnInstance_")
	if err != nil {
		return nil, errors.Wrap(err, "calling SpawnInstance_")
	}
	if err := newVMInstance.Set("ElementName", name); err != nil {
		return nil, errors.Wrap(err, "Set ElementName")
	}
	if err := newVMInstance.Set("VirtualSystemSubType", string(generation)); err != nil {
		return nil, errors.Wrap(err, "Set VirtualSystemSubType")
	}
	if err := newVMInstance.Set("ConfigurationDataRoot", fmt.Sprintf(`D:\Hyper-V\%s`, name)); err != nil {
		return nil, errors.Wrap(err, "Set ConfigurationDataRoot")
	}
	vmText, err := newVMInstance.GetText(1)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to get VM instance XML")
	}
	jobPath := ole.VARIANT{}
	resultingSystem := ole.VARIANT{}
	jobState, err := vmm.svc.Get("DefineSystem", vmText, nil, nil, &resultingSystem, &jobPath)
	if err != nil {
		return nil, errors.Wrap(err, "calling DefineSystem")
	}
	if jobState.Value().(int32) == wmi.JobStatusStarted {
		err := wmi.WaitForJob(jobPath.Value().(string))
		if err != nil {
			return nil, errors.Wrap(err, "waiting for job")
		}
	}
	// The resultingSystem value for DefineSystem is always a string containing the
	// location of the newly created resource
	locationURI := resultingSystem.Value().(string)
	loc, err := wmi.NewLocation(locationURI)
	if err != nil {
		return nil, errors.Wrap(err, "getting location")
	}

	result, err := loc.GetResult()
	if err != nil {
		return nil, errors.Wrap(err, "getting result")
	}
	// The name field of the returning class is actually the InstanceID...
	id, err := result.GetProperty("Name")
	if err != nil {
		return nil, errors.Wrap(err, "fetching VM ID")
	}

	vm, err := vmm.GetVM(id.Value().(string))
	if err != nil {
		return nil, errors.Wrap(err, "fetching VM")
	}

	if err := vm.SetMemory(memMB); err != nil {
		return nil, errors.Wrap(err, "setting memory limit")
	}

	if err := vm.SetCPUs(cpu, false); err != nil {
		return nil, errors.Wrap(err, "setting CPU limit")
	}

	return vm.computerSystem, nil
}

func GetVmRefByName(name string) (*wmi.Result, error) {
	qParams := []wmi.Query{
		&wmi.AndQuery{
			QueryFields: wmi.QueryFields{Key: "ElementName", Value: name, Type: wmi.Equals},
		},
	}
	result, err := Con.GetOne(ComputerSystemClass, nil, qParams)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func ReStartVM(name string) (bool, error) {
	_, err := StopVM(name)
	if err != nil {
		return false, nil
	}
	return StartVM(name)
}

func StartVM(name string) (bool, error) {
	vm, err := GetVmRefByName(name)
	if err != nil {
		return false, err
	}
	jobPath := ole.VARIANT{}
	jobState, _ := vm.Get("RequestStateChange", StateRunning, &jobPath)
	if jobState.Value().(int32) == wmi.JobStatusStarted {
		err := wmi.WaitForJob(jobPath.Value().(string))
		if err != nil {
			return false, errors.Wrap(err, "waiting for job")
		}
	}
	return true, nil
}

func StopVM(name string) (bool, error) {
	vm, err := GetVmRefByName(name)
	if err != nil {
		return false, err
	}
	jobPath := ole.VARIANT{}
	jobState, _ := vm.Get("RequestStateChange", StateOff, &jobPath)
	if jobState.Value().(int32) == wmi.JobStatusStarted {
		err := wmi.WaitForJob(jobPath.Value().(string))
		if err != nil {
			return false, errors.Wrap(err, "waiting for job")
		}
	}
	return true, nil
}
