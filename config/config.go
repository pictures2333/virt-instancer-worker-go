package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// database (sqlite)
var Database string

// minio
var MinIOEndpoint string
var MinIOAccesskey string
var MinIOSecretkey string

// master
var MasterUrl string

// file storage
var FileDir string

// synchelper
var SynchelperSleep int

// vmhelper
var LibvirtGroup string
var QemuUrl string

// vmhelper - network
var MacAddressPrefix string
var BaseNetwork string
var SubnetPrefix int
var PortMin int
var PortMax int

// vmhelper - vnc
var VNCListenHost string

// VPN (wireguard)
var WireguardAddress string // as "base network"
var WireguardHost string    // for clients to connect
var VPNMGRListenHost string // for vpn manager "http" server

func loadenv() {
	var err error

	err = godotenv.Load()
	if err != nil {
		log.Fatal("Failed to load .env file")
	}
}

func Init() {
	var err error

	loadenv()

	Database = os.Getenv("DATABASE")

	MinIOEndpoint = os.Getenv("MINIO_ENDPOINT")
	MinIOAccesskey = os.Getenv("MINIO_ACCESSKEY")
	MinIOSecretkey = os.Getenv("MINIO_SECRETKEY")

	MasterUrl = os.Getenv("MASTER_URL")

	FileDir = os.Getenv("FILE_DIR")

	SynchelperSleep, err = strconv.Atoi(os.Getenv("SYNCHELPER_SLEEP"))
	if err != nil {
		log.Fatal(err)
	}

	LibvirtGroup = os.Getenv("LIBVIRT_GROUP")
	QemuUrl = os.Getenv("QEMU_URL")

	MacAddressPrefix = os.Getenv("MAC_ADDRESS_PREFIX")
	BaseNetwork = os.Getenv("BASE_NETWORK")
	SubnetPrefix, err = strconv.Atoi(os.Getenv("SUBNET_PREFIX"))
	if err != nil {
		log.Fatal(err)
	}
	PortMin, err = strconv.Atoi(os.Getenv("PORT_MIN"))
	if err != nil {
		log.Fatal(err)
	}
	PortMax, err = strconv.Atoi(os.Getenv("PORT_MAX"))
	if err != nil {
		log.Fatal(err)
	}

	VNCListenHost = os.Getenv("VNC_LISTEN_HOST")

	WireguardAddress = os.Getenv("WIREGUARD_ADDRESS")
	WireguardHost = os.Getenv("WIREGUARD_HOST")
	VPNMGRListenHost = os.Getenv("VPN_MGR_LISTEN_HOST")
}
