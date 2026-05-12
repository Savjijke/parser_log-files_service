package core

import "time"

const (
	constStatusParsing = "parsing"
	constStatusReady   = "ready"
)

type Log struct {
	ID         string
	Status     string
	CreatedAt  time.Time
	NodesCount int
	PortsCount int
}

type Node struct {
	//START_NODES
	LogID    string
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
	LogID         string
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
	LogID                  string
	NodeGUID               string
	Endianness             int
	EnableEndiannessPerJob int
	ReproducibilityDisable int
}
