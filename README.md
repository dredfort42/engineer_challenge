# Task: Golang Development for Event Handling and Docker Integration

## Objective

The task is designed microservices using Golang and Docker, focusing on real-time data handling and querying from a time-series database. Full details are in the [TASK.md](TASK.md) file.

General schema of the services used in this project is shown in the following diagram:

-   bpmn.io diagram: [services.bpmn](services.bpmn)
-   pdf file: [services.pdf](services.pdf)

---

## Setup

This project uses **Docker** and **Docker Compose** to manage the services.

**Before running** the services, ensure you have Docker and Docker Compose installed on your system. To check if they are installed, run the following commands:

```bash
docker --version
docker-compose --version
```

If you don't have them installed, please follow the instructions below:

-   [Docker](https://docs.docker.com/engine/install/)
-   [Docker Compose](https://docs.docker.com/compose/install/)

Usually, Docker Compose comes with Docker Desktop for Windows and macOS. If you are using Linux, you may need to install it separately.

---

## Prepare the Environment

1.  Clone the repository and go to the project directory:

    ```bash
    git clone https://github.com/dredfort42/engineer_challenge && cd engineer_challenge
    ```

2.  Create the secrets directory and go to it:

    ```bash
    mkdir -p secrets &&
    cd secrets
    ```

3.  Set the InfluxDB admin `username` and `password`:

    ```bash
    echo -n "YourAdminUsernameForInfluxDB" > .env.influxdb2-admin-username &&
    echo -n "YourAdminPasswordForInfluxDB" > .env.influxdb2-admin-password
    ```

    > Replace `YourAdminUsernameForInfluxDB`, and `YourAdminPasswordForInfluxDB` with your InfluxDB admin username and password.

4.  Set the InfluxDB admin `token`:

    ```bash
    echo -n "YourAdminTokenForInfluxDB" > .env.influxdb2-admin-token
    ```

5.  Return to the project directory:

    ```bash
    cd ..
    ```

6.  If you want to use persistent storage for InfluxDB, create the `influxdb-data` directory:

    ```bash
    mkdir -p storage/influxdb-data
    ```

    And uncomment the `volumes` section in the `influxdb2` service in the `docker-compose.yml` file.

---

## Prepare the Services Configurations

### Daemon Service

Daemon service is responsible for generating random JSON events and publish these events on NATS using the subject `events`.

Configuration file: [daemon_config.json](configs/daemon_config.json)

```json
{
    "debug": true, // Enable debug mode
    "nats_url": "nats://nats:4222", // NATS server URL
    "subject": "event", // NATS subject to publish events
    "event_frequency_ms": 200 // Event generation frequency in milliseconds
}
```

### Writer Service

Writer service is responsible for subscribing to events published on NATS and write them on InfluxDB.

Configuration file: [writer_config.json](configs/writer_config.json)

```json
{
    "debug": true, // Enable debug mode
    "nats_url": "nats://nats:4222", // NATS server URL
    "subject": "event", // NATS subject to subscribe to events
    "influxdb_url": "http://influxdb2:8086", // InfluxDB server URL
    "influxdb_org": "iot", // InfluxDB organization
    "influxdb_bucket": "events_bucket", // InfluxDB bucket
    "influxdb_measurement": "event", // InfluxDB measurement
    "path_to_influxdb_token": "/app/.env.influxdb2-admin-token" // Path to InfluxDB token
}
```

### Reader Service

Reader service is responsible for receiving queries from the NATS and querying InfluxDB to get the results.

Configuration file: [reader_config.json](configs/reader_config.json)

```json
{
    "debug": true, // Enable debug mode
    "nats_url": "nats://nats:4222", // NATS server URL
    "subject": "query", // NATS subject to receive queries
    "influxdb_url": "http://influxdb2:8086", // InfluxDB server URL
    "influxdb_org": "iot", // InfluxDB organization
    "influxdb_bucket": "events_bucket", // InfluxDB bucket
    "influxdb_measurement": "event", // InfluxDB measurement
    "path_to_influxdb_token": "/app/.env.influxdb2-admin-token" // Path to InfluxDB token
}
```

### Client Service

Client service is responsible for sending queries to the NATS and printing the results.

Configuration file: [client_config.json](configs/client_config.json)

```json
{
    "debug": true, // Enable debug mode
    "nats_url": "nats://nats:4222", // NATS server URL
    "subject": "query", // NATS subject to send queries
    "event_frequency_ms": 10000 // Query frequency in milliseconds
}
```

Minimal criticality for queries **_should be set_** as environment variable `MIN_CRITICALITY`.

---

## Run the Services

To build the services, execute the following command:

```bash
docker-compose build
```

To run the services, execute the following command:

```bash
# Run the services in the foreground
docker-compose up

# Run the services in the background
docker-compose up -d

# Or run a specific service
docker-compose up <service_name>
```

To check the logs of the services, execute the following command:

```bash
# Check the logs of all services
docker-compose logs -f

# Or check the logs of a specific service
docker-compose logs -f <service_name>
```

To chack the status of the services, execute the following command:

```bash
# Check the status of all services
docker-compose ps

# Check statistics of all services
docker-compose stats
```

To restart the services, execute the following command:

```bash
# Restart the services
docker-compose restart

# Or restart a specific service
docker-compose restart <service_name>
```

To stop the services, execute the following command:

```bash
# Stop the services
docker-compose down

# Or stop a specific service
docker-compose stop <service_name>
```

---

## Debug the Services

To debug the services, you can set the `debug` flag to `true` in the configuration files of the services.

For example, to debug the daemon service, set the `debug` flag to `true` in the `daemon_config.json` file:

```json
{
    "debug": true,
    "nats_url": "nats://nats:4222",
    "subject": "event",
    "event_frequency_ms": 200
}
```

And then restart the daemon service:

```bash
docker-compose restart daemon
```

For the other services, you can follow the same steps.

If you want to debug the services using mocks you can uncomment NATS services in the `docker-compose.yml` file and run the services:

-   NATS test sub client: `nats_box_sub` for subscribing to the `events` subjects.
-   NATS test pub client: `nats_box_pub` for publishing to the `events` subjects.
-   NATS test query client: `nats_box_req` for requesting to the `query` subjects.

If you want to debug the services not in the Docker environment, you can run the services locally.
**_Note:_** You need to open the ports of the services in the `docker-compose.yml` file.

---

## Official documentation

Documentation for the services used in this project:

-   [Docker](https://docs.docker.com/)
-   [Docker Compose](https://docs.docker.com/compose/)
-   [NATS](https://docs.nats.io/)
-   [InfluxDB](https://docs.influxdata.com/influxdb/v2.0/)

If you have any questions or need help, please feel free to contact me at [dredfort.42@gmail.com](mailto:dredfort.42@gmail.com).

---

## TODO

### General tasks for this project:

-   [+] Use NATS as the message broker.
-   [+] Use InfluxDB as the time-series database.
-   [+] Implement the daemon service.
-   [+] Implement the writer service.
-   [+] Implement the reader service.
-   [+] Implement the client service.
-   [+] Implement the Docker Compose file.
-   [+] Implement the Dockerfiles.
-   [+] Implement the README file.
-   [+] Implement the BPMN diagram.
-   [-] Use optimal configurations for the services.
-   [-] Use ooptimal queries for the InfluxDB.

### Security tasks for this project:

-   [+] Use different networks for the services.
-   [+] Close the unnecessary ports for the services.
-   [+] Use secrets for the sensitive data.
-   [+] Use the latest versions of the base images.
-   [+] Use the latest versions of the packages.
-   [+] Log errors and security events for the services.
-   [-] Use monitoring for the services.
-   [-] Use different users for the services.
-   [-] Use SSL/TLS for the services.
-   [-] Use Hashicorp Vault for the secrets.
-   [-] Use Unit tests for the services with coverage more than 80%.
-   [-] Use static security analysis for the services.
-   [-] Use Integration tests for the services.
