BEGIN;

CREATE TABLE file_logs (
    id BIGSERIAL PRIMARY KEY,
    status TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    nodes_count INT NOT NULL DEFAULT 0,
    ports_count INT NOT NULL DEFAULT 0
);

CREATE TABLE nodes (
    log_id BIGINT NOT NULL REFERENCES file_logs(id) ON DELETE CASCADE,

    node_guid TEXT NOT NULL,
    node_desc TEXT,
    node_type INT,
    num_ports INT,

    serial_number TEXT,
    product_name TEXT,

    PRIMARY KEY (log_id, node_guid)
);

CREATE INDEX idx_nodes_log_id ON nodes(log_id);


CREATE TABLE ports (
    log_id BIGINT NOT NULL REFERENCES file_logs(id) ON DELETE CASCADE,

    node_guid TEXT NOT NULL,
    port_guid TEXT NOT NULL,

    port_num INT,
    port_state INT,
    port_phy_state INT,
    link_speed_actv INT,
    link_width_actv INT,

    PRIMARY KEY (log_id, node_guid, port_guid)
);

CREATE INDEX idx_ports_log_id ON ports(log_id);
CREATE INDEX idx_ports_node_guid ON ports(node_guid);


CREATE TABLE switch_settings (
    log_id BIGINT NOT NULL REFERENCES file_logs(id) ON DELETE CASCADE,

    node_guid TEXT NOT NULL,

    endianness INT,
    enable_endianness_per_job INT,
    reproducibility_disable INT,

    PRIMARY KEY (log_id, node_guid)
);

CREATE INDEX idx_switch_settings_log_id ON switch_settings(log_id);

COMMIT;