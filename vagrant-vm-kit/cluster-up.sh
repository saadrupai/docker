#!/usr/bin/env bash
set -e

MANAGER_FILE="Vagrantfile.manager"
WORKERS_FILE="Vagrantfile.worker"

echo "🚀 Step 1: Starting Docker Swarm manager..."
VAGRANT_VAGRANTFILE=$MANAGER_FILE vagrant up manager

echo "⏳ Step 2: Detecting manager IP..."
MANAGER_IP=$(VAGRANT_VAGRANTFILE=$MANAGER_FILE vagrant ssh manager -c "hostname -I | awk '{print \$1}'" | tr -d '\r')
if [ -z "$MANAGER_IP" ]; then
    echo "❌ Could not detect manager IP"
    exit 1
fi
echo "   Manager detected at $MANAGER_IP"

echo "⏳ Step 3: Waiting for manager ($MANAGER_IP) to become reachable..."
until ping -c1 $MANAGER_IP &>/dev/null; do
    echo "   Manager not up yet... retrying in 5s"
    sleep 5
done
echo "✅ Manager is up!"

echo "🔑 Step 4: Fetching worker join token..."
WORKER_JOIN_CMD=$(VAGRANT_VAGRANTFILE=$MANAGER_FILE vagrant ssh manager -c "
  if ! docker info --format '{{.Swarm.LocalNodeState}}' | grep -q active; then
    docker swarm init --advertise-addr $MANAGER_IP
  fi
  docker swarm join-token worker -q
" | tr -d '\r')

if [ -z "$WORKER_JOIN_CMD" ]; then
    echo "❌ Failed to get worker join token"
    exit 1
fi

echo "docker swarm join --token $WORKER_JOIN_CMD $MANAGER_IP:2377" > worker_token.sh
chmod +x worker_token.sh
echo "✅ Worker token saved to worker_token.sh"

echo "🚀 Step 5: Starting Docker Swarm workers..."
VAGRANT_VAGRANTFILE=$WORKERS_FILE vagrant up

echo "🎉 Cluster is ready!"
echo "👉 To check, run:"
echo "   VAGRANT_VAGRANTFILE=$MANAGER_FILE vagrant ssh manager -c 'docker node ls'"
