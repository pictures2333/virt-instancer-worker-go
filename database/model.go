package database

import "time"

// Instance (VM)
type Instance struct {
	ID uint `gorm:"primaryKey;not null;unique;autoIncrement"`

	// from master
	VMUUID string `gorm:"not null;unique"`
	//VPNID  uint   // bind to a VPN client

	// bridge network
	NetworkUUID       string `gorm:"not null;unique"`
	NetworkBridgeName string `gorm:"not null;unique"`
	Subnet            string `gorm:"not null;unique"`

	// vm network
	MacAddress int  `gorm:"not null;unique"`
	Ports      Port `gorm:"foreignKey:VMUUID;references:VMUUID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`

	// data - schema.InstanceData
	Data string `gorm:"not null"`
}

type Port struct {
	ID uint `gorm:"primaryKey;not null;unique;autoIncrement"`

	VMUUID string `gorm:"not null"`
	Port   int    `gorm:"not null;unique"`
}

// synchelper

type FileLink struct {
	ID uint `gorm:"primaryKey;not null;unique;autoIncrement"`

	// fulename of file on MinIO
	// structure : bucket_folder1_folder2_name
	Filename string `gorm:"not null;unique"`

	// OnDelete:SET NULL -> set fileobj to orphan (unlink)
	FileObj FileObj `gorm:"foreignKey:Filename;references:Filename;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
}

type FileObj struct {
	ID uint `gorm:"primaryKey;not null;unique;autoIncrement"`

	Filename string

	// filename of file on disk
	FilenameReal string `gorm:"not null;unique"`

	// from MinIO "last_modified" (for version controlling)
	LastModified time.Time `gorm:"not null"`

	// Status: CREATING,READY
	Status string `gorm:"not null"`

	// OnDelete:RESTRICT -> Cannot delete a FileObj which has Placeholders
	Placeholders []Placeholder `gorm:"foreignKey:FilenameReal;references:FilenameReal;constraint:OnDelete:RESTRICT"`
}

type Placeholder struct {
	ID           uint `gorm:"primaryKey;not null;unique;autoIncrement"`
	FilenameReal string
}

// VPN client

//type VPNClient struct {
//	ID         uint   `gorm:"primaryKey;not null;unique;autoIncrement"`
//	UserID     uint   `gorm:"not null;unique;autoIncrement"` // from master
//	AllowedIPs string `gorm:"not null;unique"`               // wireguard client ip
//
//	// OnDelete:RESTRICT -> Cannot delete a VPNClient which has Instances
//	Instances []Instance `gorm:"foreignKey:VPNID;references:ID;constraint:OnDelete:RESTRICT"`
//}
