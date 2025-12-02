package vmhelper

import (
	"Instancer-worker-go/config"
	"Instancer-worker-go/database"
	"Instancer-worker-go/libvirtHelper"
	"Instancer-worker-go/schema"
	"Instancer-worker-go/utils"
	"fmt"
	"iter"
	"log"
	"net"
	"sync"

	"github.com/google/uuid"
)

var once sync.Once

func Init() {
	once.Do(func() {
		var err error

		// mkdir
		if err = utils.CheckDir(config.FileDir); err != nil {
			log.Fatalf("Error while checking directory %s : %v", config.FileDir, err)
		}

		fileVmfilesDir := config.FileDir + "/vmfiles"
		if err = VMMustMkdir(fileVmfilesDir); err != nil {
			log.Fatalf("Failed to create directory %s : %v", fileVmfilesDir, err)
		}

		// generate allSubnets
		if err = initializeAllSubnets(config.BaseNetwork, config.SubnetPrefix); err != nil {
			log.Fatalf("Failed to generate all subnets for vmhelper : %v", err)
		}
	})
}

func Create(data *schema.InstanceData) (err error) {
	var ok bool

	// lock
	if err = vmlock(data.VMUUID); err != nil {
		return err
	}
	defer vmunlock(data.VMUUID)

	// check exists (instance)
	var instances []database.Instance
	if instances, err = database.ReadInstance(&data.VMUUID); err != nil {
		return fmt.Errorf("failed to read database : %v", err)
	}
	if len(instances) != 0 {
		return fmt.Errorf("Instance %s has been created", data.VMUUID)
	}

	// network resources

	// - bridge
	var subnet *net.IPNet
	if subnet, err = GetSubnetAvaliable(); err != nil {
		return err
	}
	subnetNext, subnetStop := iter.Pull(EnumerateIPinSubnet(subnet))
	defer subnetStop()

	networkUUID := uuid.New().String()
	networkBridgeName := fmt.Sprintf("br%s", networkUUID[:8])
	networkIP := subnet.IP.String()
	networkNetmask := fmt.Sprintf("%d.%d.%d.%d",
		subnet.Mask[0],
		subnet.Mask[1],
		subnet.Mask[2],
		subnet.Mask[3],
	)

	// - VM
	var vmMacAddressInt int
	if vmMacAddressInt, err = GetMacAddressAvaliable(); err != nil {
		return err
	}
	vmMacAddress := MAC_Int2Str(vmMacAddressInt)
	var vmIP string
	if vmIP, ok = subnetNext(); !ok {
		return fmt.Errorf("No IP avaliable in subnet")
	}

	// - VM.VNC
	var (
		vncPort   int
		vncPasswd string
	)
	if data.Config.VNC {
		if vncPort, err = GetPortAvaliable(); err != nil {
			return err
		}
		vncPasswd = uuid.New().String()[:8]
	}

	// prepare data
	data.Network = &schema.Network{
		UUID:    networkUUID,
		Brname:  networkBridgeName,
		IP:      networkIP,
		Netmask: networkNetmask,
		VM: schema.VMNetwork{
			MacAddress: vmMacAddress,
			IP:         vmIP,
		},
	}

	if data.Config.VNC {
		data.Network.VM.VNC = &schema.VNC{
			IP:     config.VNCListenHost,
			Port:   vncPort,
			Passwd: vncPasswd,
		}
	}

	// write db
	if err = database.CreateInstance(
		data.VMUUID,
		networkUUID, networkBridgeName, networkIP,
		vmMacAddressInt, data,
	); err != nil {
		return err
	}
	defer func(err *error) {
		if *err != nil {
			if err := database.DeleteInstance(data.VMUUID); err != nil {
				utils.Showerr(fmt.Sprintf("failed to delete Instance from database : %v", err), true)
			}
		}
	}(&err)

	if data.Config.VNC {
		if err = database.AllocatePort(data.VMUUID, vncPort); err != nil {
			return err
		}
	}

	// mkdir
	dir := fmt.Sprintf("%s/vmfiles/%s", config.FileDir, data.VMUUID)
	if err = VMMustMkdir(dir); err != nil {
		return fmt.Errorf("Failed to create directory %s : %v", dir, err)
	}
	defer func(err *error) {
		if *err != nil {
			if err := VMMustRmdir(dir); err != nil {
				utils.Showerr(fmt.Sprintf("failed to remove directory %s : %v", dir, err), true)
			}
		}
	}(&err)

	// copy file
	for _, file := range data.Config.Files {
		if err = VMCopyFile(data.VMUUID, &file); err != nil {
			return err
		}
	}

	// render xml file
	var (
		vmxml  string
		netxml string
	)
	if vmxml, err = XMLTemplate(data, "vm"); err != nil {
		return fmt.Errorf("Failed to render XML : %v", err)
	}
	if netxml, err = XMLTemplate(data, "network"); err != nil {
		return fmt.Errorf("Failed to render XML : %v", err)
	}

	// start network and vm
	// set up firewall in callback function
	if err = libvirtHelper.StartNetwork(netxml); err != nil {
		return err
	}
	defer func(err *error) {
		if *err != nil {
			if err := libvirtHelper.DeleteNetwork(data.Network.UUID); err != nil {
				utils.Showerr(fmt.Sprintf("failed to delete network %s : %v", data.Network.UUID, err), true)
			}
		}
	}(&err)

	if err = libvirtHelper.StartVM(vmxml); err != nil {
		return err
	}
	defer func(err *error) {
		if *err != nil {
			if err := libvirtHelper.DeleteVM(data.VMUUID); err != nil {
				utils.Showerr(fmt.Sprintf("failed to delete VM %s : %v", data.VMUUID, err), true)
			}
		}
	}(&err)

	return nil
}

func Delete(VMUUID string) (err error) {
	// lock
	if err = vmlock(VMUUID); err != nil {
		return err
	}
	defer vmunlock(VMUUID)

	// check exists
	var instances []database.Instance
	if instances, err = database.ReadInstance(&VMUUID); err != nil {
		return fmt.Errorf("failed to read VM %s : %v", VMUUID, err)
	}
	if len(instances) != 1 {
		return fmt.Errorf("VM %s not found", VMUUID)
	}
	instance := instances[0]

	// delete vm
	if err = libvirtHelper.DeleteVM(VMUUID); err != nil {
		return fmt.Errorf("failed to delete VM %s : %v", VMUUID, err)
	}

	// delete network
	if err = libvirtHelper.DeleteNetwork(instance.NetworkUUID); err != nil {
		return fmt.Errorf("failed to delete network %s : %v", instance.NetworkUUID, err)
	}

	// delete files
	dir := fmt.Sprintf("%s/vmfiles/%s", config.FileDir, instance.VMUUID)
	if err = VMMustRmdir(dir); err != nil {
		// todo: 考慮要不要忽略這個錯誤
		return fmt.Errorf("failed to remove directory %s : %v", dir, err)
	}

	// delete db record
	// free ports (foreign key CASCADE)
	if err = database.DeleteInstance(VMUUID); err != nil {
		return fmt.Errorf("failed to delete Instance from database : %v", err)
	}

	return nil
}
