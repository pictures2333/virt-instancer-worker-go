package vmhelper

import (
	"Instancer-worker-go/config"
	"Instancer-worker-go/database"
	"fmt"
	"net"
	"net/netip"
	"slices"

	"github.com/apparentlymart/go-cidr/cidr"
)

// Mac Address

func GetMacAddressAvaliable() (result int, err error) {
	// query
	var usedMacAddress *[]int
	if usedMacAddress, err = database.ReadInstanceMacAddressAll(); err != nil {
		return -1, err
	}

	// find
	for i := range 2 << 24 {
		if !slices.Contains(*usedMacAddress, i) {
			return i, nil
		}
	}

	return -1, fmt.Errorf("no MAC address avaliable")
}

func MAC_Int2Str(mac int) string {
	return fmt.Sprintf("%s:%02x:%02x:%02x",
		config.MacAddressPrefix,
		(mac>>16)&0xff,
		(mac>>8)&0xff,
		mac&0xff,
	)
}

// IP address (IPv4)

var allSubnets []*net.IPNet

// initialize all subnets
func initializeAllSubnets(baseNetworkStr string, subnetPrefix int) (err error) {
	var baseNetwork *net.IPNet
	if _, baseNetwork, err = net.ParseCIDR(baseNetworkStr); err != nil {
		return err
	}

	// check
	baseNetworkPrefix, _ := baseNetwork.Mask.Size()
	if subnetPrefix <= baseNetworkPrefix {
		return fmt.Errorf("subnet_prefix should be larger than base_network_prefix")
	}
	if subnetPrefix > 30 {
		return fmt.Errorf("subnet_prefix should not be larger than 30")
	}

	// generate
	newBits := subnetPrefix - baseNetworkPrefix
	numSubnets := 1 << newBits
	allSubnets = make([]*net.IPNet, 0, numSubnets)
	for i := 0; i < numSubnets; i++ {
		var subnet *net.IPNet
		if subnet, err = cidr.Subnet(baseNetwork, newBits, i); err != nil {
			return fmt.Errorf("Error while calculating subnet %d : %v", i, err)
		}
		allSubnets = append(allSubnets, subnet)
	}

	return nil
}

func GetSubnetAvaliable() (result *net.IPNet, err error) {
	// query
	var usedSubnets *[]string
	if usedSubnets, err = database.ReadInstanceSubnetAll(); err != nil {
		return nil, err
	}

	// find
	for _, subnet := range allSubnets {
		if !slices.Contains(*usedSubnets, subnet.IP.String()) {
			return subnet, nil
		}
	}
	return nil, fmt.Errorf("no subnet avaliable")
}

func EnumerateIPinSubnet(subnet *net.IPNet) func(yield func(string) bool) {
	return func(yield func(string) bool) {
		addr, ok := netip.AddrFromSlice(subnet.IP)
		if !ok {
			return
		}

		// 取得遮罩長度（例如 /24）
		ones, _ := subnet.Mask.Size()

		// 建立 Prefix 物件
		prefix := netip.PrefixFrom(addr, ones)

		prefix = prefix.Masked()
		startIP := prefix.Addr()
		for ip := startIP; prefix.Contains(ip); ip = ip.Next() {
			// 跳過網段 IP（第一個 IP）
			if ip == startIP {
				continue
			}

			// 跳過廣播 IP（最後一個 IP）
			if !prefix.Contains(ip.Next()) {
				continue
			}

			if !yield(ip.String()) {
				return
			}
		}
	}
}

// Port

func GetPortAvaliable() (result int, err error) {
	// query
	var usedPortsDB []database.Port
	if usedPortsDB, err = database.ReadPort(); err != nil {
		return -1, err
	}
	var usedPort []int
	for _, p := range usedPortsDB {
		usedPort = append(usedPort, p.Port)
	}

	// find
	for p := config.PortMin; p <= config.PortMax; p++ {
		if !slices.Contains(usedPort, p) {
			return p, nil
		}
	}
	return -1, fmt.Errorf("no port avaliable")
}
