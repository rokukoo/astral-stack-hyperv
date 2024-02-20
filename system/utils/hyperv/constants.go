package hyperv

type State = int32

type IDE_CONTROLLER_SLOT = string

const (
	IDE_SLOT_0 IDE_CONTROLLER_SLOT = "0"
	IDE_SLOT_1                     = "1"
)

type IDE_CONTROLLER_NUM = string

const (
	IDE_CTRL_0 IDE_CONTROLLER_NUM = "0"
	IDE_CTRL_1                    = "1"
)

const (
	StateUnknown State = iota
	StateOther
	StateRunning
	StateOff
	StateShuttingDown
	StateNotApplicable
	StateEnabledButOffline
	StateInTest
	StateDeferred
	StateQuiesce
	StateStarting
	StateReset
	StateSaving   = 32773
	StatePausing  = 32776
	StateResuming = 32777
)

// Hyper-V virtual machine specific constants
const (
	VMManagementService                         = "Msvm_VirtualSystemManagementService"
	SettingsDefineStateClass                    = "Msvm_SettingsDefineState"               // _SETTINGS_DEFINE_STATE_CLASS
	VirtualSystemSettingDataClass               = "Msvm_VirtualSystemSettingData"          // _VIRTUAL_SYSTEM_SETTING_DATA_CLASS
	ResourceAllocSettingDataClass               = "Msvm_ResourceAllocationSettingData"     // _RESOURCE_ALLOC_SETTING_DATA_CLASS
	ProcessorSettingDataClass                   = "Msvm_ProcessorSettingData"              // _PROCESSOR_SETTING_DATA_CLASS
	MemorySettingDataClass                      = "Msvm_MemorySettingData"                 // _MEMORY_SETTING_DATA_CLASS
	SyntheticEthernetPortSettingDataClass       = "Msvm_SyntheticEthernetPortSettingData"  // _SYNTHETIC_ETHERNET_PORT_SETTING_DATA_CLASS
	EmulatedEthernetPortSettingDataClass        = "Msvm_EmulatedEthernetPortSettingData"   // _EMULATED_ETHERNET_PORT_SETTING_DATA_CLASS
	AffectedJobElementClass                     = "Msvm_AffectedJobElement"                // _AFFECTED_JOB_ELEMENT_CLASS
	ShutdownComponentClass                      = "Msvm_ShutdownComponent"                 // _SHUTDOWN_COMPONENT
	StorageAllocSettingDataClass                = "Msvm_StorageAllocationSettingData"      // _STORAGE_ALLOC_SETTING_DATA_CLASS
	EthernetPortAllocationSettingDataClass      = "Msvm_EthernetPortAllocationSettingData" // _ETHERNET_PORT_ALLOCATION_SETTING_DATA_CLASS
	SerialPortSettingDataClass                  = "Msvm_SerialPortSettingData"             // _TH_SERIAL_PORT_SETTING_DATA_CLASS
	ComputerSystemClass                         = "Msvm_ComputerSystem"
	VirtualHardDiskSettingDataClass             = "Msvm_VirtualHardDiskSettingData"
	ImageManagementServiceClass                 = "Msvm_ImageManagementService"
	VirtualEthernetSwitchManagementServiceClass = "Msvm_VirtualEthernetSwitchManagementService"
)

type ResourceType = int

const (
	DiskDriverType ResourceType = 17
)

// Device resource and subtypes
const (
	PhysDiskResSubType        = "Microsoft:Hyper-V:Physical Disk Drive"       // _PHYS_DISK_RES_SUB_TYPE
	DiskResSubtype            = "Microsoft:Hyper-V:Synthetic Disk Drive"      // _DISK_RES_SUB_TYPE
	DVDResSubType             = "Microsoft:Hyper-V:Synthetic DVD Drive"       // _DVD_RES_SUB_TYPE
	SCSIResSubType            = "Microsoft:Hyper-V:Synthetic SCSI Controller" // _SCSI_RES_SUBTYPE
	IDEDiskResSubType         = "Microsoft:Hyper-V:Virtual Hard Disk"         // _IDE_DISK_RES_SUB_TYPE
	IDEDVDResSubType          = "Microsoft:Hyper-V:Virtual CD/DVD Disk"       // _IDE_DVD_RES_SUB_TYPE
	IDEControllerResSubType   = "Microsoft:Hyper-V:Emulated IDE Controller"   // _IDE_CTRL_RES_SUB_TYPE
	SCSIControllerResSubType  = "Microsoft:Hyper-V:Synthetic SCSI Controller" // _SCSI_CTRL_RES_SUB_TYPE
	VFDDriveResSubType        = "Microsoft:Hyper-V:Synthetic Diskette Drive"  // _VFD_DRIVE_RES_SUB_TYPE
	VFDDiskResSubType         = "Microsoft:Hyper-V:Virtual Floppy Disk"       // _VFD_DISK_RES_SUB_TYPE
	SerialPortResSubType      = "Microsoft:Hyper-V:Serial Port"               // _SERIAL_PORT_RES_SUB_TYPE
	VirtualSystemTypeRealized = "Microsoft:Hyper-V:System:Realized"           // _VIRTUAL_SYSTEM_TYPE_REALIZED
)

// DriveType represents the type of drive that can get attached
// to an IDE or SCSI controller
type DriveType string

// Disk drive types
var (
	DiskDrive DriveType = DiskResSubtype
	DVDDrive  DriveType = DVDResSubType
)

// GenerationType defines the VM generation type
type GenerationType string

// VM generation constants
var (
	Generation1 GenerationType = "Microsoft:Hyper-V:SubType:1"
	Generation2 GenerationType = "Microsoft:Hyper-V:SubType:2"
)
