package hyperv

import (
	"astralstack-hyperv/system/utils/wmi"
	"github.com/go-ole/go-ole"
)

//func NewVirtualEthernetAdapter() {
//	vsms, err := NewVMManager()
//	if err != nil {
//		return
//	}
//	host, err := vsms.GetHost()
//	if err != nil {
//		return
//	}
//	location, err := wmi.NewLocation(host.path)
//	if err != nil {
//		return
//	}
//	server := location.Server
//	namespace := location.Namespace
//	path := fmt.Sprintf("\\\\%s\\%s:%s", server, namespace, `Msvm_SyntheticEthernetPortSettingData.InstanceID="Microsoft:Definition\\6A45335D-4C3A-44B7-B61F-C9808BBDF8ED\\Default"`)
//	object, err := GetObject(path)
//	if err != nil {
//		return
//	}
//	text, err := object.GetText(1)
//	if err != nil {
//		return
//	}
//	_, err = AddResourceSetting(vsms.svc, host.path, []string{text})
//	if err != nil {
//		return
//	}
//}

func CreateVSwitch(name string) error {
	vsms, err := NewSwitchManager()
	con := vsms.con
	if err != nil {
		return err
	}
	instance, err := SpawnInstance("Msvm_VirtualEthernetSwitchSettingData", con)
	if err != nil {
		return err
	}
	if err := instance.Set("ElementName", name); err != nil {
		return err
	}
	vmText, err := instance.GetText(1)
	if err != nil {
		return err
	}
	jobPath := ole.VARIANT{}
	resultingSystem := ole.VARIANT{}
	jobState, err := vsms.svc.Get("DefineSystem", vmText, nil, nil, &resultingSystem, &jobPath)
	if err != nil {
		return err
	}
	if jobState.Value().(int32) == wmi.JobStatusStarted {
		err := wmi.WaitForJob(jobPath.Value().(string))
		if err != nil {
			return err
		}
	}
	// The resultingSystem value for DefineSystem is always a string containing the
	// location of the newly created resource
	locationURI := resultingSystem.Value().(string)
	loc, err := wmi.NewLocation(locationURI)
	if err != nil {
		return err
	}

	result, err := loc.GetResult()
	//_, err = loc.GetResult()
	if err != nil {
		return err
	}

	//portName, err := result.GetProperty("Name")
	//if err != nil {
	//	return err
	//}

	parent, err := result.Path()
	if err != nil {
		return err
	}

	vmManager, err := NewVMManager()
	if err != nil {
		return err
	}

	host, err := vmManager.GetHost()
	if err != nil {
		return err
	}

	settingData, err := SpawnInstance("Msvm_EthernetPortAllocationSettingData", vsms.con)
	path := host.path
	settingData.Set("ElementName", name)
	settingData.Set("ResourceType", 33)
	settingData.Set("ResourceSubType", "Microsoft:Hyper-V:Ethernet Connection")
	settingData.Set("HostResource", []string{path})
	settingData.Set("PortName", "7D21B94A-B258-4170-B221-80FC046D88DD")
	println(parent)

	text, err := settingData.GetText(1)
	if err != nil {
		return err
	}
	_, err = AddResourceSetting(vsms.svc, host.path, []string{text})
	if err != nil {
		return err
	}
	return nil
}
