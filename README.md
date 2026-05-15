# 📦 Log Parser Service

Microservice in Go for parsing log archives, extracting nodes/ports/settings and storing results in PostgreSQL.
Provides REST API for accessing parsed data.

---

## 🚀 How to run

### 1. Clone project

```bash
git clone <repo-url>
cd parser_log-files_service
```

---

### 2. Start with Docker

```bash
docker compose up --build
```

---

### 3. Services

| Service    | URL                                            |
| ---------- | ---------------------------------------------- |
| API        | [http://localhost:8080](http://localhost:8080) |
| PostgreSQL | localhost:5432                                 |
| pgAdmin    | [http://localhost:5050](http://localhost:5050) |

pgAdmin credentials:

```
email: admin@admin.com
password: admin
```

---

## ⚙️ Environment

Configured via docker-compose:

```env
DATABASE_URL=postgres://postgres:postgres@db:5432/parser?sslmode=disable
PORT=8080
LOG_LEVEL=info
```

---

## 📡 API

### 📤 Parse log

**POST** `/api/v1/parse/`

Response:

```json
{
  "log_id": 1
}
```

---

### 📊 Get log info

**GET** `/api/v1/log/{log_id}`

---

### 🧠 Get node

**GET** `/api/v1/node/{node_id}`

---

### 🔌 Get node ports

**GET** `/api/v1/port/{node_id}`

---

### 🌐 Get topology

**GET** `/api/v1/topology/{log_id}`

---

## 🧪 curl examples

### Parse log

```bash
curl -X POST http://localhost:8080/api/v1/parse/ \
  -H "Content-Type: application/json" \
  -d '{"path":"/app/data/ibdiagnet2.zip"}'
```

Response:

```json
{
  "log_id": 1
}
```

---

### Get log info

```bash
curl http://localhost:8080/api/v1/log/1
```

Response:

```json
{
  "id": 1,
  "status": "ready",
  "created_at": "2026-05-14T23:37:47Z",
  "nodes_count": 5,
  "ports_count": 5
}
```

---

### Get node

```bash
curl -X GET http://localhost:8080/api/v1/node/0xswitch1?logID=1
```

Response:

```json
{
  "log_id": 1,
  "node_guid": "0xswitch1",
  "node_desc": "SWITCH_1",
  "node_type": 2,
  "num_ports": 65,
  "serial_number": "SOS123",
  "product_name": "Gorilla"
}
```

---

### Get ports

```bash
curl -X GET http://localhost:8080/api/v1/port/0xswitch1?logID=1
```

Response:

```json
["0xswitch1"]
```

---

### Get topology

```bash
curl -X GET http://localhost:8080/api/v1/topology/1
```

Response:

```json
{
  "groups": [
    {
      "type": "host",
      "nodes": [
        {
          "node": {
            "node_guid": "0xhost1",
            "node_desc": "HOST_1",
            "node_type": "host",
            "num_ports": 1,
            "serial_number": "",
            "product_name": ""
          },
          "ports": [
            {
              "port_guid": "0xhost1",
              "port_num": 1,
              "port_state": 5,
              "port_phy_state": 4,
              "link_speed_actv": 2048,
              "link_width_actv": 2
            }
          ]
        }
      ]
    },
    {
      "type": "switch",
      "nodes": [
        {
          "node": {
            "node_guid": "0xswitch1",
            "node_desc": "SWITCH_1",
            "node_type": "switch",
            "num_ports": 65,
            "serial_number": "SOS123",
            "product_name": "Gorilla"
          },
          "ports": [
            {
              "port_guid": "0xswitch1",
              "port_num": 65,
              "port_state": 5,
              "port_phy_state": 4,
              "link_speed_actv": 2052,
              "link_width_actv": 2
            }
          ],
          "settings": {
            "endianness": 10,
            "enable_endianness_per_job": 0,
            "reproducibility_disable": 0
          }
        },
        {
          "node": {
            "node_guid": "0xswitch2",
            "node_desc": "SWITCH_2",
            "node_type": "switch",
            "num_ports": 65,
            "serial_number": "PTSR50",
            "product_name": "Gorilla Prod"
          },
          "ports": [
            {
              "port_guid": "0xswitch2",
              "port_num": 65,
              "port_state": 5,
              "port_phy_state": 4,
              "link_speed_actv": 2052,
              "link_width_actv": 2
            }
          ],
          "settings": {
            "endianness": 0,
            "enable_endianness_per_job": 0,
            "reproducibility_disable": 30
          }
        },
        {
          "node": {
            "node_guid": "0xswitch3",
            "node_desc": "SWITCH_3",
            "node_type": "switch",
            "num_ports": 65,
            "serial_number": "GGVP79",
            "product_name": "Gorilla CLust"
          },
          "ports": [
            {
              "port_guid": "0xswitch3",
              "port_num": 65,
              "port_state": 5,
              "port_phy_state": 4,
              "link_speed_actv": 2052,
              "link_width_actv": 2
            }
          ],
          "settings": {
            "endianness": 0,
            "enable_endianness_per_job": 0,
            "reproducibility_disable": 4
          }
        },
        {
          "node": {
            "node_guid": "0xswitch4",
            "node_desc": "SWITCH_4",
            "node_type": "switch",
            "num_ports": 65,
            "serial_number": "LOX505",
            "product_name": "Gorilla"
          },
          "ports": [
            {
              "port_guid": "0xswitch4",
              "port_num": 65,
              "port_state": 5,
              "port_phy_state": 4,
              "link_speed_actv": 2052,
              "link_width_actv": 2
            }
          ],
          "settings": {
            "endianness": 0,
            "enable_endianness_per_job": 1,
            "reproducibility_disable": 0
          }
        }
      ]
    }
  ]
}
```

---

## 🟢 Done

Service runs fully via Docker Compose and stores parsed log data in PostgreSQL.
