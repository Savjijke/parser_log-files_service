package core

import "time"

const (
	constStatusParsing = "parsing"
	constStatusReady   = "ready"
	constStatusFailed  = "failed"
)

type FileLog struct {
	ID         int
	Status     string
	CreatedAt  time.Time
	NodesCount int
	PortsCount int
}

type Node struct {
	//START_NODES
	LogID    int
	NodeGUID string
	NodeDesc string
	NodeType int
	NumPorts int
	//START_SYSTEM_GENERAL_INFORMATION
	SerialNumber string
	ProductName  string
}

//START_PORTS
type Port struct {
	LogID         int
	NodeGUID      string
	PortGUID      string
	PortNum       int
	PortState     int
	PortPhyState  int
	LinkSpeedActv int
	LinkWidthActv int
}

// file with extantions .sharp_an_info
type SwitchSettings struct {
	LogID                  int
	NodeGUID               string
	Endianness             int
	EnableEndiannessPerJob int
	ReproducibilityDisable int
}

type Payload struct {
	Nodes    []Node
	Ports    []Port
	Settings []SwitchSettings
}

type Topology struct {
	Groups []TopologyGroup
	Edges  []TopologyEdge
}

type TopologyGroup struct {
	Type  string
	Nodes []TopologyNode
}

type TopologyNode struct {
	Node     Node
	Ports    []Port
	Settings *SwitchSettings
}

type TopologyEdge struct {
	FromPort string
	ToPort   string
}
