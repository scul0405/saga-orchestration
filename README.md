# SAGA ORCHESTRATION

Microservice architecture with Saga Orchestration pattern

## What has been used:
- Traefik:  edge proxy that is responsible for external traffic routing and internal grpc load-balancing.
- Services: 6 services are implemented in this project.
  - Account service: responsible for managing user accounts, tokens.
  - Product service: responsible for managing products, categories.
  - Order service: responsible for managing orders.
  - Payment service: responsible for managing payments.
  - Purchase service: responsible for managing purchases.
  - Orchestrator service: responsible for managing the saga orchestration.
- Database: 4 databases are used in this project.
  - Account database (PostgreSQL 15): responsible for storing user accounts, tokens.
  - Product database (PostgreSQL 15): responsible for storing products, categories.
  - Order database (PostgreSQL 15): responsible for storing orders.
  - Payment database (PostgreSQL 15): responsible for storing payments.
- Six-node Redis cluster
  - In-memory data store for caching.
  - Cuckoo filter for prevent cache penetration.
  - Distributed lock for preventing cache stampede.
- Kafka:  distributed event streaming platform.
  - Used for SAGA command and event.
## Quick start

Clone this repository:
```sh
git clone https://github.com/scul0405/saga-orchestration.git
cd saga-orchestration
```
#### Docker usage

Run this command:
```sh
make docker-up
```

## Monitor

### Kafkdrop
[http://localhost:9000](http://localhost:9000)

### Traefik
[http://localhost:8080](http://localhost:8080)