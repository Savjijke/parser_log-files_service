package core

import (
	"context"
)

type DB interface {
	CreateLog(
		ctx context.Context,
		status string,
	) (int, error)
	GetFileLog(
		ctx context.Context,
		id int,
	) (FileLog, error)
	SaveParsedLog(
		ctx context.Context,
		logID int,
		status string,
		nodes []Node,
		ports []Port,
		settings []SwitchSettings,
	) error
	UpdateLogStatus(
		ctx context.Context,
		logID int,
		status string,
	) error
	GetNode(
		ctx context.Context,
		logID int,
		nodeGUID string,
	) (Node, error)
	GetPortGUIDsByNode(
		ctx context.Context,
		nodeGUID string,
	) ([]string, error)
	GetNodesByLogID(
		ctx context.Context,
		logID int,
	) ([]Node, error)
	GetPortsByLogID(
		ctx context.Context,
		logID int,
	) ([]Port, error)
	GetSettingsByLogID(
		ctx context.Context,
		logID int,
	) ([]SwitchSettings, error)
}

type Parser interface {
	Parse(path string, logID int) (Payload, error)
}

type ParserService interface {
	Parse(ctx context.Context, path string, logID int) error
	CreateID(ctx context.Context) (int, error)
	StatsFileLog(ctx context.Context, id int) (FileLog, error)
	GetDetailsNode(ctx context.Context, logID int, nodeGUID string) (Node, error)
	GetPorts(ctx context.Context, logID int, nodeGUID string) ([]string, error)
	GetTopology(ctx context.Context, logID int) (Topology, error)
}
