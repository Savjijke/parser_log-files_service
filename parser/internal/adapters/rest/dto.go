package rest

type RequestParse struct {
	Path string `json:"path"`
}

type ReplyNode struct {
	LogID        int    `json:"log_id" `
	NodeGUID     string `json:"node_guid"`
	NodeDesc     string `json:"node_desc"`
	NodeType     int    `json:"node_type"`
	NumPorts     int    `json:"num_ports"`
	SerialNumber string `json:"serial_number"`
	ProductName  string `json:"product_name"`
}

type NodeDTO struct {
	NodeGUID     string `json:"node_guid"`
	NodeDesc     string `json:"node_desc"`
	NodeType     string `json:"node_type"`
	NumPorts     int    `json:"num_ports"`
	SerialNumber string `json:"serial_number"`
	ProductName  string `json:"product_name"`
}

type PortDTO struct {
	PortGUID      string `json:"port_guid"`
	PortNum       int    `json:"port_num"`
	PortState     int    `json:"port_state"`
	PortPhyState  int    `json:"port_phy_state"`
	LinkSpeedActv int    `json:"link_speed_actv"`
	LinkWidthActv int    `json:"link_width_actv"`
}

type SwitchSettingsDTO struct {
	Endianness             int `json:"endianness"`
	EnableEndiannessPerJob int `json:"enable_endianness_per_job"`
	ReproducibilityDisable int `json:"reproducibility_disable"`
}

type FileLogDTO struct {
	ID         int    `json:"id"`
	Status     string `json:"status"`
	CreatedAt  string `json:"created_at"`
	NodesCount int    `json:"nodes_count"`
	PortsCount int    `json:"ports_count"`
}

type TopologyGroupDTO struct {
	Type  string            `json:"type"` // host | switch
	Nodes []TopologyNodeDTO `json:"nodes"`
}

type TopologyDTO struct {
	Groups []TopologyGroupDTO `json:"groups"`
}
type TopologyNodeDTO struct {
	Node     NodeDTO            `json:"node"`
	Ports    []PortDTO          `json:"ports"`
	Settings *SwitchSettingsDTO `json:"settings,omitempty"`
}
