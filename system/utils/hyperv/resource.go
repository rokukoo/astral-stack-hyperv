package hyperv

import "astralstack-hyperv/system/utils/wmi"

func SpawnInstance(className string, con *wmi.WMI) (*wmi.Result, error) {
	class, err := con.Get(className, nil, nil)
	if err != nil {
		return nil, err
	}
	return class.Get("SpawnInstance_")
}

func NewResourceSettingData(con *wmi.WMI) (*wmi.Result, error) {
	return SpawnInstance(ResourceAllocSettingDataClass, con)
}
