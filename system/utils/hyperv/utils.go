package hyperv

import (
	"astralstack-hyperv/system/utils/wmi"
	"fmt"
	"github.com/go-ole/go-ole"
	"github.com/pkg/errors"
)

// AddResourceSetting adds the resource settings to the specified VM
// affectedConfiguration: vmPath
func AddResourceSetting(svc *wmi.Result, affectedConfiguration string, settingsData []string) ([]string, error) {
	jobPath := ole.VARIANT{}
	resultingSystem := ole.VARIANT{}
	jobState, err := svc.Get("AddResourceSettings", affectedConfiguration, settingsData, &resultingSystem, &jobPath)
	if err != nil {
		return nil, errors.Wrap(err, "calling ModifyResourceSettings")
	}

	if jobState.Value().(int32) == wmi.JobStatusStarted {
		err := wmi.WaitForJob(jobPath.Value().(string))
		if err != nil {
			return nil, errors.Wrap(err, "waiting for job")
		}
	}
	safeArrayConversion := resultingSystem.ToArray()
	valArray := safeArrayConversion.ToValueArray()
	if len(valArray) == 0 {
		return nil, fmt.Errorf("no resource in resultingSystem value")
	}
	resultingSystems := make([]string, len(valArray))
	for idx, val := range valArray {
		resultingSystems[idx] = val.(string)
	}
	return resultingSystems, nil
}
