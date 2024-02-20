package compute

import (
	"astralstack-hyperv/system/utils/hyperv"
	"astralstack-hyperv/system/utils/wmi"
)

type Instance struct {
	Name  string
	State hyperv.State
}

func CreateInstance(name string) (*Instance, error) {
	instance, err := NewInstance(name)

	return instance, err
}

func NewInstance(name string) (*Instance, error) {
	return &Instance{
		Name: name,
	}, nil
}

type InstanceList = []Instance

func wrap(ref *wmi.Result) *Instance {
	elementName, _ := ref.GetProperty("ElementName")
	enabledState, _ := ref.GetProperty("EnabledState")
	ins := &Instance{
		Name:  elementName.Value().(string),
		State: enabledState.Value().(int32),
	}
	return ins
}

func GetByName(name string) (*Instance, error) {
	vmRef, err := hyperv.GetVmRefByName(name)
	if err != nil {
		return nil, err
	}
	var instance = wrap(vmRef)
	return instance, nil
}

func GetAllInstance() (*InstanceList, error) {
	var list InstanceList
	result, err := hyperv.Con.Gwmi("Msvm_ComputerSystem", nil, nil)
	count, err := result.Count()
	if err != nil {
		return nil, err
	}
	for i := 0; i < count; i++ {
		itemAtIndex, err := result.ItemAtIndex(i)
		if err != nil {
			return nil, err
		}
		ins := wrap(itemAtIndex)
		list = append(list, *ins)
	}
	return &list, nil
}

func (i *Instance) Start() (bool, error) {
	return hyperv.StartVM(i.Name)
}

func (i *Instance) Stop() (bool, error) {
	return hyperv.StopVM(i.Name)
}

func (i *Instance) ReStart() (bool, error) {
	return hyperv.ReStartVM(i.Name)
}
