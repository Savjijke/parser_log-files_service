package db

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"

	"github.com/jmoiron/sqlx"
	"github.com/savjijke/parser-log-files-service/internal/core"
)

const locationForLogger = "adapters/db/"

type DB struct {
	log  *slog.Logger
	conn *sqlx.DB
}

func NewDB(dsn string, log *slog.Logger) (*DB, error) {

	db, err := sqlx.Connect("pgx", dsn)
	if err != nil {
		log.Error("connection problem", "address", dsn, "error", err)
		return nil, err
	}
	log.Debug("DB connected", "db dsn", dsn)
	return &DB{log: log, conn: db}, nil
}

func (db *DB) CreateLog(ctx context.Context, status string) (int, error) {
	log := db.log.With(
		slog.String("location", locationForLogger+"CreateLog"),
	)
	const query = `
		INSERT INTO file_logs (
			status,
			created_at,
			nodes_count,
			ports_count
		)
		VALUES ($1, NOW(), 0, 0)
		RETURNING id;
	`

	var id int64

	err := db.conn.QueryRowContext(
		ctx,
		query,
		status,
	).Scan(&id)

	if err != nil {
		log.Error("failed query row", "err", err)
		return 0, err
	}

	return int(id), nil
}

func (db *DB) GetFileLog(ctx context.Context, id int) (core.FileLog, error) {
	log := db.log.With(
		slog.String("location", locationForLogger+"GetFileLog"),
	)
	const query = `
		SELECT
			id,
			status,
			created_at,
			nodes_count,
			ports_count
		FROM file_logs
		WHERE id = $1;
	`

	var logFile core.FileLog

	err := db.conn.QueryRowContext(
		ctx,
		query,
		id,
	).Scan(
		&logFile.ID,
		&logFile.Status,
		&logFile.CreatedAt,
		&logFile.NodesCount,
		&logFile.PortsCount,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Warn("not found")
			return core.FileLog{}, core.ErrNotFound
		}
		log.Error("failed query row", "err", err)
		return core.FileLog{}, err
	}

	return logFile, nil
}

func (db *DB) UpdateFileLog(
	ctx context.Context,
	tx *sql.Tx,
	logID int,
	status string,
) error {
	log := db.log.With(
		slog.String("location", locationForLogger+"UpdateFileLog"),
	)

	const query = `
		UPDATE file_logs
		SET
			status = $2,
			nodes_count = (
				SELECT COUNT(*)
				FROM nodes
				WHERE log_id = $1
			),
			ports_count = (
				SELECT COUNT(*)
				FROM ports
				WHERE log_id = $1
			)
		WHERE id = $1;
	`

	_, err := tx.ExecContext(
		ctx,
		query,
		logID,
		status,
	)

	if err != nil {
		log.Error("failed update file_logs", "err", err)
		return err
	}

	return nil
}

func (db *DB) UpsertNode(ctx context.Context, nodes []core.Node) error {
	log := db.log.With(
		slog.String("location", locationForLogger+"UpdateNode"),
	)

	if len(nodes) == 0 {
		return nil
	}

	const query = `
		INSERT INTO nodes (
			log_id,
			node_guid,
			node_desc,
			node_type,
			num_ports,
			serial_number,
			product_name
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (log_id, node_guid)
		DO UPDATE SET
			node_desc = EXCLUDED.node_desc,
			node_type = EXCLUDED.node_type,
			num_ports = EXCLUDED.num_ports,
			serial_number = EXCLUDED.serial_number,
			product_name = EXCLUDED.product_name;
	`

	tx, err := db.conn.BeginTx(ctx, nil)
	if err != nil {
		log.Error("failed begin tx", "err", err)
		return err
	}

	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		log.Error("failed prepare stmt", "err", err)
		return err
	}

	defer stmt.Close()

	for _, node := range nodes {
		_, err := stmt.ExecContext(
			ctx,
			node.LogID,
			node.NodeGUID,
			node.NodeDesc,
			node.NodeType,
			node.NumPorts,
			node.SerialNumber,
			node.ProductName,
		)

		if err != nil {
			log.Error("failed exec", "err", err)
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		log.Error("failed commit", "err", err)
		return err
	}

	return nil
}

func (db *DB) UpsertPort(ctx context.Context, ports []core.Port) error {
	log := db.log.With(
		slog.String("location", locationForLogger+"UpdatePort"),
	)

	if len(ports) == 0 {
		return nil
	}

	const query = `
		INSERT INTO ports (
			log_id,
			node_guid,
			port_guid,
			port_num,
			port_state,
			port_phy_state,
			link_speed_actv,
			link_width_actv
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		ON CONFLICT (log_id, node_guid, port_guid)
		DO UPDATE SET
			port_num = EXCLUDED.port_num,
			port_state = EXCLUDED.port_state,
			port_phy_state = EXCLUDED.port_phy_state,
			link_speed_actv = EXCLUDED.link_speed_actv,
			link_width_actv = EXCLUDED.link_width_actv;
	`

	tx, err := db.conn.BeginTx(ctx, nil)
	if err != nil {
		log.Error("failed begin tx", "err", err)
		return err
	}

	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		log.Error("failed prepare stmt", "err", err)
		return err
	}

	defer stmt.Close()

	for _, port := range ports {
		_, err := stmt.ExecContext(
			ctx,
			port.LogID,
			port.NodeGUID,
			port.PortGUID,
			port.PortNum,
			port.PortState,
			port.PortPhyState,
			port.LinkSpeedActv,
			port.LinkWidthActv,
		)

		if err != nil {
			log.Error("failed exec", "err", err)
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		log.Error("failed commit", "err", err)
		return err
	}

	return nil
}

func (db *DB) UpsertSettings(ctx context.Context, sets []core.SwitchSettings) error {
	log := db.log.With(
		slog.String("location", locationForLogger+"UpdateSettings"),
	)

	if len(sets) == 0 {
		return nil
	}

	const query = `
		INSERT INTO switch_settings (
			log_id,
			node_guid,
			endianness,
			enable_endianness_per_job,
			reproducibility_disable
		)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (log_id, node_guid)
		DO UPDATE SET
			endianness = EXCLUDED.endianness,
			enable_endianness_per_job = EXCLUDED.enable_endianness_per_job,
			reproducibility_disable = EXCLUDED.reproducibility_disable;
	`

	tx, err := db.conn.BeginTx(ctx, nil)
	if err != nil {
		log.Error("failed begin tx", "err", err)
		return err
	}

	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		log.Error("failed prepare stmt", "err", err)
		return err
	}

	defer stmt.Close()

	for _, set := range sets {
		_, err := stmt.ExecContext(
			ctx,
			set.LogID,
			set.NodeGUID,
			set.Endianness,
			set.EnableEndiannessPerJob,
			set.ReproducibilityDisable,
		)

		if err != nil {
			log.Error("failed exec", "err", err)
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		log.Error("failed commit", "err", err)
		return err
	}

	return nil
}

func (db *DB) UpdateLogStatus(
	ctx context.Context,
	logID int,
	status string,
) error {
	log := db.log.With(
		slog.String("location", locationForLogger+"UpdateLogStatus"),
	)

	const query = `
		UPDATE file_logs
		SET status = $2
		WHERE id = $1;
	`

	result, err := db.conn.ExecContext(
		ctx,
		query,
		logID,
		status,
	)

	if err != nil {
		log.Error("failed exec update status", "err", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Error("failed get rows affected", "err", err)
		return err
	}

	if rowsAffected == 0 {
		log.Error("log not found", "log_id", logID)
		return core.ErrNotFound
	}

	return nil
}

func (db *DB) SaveParsedLog(
	ctx context.Context,
	logID int,
	status string,
	nodes []core.Node,
	ports []core.Port,
	settings []core.SwitchSettings,
) error {
	log := db.log.With(
		slog.String("location", locationForLogger+"SaveParsedLog"),
	)

	tx, err := db.conn.BeginTx(ctx, nil)
	if err != nil {
		log.Error("failed begin tx", "err", err)
		return err
	}

	defer tx.Rollback()

	if len(nodes) > 0 {
		const query = `
			INSERT INTO nodes (
				log_id,
				node_guid,
				node_desc,
				node_type,
				num_ports,
				serial_number,
				product_name
			)
			VALUES ($1, $2, $3, $4, $5, $6, $7)

			ON CONFLICT (log_id, node_guid)
			DO UPDATE SET
				node_desc = EXCLUDED.node_desc,
				node_type = EXCLUDED.node_type,
				num_ports = EXCLUDED.num_ports,
				serial_number = EXCLUDED.serial_number,
				product_name = EXCLUDED.product_name;
		`

		stmt, err := tx.PrepareContext(ctx, query)
		if err != nil {
			log.Error("failed prepare nodes stmt", "err", err)
			return err
		}

		for _, node := range nodes {
			_, err := stmt.ExecContext(
				ctx,
				node.LogID,
				node.NodeGUID,
				node.NodeDesc,
				node.NodeType,
				node.NumPorts,
				node.SerialNumber,
				node.ProductName,
			)

			if err != nil {
				stmt.Close()

				log.Error("failed exec nodes", "err", err)
				return err
			}
		}

		stmt.Close()
	}

	if len(ports) > 0 {
		const query = `
			INSERT INTO ports (
				log_id,
				node_guid,
				port_guid,
				port_num,
				port_state,
				port_phy_state,
				link_speed_actv,
				link_width_actv
			)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8)

			ON CONFLICT (log_id, node_guid, port_guid)
			DO UPDATE SET
				port_num = EXCLUDED.port_num,
				port_state = EXCLUDED.port_state,
				port_phy_state = EXCLUDED.port_phy_state,
				link_speed_actv = EXCLUDED.link_speed_actv,
				link_width_actv = EXCLUDED.link_width_actv;
		`

		stmt, err := tx.PrepareContext(ctx, query)
		if err != nil {
			log.Error("failed prepare ports stmt", "err", err)
			return err
		}

		for _, port := range ports {
			_, err := stmt.ExecContext(
				ctx,
				port.LogID,
				port.NodeGUID,
				port.PortGUID,
				port.PortNum,
				port.PortState,
				port.PortPhyState,
				port.LinkSpeedActv,
				port.LinkWidthActv,
			)

			if err != nil {
				stmt.Close()

				log.Error("failed exec ports", "err", err)
				return err
			}
		}

		stmt.Close()
	}

	if len(settings) > 0 {
		const query = `
			INSERT INTO switch_settings (
				log_id,
				node_guid,
				endianness,
				enable_endianness_per_job,
				reproducibility_disable
			)
			VALUES ($1, $2, $3, $4, $5)

			ON CONFLICT (log_id, node_guid)
			DO UPDATE SET
				endianness = EXCLUDED.endianness,
				enable_endianness_per_job = EXCLUDED.enable_endianness_per_job,
				reproducibility_disable = EXCLUDED.reproducibility_disable;
		`

		stmt, err := tx.PrepareContext(ctx, query)
		if err != nil {
			log.Error("failed prepare settings stmt", "err", err)
			return err
		}

		for _, set := range settings {
			_, err := stmt.ExecContext(
				ctx,
				set.LogID,
				set.NodeGUID,
				set.Endianness,
				set.EnableEndiannessPerJob,
				set.ReproducibilityDisable,
			)

			if err != nil {
				stmt.Close()

				log.Error("failed exec settings", "err", err)
				return err
			}
		}

		stmt.Close()
	}

	err = db.UpdateFileLog(
		ctx,
		tx,
		logID,
		status,
	)

	if err != nil {
		log.Error("failed update file log", "err", err)
		return err
	}

	if err := tx.Commit(); err != nil {
		log.Error("failed commit tx", "err", err)
		return err
	}

	return nil
}

func (db *DB) GetNode(
	ctx context.Context,
	logID int,
	nodeGUID string,
) (core.Node, error) {
	log := db.log.With(
		slog.String("location", locationForLogger+"GetNode"),
	)

	const query = `
		SELECT
			log_id,
			node_guid,
			node_desc,
			node_type,
			num_ports,
			serial_number,
			product_name
		FROM nodes
		WHERE log_id = $1
		  AND node_guid = $2;
	`

	var node core.Node

	err := db.conn.QueryRowContext(
		ctx,
		query,
		logID,
		nodeGUID,
	).Scan(
		&node.LogID,
		&node.NodeGUID,
		&node.NodeDesc,
		&node.NodeType,
		&node.NumPorts,
		&node.SerialNumber,
		&node.ProductName,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return core.Node{}, core.ErrNotFound
		}

		log.Error("failed query row", "err", err)
		return core.Node{}, err
	}

	return node, nil
}

func (db *DB) GetPortGUIDsByNode(
	ctx context.Context,
	nodeGUID string,
) ([]string, error) {
	log := db.log.With(
		slog.String("location", locationForLogger+"GetPortGUIDsByNode"),
	)

	const query = `
		SELECT port_guid
		FROM ports
		WHERE node_guid = $1
		ORDER BY port_num;
	`

	rows, err := db.conn.QueryContext(ctx, query, nodeGUID)
	if err != nil {
		log.Error("failed query ports", "err", err)
		return nil, err
	}
	defer rows.Close()

	var guids []string

	for rows.Next() {
		var guid string

		if err := rows.Scan(&guid); err != nil {
			log.Error("failed scan port_guid", "err", err)
			return nil, err
		}

		guids = append(guids, guid)
	}

	if err := rows.Err(); err != nil {
		log.Error("rows error", "err", err)
		return nil, err
	}

	if len(guids) == 0 {
		return nil, core.ErrNotFound
	}

	return guids, nil
}

func (db *DB) GetNodesByLogID(
	ctx context.Context,
	logID int,
) ([]core.Node, error) {
	log := db.log.With(
		slog.String("location", locationForLogger+"GetNodesByLogID"),
	)

	const query = `
		SELECT
			log_id,
			node_guid,
			node_desc,
			node_type,
			num_ports,
			serial_number,
			product_name
		FROM nodes
		WHERE log_id = $1;
	`

	rows, err := db.conn.QueryContext(ctx, query, logID)
	if err != nil {
		log.Error("failed query nodes", "err", err)
		return nil, err
	}
	defer rows.Close()

	var nodes []core.Node

	for rows.Next() {
		var n core.Node

		if err := rows.Scan(
			&n.LogID,
			&n.NodeGUID,
			&n.NodeDesc,
			&n.NodeType,
			&n.NumPorts,
			&n.SerialNumber,
			&n.ProductName,
		); err != nil {
			log.Error("failed scan node", "err", err)
			return nil, err
		}

		nodes = append(nodes, n)
	}

	if err := rows.Err(); err != nil {
		log.Error("rows error", "err", err)
		return nil, err
	}

	if len(nodes) == 0 {
		return nil, core.ErrNotFound
	}

	return nodes, nil
}

func (db *DB) GetPortsByLogID(
	ctx context.Context,
	logID int,
) ([]core.Port, error) {
	log := db.log.With(
		slog.String("location", locationForLogger+"GetPortsByLogID"),
	)

	const query = `
		SELECT
			log_id,
			node_guid,
			port_guid,
			port_num,
			port_state,
			port_phy_state,
			link_speed_actv,
			link_width_actv
		FROM ports
		WHERE log_id = $1
		ORDER BY node_guid, port_num;
	`

	rows, err := db.conn.QueryContext(ctx, query, logID)
	if err != nil {
		log.Error("failed query ports", "err", err)
		return nil, err
	}
	defer rows.Close()

	var ports []core.Port

	for rows.Next() {
		var p core.Port

		if err := rows.Scan(
			&p.LogID,
			&p.NodeGUID,
			&p.PortGUID,
			&p.PortNum,
			&p.PortState,
			&p.PortPhyState,
			&p.LinkSpeedActv,
			&p.LinkWidthActv,
		); err != nil {
			log.Error("failed scan port", "err", err)
			return nil, err
		}

		ports = append(ports, p)
	}

	if err := rows.Err(); err != nil {
		log.Error("rows error", "err", err)
		return nil, err
	}

	if len(ports) == 0 {
		return nil, core.ErrNotFound
	}

	return ports, nil
}

func (db *DB) GetSettingsByLogID(
	ctx context.Context,
	logID int,
) ([]core.SwitchSettings, error) {
	log := db.log.With(
		slog.String("location", locationForLogger+"GetSettingsByLogID"),
	)

	const query = `
		SELECT
			log_id,
			node_guid,
			endianness,
			enable_endianness_per_job,
			reproducibility_disable
		FROM switch_settings
		WHERE log_id = $1;
	`

	rows, err := db.conn.QueryContext(ctx, query, logID)
	if err != nil {
		log.Error("failed query settings", "err", err)
		return nil, err
	}
	defer rows.Close()

	var settings []core.SwitchSettings

	for rows.Next() {
		var s core.SwitchSettings

		if err := rows.Scan(
			&s.LogID,
			&s.NodeGUID,
			&s.Endianness,
			&s.EnableEndiannessPerJob,
			&s.ReproducibilityDisable,
		); err != nil {
			log.Error("failed scan settings", "err", err)
			return nil, err
		}

		settings = append(settings, s)
	}

	if err := rows.Err(); err != nil {
		log.Error("rows error", "err", err)
		return nil, err
	}

	if len(settings) == 0 {
		return nil, core.ErrNotFound
	}

	return settings, nil
}
