package rest

import (
	"time"

	"github.com/savjijke/parser-log-files-service/internal/core"
)

func toNodeDTO(logID int, n core.Node) NodeDTO {
	return NodeDTO{
		LogID:        logID,
		NodeGUID:     n.NodeGUID,
		NodeDesc:     n.NodeDesc,
		NodeType:     n.NodeType,
		NumPorts:     n.NumPorts,
		SerialNumber: n.SerialNumber,
		ProductName:  n.ProductName,
	}
}

func toFileLogDTO(l core.FileLog) FileLogDTO {
	return FileLogDTO{
		ID:         l.ID,
		Status:     l.Status,
		CreatedAt:  l.CreatedAt.Format(time.RFC3339),
		NodesCount: l.NodesCount,
		PortsCount: l.PortsCount,
	}
}

func mapTopologyToDTO(t core.Topology) TopologyDTO {
	nodes := make([]TopologyNodeDTO, 0, len(t.Nodes))

	for _, n := range t.Nodes {
		nodes = append(nodes, mapTopologyNodeToDTO(n))
	}

	settings := make([]SwitchSettingsDTO, 0, len(t.Settings))
	for _, s := range t.Settings {
		settings = append(settings, mapSwitchSettingsToDTO(s))
	}

	return TopologyDTO{
		Nodes:    nodes,
		Settings: settings,
	}
}

func mapTopologyNodeToDTO(n core.TopologyNode) TopologyNodeDTO {
	return TopologyNodeDTO{
		Node:  mapNodeToDTO(n.Node),
		Ports: mapPortsToDTO(n.Ports),
	}
}

func mapNodeToDTO(n core.Node) NodeDTO {
	return NodeDTO{
		LogID:        n.LogID,
		NodeGUID:     n.NodeGUID,
		NodeDesc:     n.NodeDesc,
		NodeType:     n.NodeType,
		NumPorts:     n.NumPorts,
		SerialNumber: n.SerialNumber,
		ProductName:  n.ProductName,
	}
}

func mapPortsToDTO(ports []core.Port) []PortDTO {
	result := make([]PortDTO, 0, len(ports))

	for _, p := range ports {
		result = append(result, PortDTO{
			LogID:         p.LogID,
			NodeGUID:      p.NodeGUID,
			PortGUID:      p.PortGUID,
			PortNum:       p.PortNum,
			PortState:     p.PortState,
			PortPhyState:  p.PortPhyState,
			LinkSpeedActv: p.LinkSpeedActv,
			LinkWidthActv: p.LinkWidthActv,
		})
	}

	return result
}

func mapSwitchSettingsToDTO(s core.SwitchSettings) SwitchSettingsDTO {
	return SwitchSettingsDTO{
		LogID:                  s.LogID,
		NodeGUID:               s.NodeGUID,
		Endianness:             s.Endianness,
		EnableEndiannessPerJob: s.EnableEndiannessPerJob,
		ReproducibilityDisable: s.ReproducibilityDisable,
	}
}
