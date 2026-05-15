package rest

import (
	"time"

	"github.com/savjijke/parser-log-files-service/internal/core"
)

func toNodeDTO(logID int, n core.Node) ReplyNode {
	return ReplyNode{
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

func mapNodeToDTO(n core.Node) NodeDTO {
	return NodeDTO{
		NodeGUID:     n.NodeGUID,
		NodeDesc:     n.NodeDesc,
		NodeType:     mapNodeType(n.NodeType),
		NumPorts:     n.NumPorts,
		SerialNumber: n.SerialNumber,
		ProductName:  n.ProductName,
	}
}

func mapPortsToDTO(ports []core.Port) []PortDTO {
	result := make([]PortDTO, 0, len(ports))

	for _, p := range ports {
		result = append(result, PortDTO{
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
		Endianness:             s.Endianness,
		EnableEndiannessPerJob: s.EnableEndiannessPerJob,
		ReproducibilityDisable: s.ReproducibilityDisable,
	}
}

func mapTopologyToDTO(t core.Topology) TopologyDTO {
	groups := make([]TopologyGroupDTO, 0, len(t.Groups))

	for _, g := range t.Groups {
		groups = append(groups, mapGroupToDTO(g))
	}

	return TopologyDTO{
		Groups: groups,
	}
}

func mapGroupToDTO(g core.TopologyGroup) TopologyGroupDTO {
	nodes := make([]TopologyNodeDTO, 0, len(g.Nodes))

	for _, n := range g.Nodes {
		nodes = append(nodes, mapTopologyNodeToDTO(n))
	}

	return TopologyGroupDTO{
		Type:  g.Type,
		Nodes: nodes,
	}
}

func mapTopologyNodeToDTO(n core.TopologyNode) TopologyNodeDTO {
	var settingsDTO *SwitchSettingsDTO

	if n.Settings != nil {
		s := mapSwitchSettingsToDTO(*n.Settings)
		settingsDTO = &s
	}

	return TopologyNodeDTO{
		Node:     mapNodeToDTO(n.Node),
		Ports:    mapPortsToDTO(n.Ports),
		Settings: settingsDTO,
	}
}

func mapNodeType(t int) string {
	switch t {
	case 1:
		return "host"
	default:
		return "switch"
	}
}
