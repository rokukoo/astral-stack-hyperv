package hyperv

import (
	"astralstack-hyperv/system/utils/wmi"
	"github.com/pkg/errors"
)

// Manager offers a root\virtualization\v2 instance connection
// and an instance of Msvm_VirtualSystemManagementService/Msvm_ImageManagementService
type Manager struct {
	con *wmi.WMI
	svc *wmi.Result
}

type VMManager struct {
	*Manager
}

type ImageManager struct {
	*Manager
}

type SwitchManager struct {
	*Manager
}

// NewVMManager returns a new VMManager type
func NewVMManager() (*VMManager, error) {
	w, err := wmi.NewVirtualizationConnection()
	// Get virtual machine management service
	svc, err := w.GetOne(VMManagementService, []string{}, []wmi.Query{})
	if err != nil {
		//return nil, err
		return nil, errors.New("Hyper-v service is not enabled!")
	}
	sw := &VMManager{
		&Manager{
			con: w,
			svc: svc,
		},
	}
	return sw, nil
}

func GetObject(path string) (*wmi.Result, error) {
	loc, err := wmi.NewLocation(path)
	if err != nil {
		return nil, errors.Wrap(err, "getting location")
	}
	result, err := loc.GetResult()
	if err != nil {
		return nil, errors.Wrap(err, "getting result")
	}
	return result, nil
}

func (vmm *VMManager) GetHost() (*VirtualMachine, error) {
	hs, err := vmm.con.GetOne("Msvm_HostedService", nil, nil)
	if err != nil {
		return nil, err
	}
	path, err := hs.GetProperty("Antecedent")
	if err != nil {
		return nil, err
	}
	sPath := path.Value().(string)
	result, err := GetObject(sPath)
	if err != nil {
		return nil, errors.Wrap(err, "getting object")
	}
	instanceId, err := result.GetProperty("Name")
	if err != nil {
		return nil, err
	}
	return &VirtualMachine{
		mgr:            vmm,
		computerSystem: result,
		// activeSettingsData:
		path:       sPath,
		instanceId: instanceId.Value().(string),
	}, nil
}

// NewImageManager returns a new ImageManager type
func NewImageManager() (*ImageManager, error) {
	w, err := wmi.NewVirtualizationConnection()
	if err != nil {
		return nil, err
	}
	// Get image management service
	svc, err := w.GetOne(ImageManagementServiceClass, []string{}, []wmi.Query{})
	if err != nil {
		//return nil, err
		return nil, errors.New("Hyper-v service is not enabled!")
	}
	sw := &ImageManager{
		&Manager{
			con: w,
			svc: svc,
		},
	}
	return sw, nil
}

// NewSwitchManager returns a new SwitchManager type
func NewSwitchManager() (*SwitchManager, error) {
	w, err := wmi.NewVirtualizationConnection()
	if err != nil {
		return nil, err
	}
	// Get image management service
	svc, err := w.GetOne(VirtualEthernetSwitchManagementServiceClass, []string{}, []wmi.Query{})
	if err != nil {
		//return nil, err
		return nil, errors.New("Hyper-v service is not enabled!")
	}
	sw := &SwitchManager{
		&Manager{
			con: w,
			svc: svc,
		},
	}
	return sw, nil
}
