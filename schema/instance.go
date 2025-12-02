package schema

type File struct {
	Bucket   string `json:"Bucket"`
	Filename string `json:"Filename"`
	Type     string `json:"type"`
}

type ChallengeConfig struct {
	Files        map[string]File `json:"files"`       // Files["iso"] = File{...}
	VMXML        string          `json:"vm_xml"`      // raw xml
	NetworkXML   string          `json:"network_xml"` // raw xml
	VNC          bool            `json:"vnc"`
	AliveMinutes int             `json:"alive_minutes"`
}

type VNC struct {
	IP     string `json:"ip"`
	Port   int    `json:"port"`
	Passwd string `json:"passwd"`
}

type VMNetwork struct {
	MacAddress string `json:"mac_address"`
	IP         string `json:"ip"` // VM IP
	VNC        *VNC   `json:"vnc"`
}

type Network struct {
	// bridge
	UUID    string `json:"uuid"`
	Brname  string `json:"brname"`
	IP      string `json:"ip"` // subnet IP
	Netmask string `json:"netmask"`
	// vm
	VM VMNetwork `json:"vm"`
}

// request body from master
type InstanceData struct {
	// from master
	VMUUID string `json:"vm_uuid"`
	UserID uint   `json:"user_id"`

	// from config
	Config ChallengeConfig `json:"config"`

	// from instance
	Network *Network `json:"network"`
}
