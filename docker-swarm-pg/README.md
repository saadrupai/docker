# Swarm-overlay-hands-on: Multi-Node PostgreSQL + Pgpool-II on Docker Swarm
This hands-on project sets up a high-availability PostgreSQL cluster with streaming replication and load-balancing, using Docker Swarm across 3 VMs. Pgpool-II handles query routing for reads and writes. Overlay networks enable secure, internal service communication.

## Prerequisites
- All Vagrant VMs up and reachable (vagrant up)
- Docker, Docker Compose, and Swarm initialized on all VMs
- Swarm worker nodes joined (docker swarm join ...)
- Overlay networks created on manager:
```
docker network create -d overlay --attachable --subnet 10.10.0.0/24 overlay_internal
docker network create -d overlay --attachable --subnet 10.11.0.0/24 overlay_external
```
- Public Docker Hub or local registry with needed images (update compose file if using local registry)

## Step-by-Step Deployment
1. Clone repo and enter project directory:
```
git clone <repo_url>
cd docker-swarm-pg
```
2. Verify Swarm overlay networks exist:
```
docker network ls | grep overlay_internal
docker network ls | grep overlay_external
```
3. Deploy the stack:
```
docker stack deploy -c docker-compose.yml pgcluster
```
4. Monitor your services:
```
docker stack services pgcluster
docker service ps pgcluster_postgres-master
docker service ps pgcluster_postgres-replica
docker service ps pgcluster_pgpool
```

## Testing & Verification
Connect to Pgpool (entrypoint for clients)
Pgpool is published on port 5432 by routing mesh. Use a temporary container to test:
```
docker run --rm -it --network overlay_external postgres:14 \
  psql -h pgpool -U postgres -d mydb
```
**Password: masterpass**

---

Testing & Validation:

1. Create table:
```
CREATE TABLE IF NOT EXISTS test_table (
  id serial PRIMARY KEY,
  msg text,
  created_at timestamptz DEFAULT now()
);

INSERT INTO test_table (msg) VALUES ('hello from master init');
```
2. Check backend (master/replica):
```
SELECT inet_server_addr(), inet_server_port(), pg_is_in_recovery();
```
3. Insert (writes):
```
INSERT INTO test_table (msg) VALUES ('inserted via pgpool at ' || now());
```
4. Read (should be load-balanced across replicas):
```
SELECT id, msg, created_at FROM test_table ORDER BY id DESC LIMIT 5;
```
5. Verify replication status:
- On master (pg_stat_replication):
```
MASTER_CID=$(docker ps -f "name=pgcluster_postgres-master" -q | head -n1)
docker exec -it $MASTER_CID psql -U postgres -d mydb -c "SELECT pid, state, client_addr, sync_state FROM pg_stat_replication;"
```
- On a replica:
```
REPLICA_CID=$(docker ps -f "name=pgcluster_postgres-replica" -q | head -n1)
docker exec -it $REPLICA_CID psql -U postgres -d mydb -c "SELECT pg_is_in_recovery(), now();"
```

## Disaster Recovery (Simulated Failures)
A. Crash a replica
```
docker ps --format "{{.ID}} {{.Names}}" | grep postgres-replica
docker rm -f <replica_container_id>
```
B. Stop Docker daemon on worker (node failure)
```
sudo systemctl stop docker
# Observe service rescheduling on the manager
docker service ps pgcluster_postgres-replica
```
C. Crash master
```
docker rm -f <master_container_id>
```

## Cleanup
```
docker stack rm pgcluster
docker volume ls | grep pgcluster
docker volume rm swarm-overlay-hands-on_master_data
docker volume rm swarm-overlay-hands-on_replica_data
```

For advanced scenarios (auto failover, custom configs, private registry), see docs and reach out for Patroni or repmgr setups.
