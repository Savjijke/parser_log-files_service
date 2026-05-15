package core

import (
	"context"
	"errors"
	"log/slog"
)

const locationForLogger = "core/"

type Service struct {
	db     DB
	log    *slog.Logger
	parser Parser
}

func NewService(db DB, parser Parser, log *slog.Logger) *Service {
	return &Service{db: db, parser: parser, log: log}
}

func (s *Service) Parse(ctx context.Context, path string, logID int) error {
	log := s.log.With(
		slog.String("location", locationForLogger+"Parse"),
	)

	pl, err := s.parser.Parse(path, logID)
	if err != nil {
		log.Error("failed parse", "err", err)
		err = s.db.UpdateLogStatus(ctx, logID, constStatusFailed)
		if err != nil {
			log.Error("failed update status", "err", err)
		}
		return err
	}

	err = s.db.SaveParsedLog(ctx, logID, constStatusReady, pl.Nodes, pl.Ports, pl.Settings)
	if err != nil {
		log.Error("failed save parsed log", "err", err)
	}

	return nil
}

func (s *Service) CreateID(ctx context.Context) (int, error) {
	log := s.log.With(
		slog.String("location", locationForLogger+"CreateID"),
	)

	id, err := s.db.CreateLog(ctx, constStatusParsing)
	if err != nil {
		log.Error("failed create log", "err", err)
		return 0, err
	}
	return id, nil
}

func (s *Service) StatsFileLog(ctx context.Context, id int) (FileLog, error) {
	log := s.log.With(
		slog.String("location", locationForLogger+"StatsFileLog"),
	)

	fileLog, err := s.db.GetFileLog(ctx, id)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			log.Warn("not found file log")
			return FileLog{}, ErrNotFound
		}
		log.Error("failed get gile log", "err", err)
	}
	return fileLog, nil
}

func (s *Service) GetDetailsNode(ctx context.Context, logID int, nodeGUID string) (Node, error) {
	log := s.log.With(
		slog.String("location", locationForLogger+"GetDetailsNode"),
	)
	node, err := s.db.GetNode(ctx, logID, nodeGUID)
	if err != nil {
		log.Error("failed get details node", "err", err)
		return Node{}, err
	}
	return node, nil
}

func (s *Service) GetPorts(ctx context.Context, logID int, nodeGUID string) ([]string, error) {
	log := s.log.With(
		slog.String("location", locationForLogger+"GetPorts"),
	)
	ports, err := s.db.GetPortGUIDsByNode(ctx, nodeGUID)
	if err != nil {
		log.Error("failed get portGUIDs", "err", err)
		return nil, err
	}
	return ports, nil
}

func (s *Service) GetTopology(ctx context.Context, logID int) (Topology, error) {
	log := s.log.With(
		slog.String("location", "service.GetTopology"),
		slog.Int("log_id", logID),
	)

	nodes, err := s.db.GetNodesByLogID(ctx, logID)
	if err != nil {
		log.Error("failed get nodes", "err", err)
		return Topology{}, err
	}

	ports, err := s.db.GetPortsByLogID(ctx, logID)
	if err != nil {
		log.Error("failed get ports", "err", err)
		return Topology{}, err
	}

	settings, err := s.db.GetSettingsByLogID(ctx, logID)
	if err != nil {
		log.Error("failed get settings", "err", err)
		return Topology{}, err
	}

	portMap := make(map[string][]Port, len(nodes))

	for _, p := range ports {
		portMap[p.NodeGUID] = append(portMap[p.NodeGUID], p)
	}

	topologyNodes := make([]TopologyNode, 0, len(nodes))

	for _, n := range nodes {
		topologyNodes = append(topologyNodes, TopologyNode{
			Node:  n,
			Ports: portMap[n.NodeGUID],
		})
	}

	result := Topology{
		Nodes:    topologyNodes,
		Settings: settings,
	}

	return result, nil
}
