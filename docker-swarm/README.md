# Docker Swarm Overlay Networking Hands-On

This project demonstrates setting up a Docker Swarm cluster with overlay networks using Vagrant to provision local Ubuntu VMs. It covers creating an internal API overlay, an external frontend overlay, deploying services with a local Docker registry, and inspecting low-level VXLAN networking details.

---

## Project Layout

- `Vagrantfile` - Provisions 3 Ubuntu 20.04 VMs (1 manager, 2 workers) with Docker CE installed.
- `docker-compose.yml` - Docker Swarm stack definition for API and frontend services.
- `nginx.conf` - Nginx config for frontend proxying API requests.
- `api/main.go` - Simple Go HTTP API returning hostname and client IP.
- `api/Dockerfile` - Multi-stage Go build for a lightweight API container image.
- `README.md` - This file.

---

## Setup

### Provision VMs

Run:
```
vagrant up --provider=libvirt
```
This creates and configures 3 VMs with static private IPs, Docker installed, and network settings for Swarm.

### Configure Local Registry as Insecure

On each VM (`manager`, `worker1`, `worker2`) run:

```
sudo tee /etc/docker/daemon.json <<EOF
{
"insecure-registries": ["192.168.56.101:5000"]
}
EOF
sudo systemctl restart docker
```
This allows pulling images from the local registry on the manager via HTTP.

---

## Initialize Swarm Cluster

On the manager VM:
```
docker swarm init --advertise-addr 192.168.100.101
```
Copy the printed join token command.

Then on each worker VM run the join command, for example:
```
docker swarm join --token <TOKEN> 192.168.100.101:2377
```

Verify cluster nodes on manager:
```
docker node ls
```

---

## Create Overlay Networks

On the manager:
```
docker network create -d overlay --attachable --subnet 10.10.0.0/24 overlayinternal
docker network create -d overlay --attachable --subnet 10.11.0.0/24 overlayexternal
```

- `overlayinternal` is an internal API-only network.
- `overlayexternal` is the frontend network connected to external traffic.

---

## Local Docker Registry

Start a local registry on the manager:
```
docker run -d --name registry --restart always -p 5000:5000 registry:2
```

---

## Build and Push API Image

From project root on manager:
```
docker build -t 192.168.100.101:5000/simple-api:1.0 ./api
docker push 192.168.100.101:5000/simple-api:1.0
```

---

## Deploy Stack

Deploy the services using docker-compose.yml:
```
docker stack deploy -c docker-compose.yml myapp
```
---

Check services and their tasks:
```
docker stack services myapp
docker service ps myapp_api
docker service ps myapp_frontend
```

---

## Testing & Validation

1. Curl frontend service from host or any node IP:
```
curl http://192.168.100.101:8080
```
You should get the proxied API response showing the API container host.

2. Enter frontend container and test internal API reachability:
```
FRONTCID=$(docker ps -qf name=myapp_frontend | head -n 1)
docker exec -it $FRONTCID sh
apk add --no-cache curl
curl -s http://api:8080
```

3. Run a debug container on the internal overlay to ping or curl API service:
```
docker run --rm -it --network overlayinternal alpine sh
apk add --no-cache curl iputils
ping api
curl http://api:8080
```

4. Inspect VXLAN networking on hosts:
```
ip -d link show type vxlan
sudo bridge fdb show
sudo tcpdump -n -i any udp port 4789 -vv
```

---

## Simulate Failures

- Graceful drain worker:
```
docker node update --availability drain <worker-name>
```

- Hard failure by stopping Docker on worker:
```
sudo systemctl stop docker
```

Check rescheduling of tasks on manager.

---

## Cleanup

On manager:
```
docker stack rm myapp
docker network rm overlayinternal overlayexternal
docker stop registry && docker rm registry
docker swarm leave --force
```

On workers:
```
docker swarm leave
```

Optionally destroy Vagrant VMs:
```
vagrant destroy -f
```

---

## How This Works: Summary

- Swarm manager runs a Raft-based cluster state.
- Overlay networks create VXLAN interfaces on each node, bridging container endpoints.
- Containers on different hosts communicate via VXLAN encapsulated UDP packets (port 4789).
- Swarm ensures global endpoint state distribution and load balancing via the routing mesh.
- Local registry enables image distribution for services across nodes without external internet.

---

## Quick Commands Cheat Sheet

```
vagrant up
docker swarm init --advertise-addr <manager-ip>
docker swarm join --token <token> <manager-ip>:2377
docker network create -d overlay --attachable overlayinternal overlayexternal
docker run -d --name registry -p 5000:5000 registry:2
docker build -t <registry-ip>:5000/simple-api:1.0 ./api
docker push <registry-ip>:5000/simple-api:1.0
docker stack deploy -c docker-compose.yml myapp
```

---

This README covers all steps to reproduce, deploy, and test a multi-node Docker Swarm overlay network setup with a real application example.

---

If preferred, a ready-to-run Git repo or automation script can be provided.

---

Feel free to ask for additional scripts or files automation.
