package csvparser

import (
	"strings"

	"github.com/savjijke/parser-log-files-service/internal/core"
)

func nodesToDomain(logID string, nodes []Node, sysInfo []SystemGeneralInformation) []core.Node {
	sysInfoMap := make(map[string]SystemGeneralInformation, len(sysInfo))

	for _, s := range sysInfo {
		sysInfoMap[s.NodeGuid] = s
	}

	out := make([]core.Node, 0, len(nodes))

	for _, n := range nodes {
		info := sysInfoMap[n.NodeGUID]

		out = append(out, core.Node{
			LogID:        logID,
			NodeGUID:     n.NodeGUID,
			NodeDesc:     n.NodeDesc,
			NodeType:     n.NodeType,
			NumPorts:     n.NumPorts,
			SerialNumber: info.SerialNumber,
			ProductName:  info.ProductName,
		})
	}

	return out
}

func portsToDomain(logID string, ports []Port) []core.Port {
	out := make([]core.Port, 0, len(ports))

	for _, p := range ports {
		out = append(out, core.Port{
			LogID:         logID,
			NodeGUID:      p.NodeGUID,
			PortGUID:      p.PortGUID,
			PortNum:       p.PortNum,
			PortState:     p.PortState,
			PortPhyState:  p.PortPhyState,
			LinkSpeedActv: p.LinkSpeedActv,
			LinkWidthActv: p.LinkWidthActv,
		})
	}

	return out
}

func switchSettingsToDomain(logID string, settings []SwitchSettings) []core.SwitchSettings {
	out := make([]core.SwitchSettings, 0, len(settings))

	for _, s := range settings {
		out = append(out, core.SwitchSettings{
			LogID:                  logID,
			NodeGUID:               normalizeSwitchGUID(s.NodeGUID),
			Endianness:             s.Endianness,
			EnableEndiannessPerJob: s.EnableEndiannessPerJob,
			ReproducibilityDisable: s.ReproducibilityDisable,
		})
	}

	return out
}

func normalizeSwitchGUID(guid string) string {
	if strings.HasPrefix(guid, "0x") {
		return guid
	}

	return "0x" + guid
}
