package csvparser

import (
	"archive/zip"
	"bufio"
	"encoding/csv"
	"log/slog"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/savjijke/parser-log-files-service/internal/core"
)

const locationForLogger = "adapters/csvParser/"

type Parser struct {
	log *slog.Logger
}

func NewParser(log *slog.Logger) *Parser {
	return &Parser{log: log}
}

func (p *Parser) Parse(path string, id int) (core.Payload, error) {
	log := p.log.With(
		slog.String("location", locationForLogger+"Parse"),
	)

	reader, err := zip.OpenReader(path)
	if err != nil {
		log.Error("failed open reader", "err", err)
		return core.Payload{}, err
	}
	defer func() {
		if err := reader.Close(); err != nil {
			log.Error("failed to reader close")
		}
	}()

	var payload core.Payload

	for _, file := range reader.File {
		ext := filepath.Ext(file.Name)

		switch ext {
		case ".db_csv":
			log.Debug("start read", "file", file.Name, "ext", ".db_csv")
			nodes, ports, systemsInfo, err := parseDBCSV(file, log)
			if err != nil {
				return core.Payload{}, err
			}
			coreNodes := nodesToDomain(nodes, systemsInfo, id)
			corePorts := portsToDomain(ports, id)
			payload.Nodes = coreNodes
			payload.Ports = corePorts
		case ".sharp_an_info":
			log.Debug("start read", "file", file.Name, "ext", ".sharp_an_info")
			settings, err := parseSharpInfo(file, log)
			if err != nil {
				return core.Payload{}, err
			}
			coreSettings := switchSettingsToDomain(settings, id)
			payload.Settings = coreSettings
		default:
			log.Warn("this ext is not supported")
		}

	}

	return payload, nil
}

func parseDBCSV(file *zip.File, log *slog.Logger) (
	[]Node,
	[]Port,
	[]SystemGeneralInformation,
	error,
) {

	rc, err := file.Open()
	if err != nil {
		log.Error("failed open file", "err", err)
		return nil, nil, nil, err
	}
	defer func() {
		if err := rc.Close(); err != nil {
			log.Error("failed to reader close")
		}
	}()

	scanner := bufio.NewScanner(rc)

	buf := make([]byte, 0, 1024*1024)
	scanner.Buffer(buf, 10*1024*1024)

	var (
		nodes   []Node
		ports   []Port
		sysInfo []SystemGeneralInformation
	)

	type sectionType int

	const (
		sectionNone sectionType = iota
		sectionNodes
		sectionPorts
		sectionSystemInfo
	)

	currentSection := sectionNone
	headerSkipped := false

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if line == "" {
			continue
		}

		switch line {

		case "START_NODES":
			currentSection = sectionNodes
			headerSkipped = false
			continue
		case "END_NODES":
			currentSection = sectionNone
			continue

		case "START_PORTS":
			currentSection = sectionPorts
			headerSkipped = false
			continue
		case "END_PORTS":
			currentSection = sectionNone
			continue

		case "START_SYSTEM_GENERAL_INFORMATION":
			currentSection = sectionSystemInfo
			headerSkipped = false
			continue
		case "END_SYSTEM_GENERAL_INFORMATION":
			currentSection = sectionNone
			continue
		}

		if !headerSkipped {
			headerSkipped = true
			continue
		}

		r := csv.NewReader(strings.NewReader(line))
		record, err := r.Read()
		if err != nil {
			log.Error("csv parse error", "err", err)
			return nil, nil, nil, err
		}

		switch currentSection {

		case sectionNodes:
			nodes = append(nodes, Node{
				NodeDesc:        record[0],
				NumPorts:        atoi(record[1]),
				NodeType:        atoi(record[2]),
				ClassVersion:    atoi(record[3]),
				BaseVersion:     atoi(record[4]),
				SystemImageGUID: record[5],
				NodeGUID:        record[6],
				PortGUID:        record[7],
			})

		case sectionPorts:
			ports = append(ports, Port{
				NodeGUID:      record[0],
				PortGUID:      record[1],
				PortNum:       atoi(record[2]),
				PortState:     atoi(record[19]),
				PortPhyState:  atoi(record[20]),
				LinkSpeedActv: atoi(record[15]),
				LinkWidthActv: atoi(record[10]),
			})

		case sectionSystemInfo:
			sysInfo = append(sysInfo, SystemGeneralInformation{
				NodeGuid:     record[0],
				SerialNumber: record[1],
				PartNumber:   record[2],
				Revision:     record[3],
				ProductName:  record[4],
			})
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, nil, nil, err
	}

	log.Info("parsed db_csv",
		"nodes", len(nodes),
		"ports", len(ports),
		"system_info", len(sysInfo),
	)

	return nodes, ports, sysInfo, nil
}

func atoi(s string) int {
	v, _ := strconv.Atoi(s)

	return v
}

func parseSharpInfo(file *zip.File, log *slog.Logger) ([]SwitchSettings, error) {
	rc, err := file.Open()
	if err != nil {
		log.Error("failed open file", "err", err)
		return nil, err
	}
	defer func() {
		if err := rc.Close(); err != nil {
			log.Error("failed to reader close")
		}
	}()
	log.Debug("start parse")

	scanner := bufio.NewScanner(rc)
	var result []SwitchSettings
	var current *SwitchSettings

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if line == "" || strings.HasPrefix(line, "----") {
			continue
		}

		if strings.HasPrefix(line, "SW_GUID=") {
			if current != nil {
				result = append(result, *current)
			}

			current = &SwitchSettings{
				NodeGUID: strings.TrimPrefix(line, "SW_GUID="),
			}

			continue
		}

		if current == nil {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		valueStr := strings.TrimSpace(parts[1])

		switch key {
		case "endianness":
			v, _ := strconv.Atoi(valueStr)
			current.Endianness = v

		case "enable_endianness_per_job":
			v, _ := strconv.Atoi(valueStr)
			current.EnableEndiannessPerJob = v

		case "reproducibility_disable":
			v, _ := strconv.Atoi(valueStr)
			current.ReproducibilityDisable = v
		}
	}

	if current != nil {
		result = append(result, *current)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	log.Debug("parse finished", "len", len(result))

	return result, nil
}
