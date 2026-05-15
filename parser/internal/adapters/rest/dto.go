package rest

type RequestParse struct {
	Path string `json:"path"`
}

type NodeDTO struct {
	LogID        int    `json:"log_id"`
	NodeGUID     string `json:"node_guid"`
	NodeDesc     string `json:"node_desc"`
	NodeType     int    `json:"node_type"`
	NumPorts     int    `json:"num_ports"`
	SerialNumber string `json:"serial_number"`
	ProductName  string `json:"product_name"`
}

type PortDTO struct {
	LogID         int    `json:"logId"`
	NodeGUID      string `json:"nodeGuid"`
	PortGUID      string `json:"portGuid"`
	PortNum       int    `json:"portNum"`
	PortState     int    `json:"portState"`
	PortPhyState  int    `json:"portPhyState"`
	LinkSpeedActv int    `json:"linkSpeedActv"`
	LinkWidthActv int    `json:"linkWidthActv"`
}

type SwitchSettingsDTO struct {
	LogID                  int    `json:"logId"`
	NodeGUID               string `json:"nodeGuid"`
	Endianness             int    `json:"endianness"`
	EnableEndiannessPerJob int    `json:"enableEndiannessPerJob"`
	ReproducibilityDisable int    `json:"reproducibilityDisable"`
}

type FileLogDTO struct {
	ID         int    `json:"id"`
	Status     string `json:"status"`
	CreatedAt  string `json:"created_at"`
	NodesCount int    `json:"nodes_count"`
	PortsCount int    `json:"ports_count"`
}

type TopologyDTO struct {
	Nodes    []TopologyNodeDTO   `json:"nodes"`
	Settings []SwitchSettingsDTO `json:"settings"`
}
type TopologyNodeDTO struct {
	Node  NodeDTO   `json:"node"`
	Ports []PortDTO `json:"ports"`
}

