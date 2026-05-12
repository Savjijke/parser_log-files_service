BEGIN;

CREATE TABLE logs (
    id BIGSERIAL PRIMARY KEY,

    log_name TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);


CREATE TABLE nodes (
    id BIGSERIAL PRIMARY KEY,

    log_id BIGINT NOT NULL,

    node_guid TEXT NOT NULL,
    system_image_guid TEXT,
    port_guid TEXT,

    node_desc TEXT,

    num_ports INTEGER,
    node_type INTEGER,

    class_version INTEGER,
    base_version INTEGER,

    CONSTRAINT fk_nodes_log
        FOREIGN KEY (log_id)
        REFERENCES logs(id)
        ON DELETE CASCADE,

    CONSTRAINT uq_nodes_log_guid
        UNIQUE (log_id, node_guid)
);


CREATE TABLE ports (
    id BIGSERIAL PRIMARY KEY,

    node_id BIGINT NOT NULL,

    port_guid TEXT,
    port_num INTEGER NOT NULL,

    lid INTEGER,
    local_port_num INTEGER,

    link_width_active INTEGER,
    link_width_supported INTEGER,
    link_width_enabled INTEGER,

    link_speed_enabled INTEGER,
    link_speed_active INTEGER,
    link_speed_supported INTEGER,

    port_phy_state INTEGER,
    port_state INTEGER,

    mtu_cap INTEGER,
    nmtu INTEGER,

    vl_cap INTEGER,
    op_vls INTEGER,

    retrans_active INTEGER,
    fec_active INTEGER,

    overrun_errors INTEGER,
    local_phy_errors INTEGER,

    raw_data JSONB,

    CONSTRAINT fk_ports_node
        FOREIGN KEY (node_id)
        REFERENCES nodes(id)
        ON DELETE CASCADE,

    CONSTRAINT uq_ports_node_port
        UNIQUE (node_id, port_num)
);


CREATE TABLE nodes_info (
    id BIGSERIAL PRIMARY KEY,

    node_id BIGINT NOT NULL,

COMMIT;