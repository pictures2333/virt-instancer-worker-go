package libvirtHelper

import (
	"Instancer-worker-go/config"
	"Instancer-worker-go/utils"
	"fmt"

	"github.com/libvirt/libvirt-go"
)

// nftables
// - isolation
// - vpn

func StartNetwork(xml string) (err error) {
	// connect
	var conn *libvirt.Connect
	if conn, err = libvirt.NewConnect(config.QemuUrl); err != nil {
		return fmt.Errorf("failed to connect libvirt : %v", err)
	}
	defer conn.Close()

	// define network
	var network *libvirt.Network
	if network, err = conn.NetworkDefineXML(xml); err != nil {
		return fmt.Errorf("failed to define network : %v", err)
	}
	defer func(err *error) {
		if *err != nil {
			if err := network.Undefine(); err != nil {
				utils.Showerr(fmt.Sprintf("failed to undefine network : %v", err), true)
			}
		}
	}(&err)

	// create network
	if err = network.Create(); err != nil {
		return fmt.Errorf("failed to create network : %v", err)
	}

	return nil
}

func DeleteNetwork(networkUUID string) (err error) {
	// connect
	var conn *libvirt.Connect
	if conn, err = libvirt.NewConnect(config.QemuUrl); err != nil {
		return fmt.Errorf("failed to connect libvirt : %v", err)
	}
	defer conn.Close()

	// find
	var network *libvirt.Network
	if network, err = conn.LookupNetworkByName(networkUUID); err != nil {
		lverr, ok := err.(libvirt.Error)

		if ok && lverr.Code == libvirt.ERR_NO_NETWORK {
			// not found (as "deleted")
			return nil
		} else {
			// type conversation failed or other error
			return fmt.Errorf("failed to lookup network %s : %v", networkUUID, err)
		}
	}

	// destory
	if err = network.Destroy(); err != nil {
		utils.Showerr(fmt.Sprintf("failed to destory network %s : %v", networkUUID, err), false)
		// ignore error and undefine the network
		err = nil
	}

	// undefine
	if err = network.Undefine(); err != nil {
		return fmt.Errorf("failed to undefine network %s : %v", networkUUID, err)
	}

	return nil
}

func StartVM(xml string) (err error) {
	// connect
	var conn *libvirt.Connect
	if conn, err = libvirt.NewConnect(config.QemuUrl); err != nil {
		return fmt.Errorf("failed to connect libvirt : %v", err)
	}
	defer conn.Close()

	// define
	var domain *libvirt.Domain
	if domain, err = conn.DomainDefineXML(xml); err != nil {
		return fmt.Errorf("failed to define domain : %v", err)
	}
	defer func(err *error) {
		if *err != nil {
			if err := domain.UndefineFlags(libvirt.DOMAIN_UNDEFINE_NVRAM); err != nil {
				utils.Showerr(fmt.Sprintf("failed to undefine domain : %v", err), true)
			}
		}
	}(&err)

	// create domain
	if err = domain.Create(); err != nil {
		return fmt.Errorf("failed to create domain : %v", err)
	}

	return nil
}

func DeleteVM(VMUUID string) (err error) {
	// connect
	var conn *libvirt.Connect
	if conn, err = libvirt.NewConnect(config.QemuUrl); err != nil {
		return fmt.Errorf("failed to connect libvirt : %v", err)
	}
	defer conn.Close()

	// find
	var domain *libvirt.Domain
	if domain, err = conn.LookupDomainByName(VMUUID); err != nil {
		lverr, ok := err.(libvirt.Error)

		if ok && lverr.Code == libvirt.ERR_NO_DOMAIN {
			// not found (as "deleted")
			return nil
		} else {
			// type converation failed or other error
			return fmt.Errorf("failed to lookup domain %s : %v", VMUUID, err)
		}
	}

	// destory
	if err = domain.Destroy(); err != nil {
		utils.Showerr(fmt.Sprintf("failed to destory domain %s : %v", VMUUID, err), false)
		// ignore error and undefine the domain
		err = nil
	}

	// undefine
	if err = domain.UndefineFlags(libvirt.DOMAIN_UNDEFINE_NVRAM); err != nil {
		return fmt.Errorf("failed to undefine domain %s : %v", VMUUID, err)
	}

	return nil
}
