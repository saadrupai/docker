# Makefile for Docker Swarm Vagrant cluster

# Scripts
UP_SCRIPT := cluster-up.sh
DOWN_SCRIPT := cluster-down.sh

# -----------------------------
# Default target
# -----------------------------
.PHONY: all
all: up

# -----------------------------
# Ensure scripts are executable
# -----------------------------
.PHONY: chmod
chmod:
	@echo "🔧 Making scripts executable..."
	chmod +x $(UP_SCRIPT) $(DOWN_SCRIPT)

# -----------------------------
# Bring the cluster up
# -----------------------------
.PHONY: up
up: chmod
	@echo "🚀 Bringing the cluster up..."
	./$(UP_SCRIPT)

# -----------------------------
# Tear the cluster down
# -----------------------------
.PHONY: down
down: chmod
	@echo "🛑 Tearing the cluster down..."
	./$(DOWN_SCRIPT)

# -----------------------------
# Restart cluster
# -----------------------------
.PHONY: restart
restart: down up

# -----------------------------
# Clean temporary files
# -----------------------------
.PHONY: clean
clean:
	@echo "🧹 Cleaning temporary files..."
	rm -f worker_token.sh

# -----------------------------
# Full clean including Vagrant state
# -----------------------------
.PHONY: cleanall
cleanall: down clean
	@echo "🧹 Removing .vagrant folder..."
	rm -rf .vagrant

# -----------------------------
# Help target
# -----------------------------
.PHONY: help
help:
	@echo ""
	@echo "Docker Swarm Vagrant Cluster Makefile"
	@echo "------------------------------------"
	@echo "Available targets:"
	@echo "  make up        - Bring the cluster up"
	@echo "  make down      - Tear the cluster down"
	@echo "  make restart   - Restart the cluster (down + up)"
	@echo "  make chmod     - Ensure scripts are executable"
	@echo "  make clean     - Remove temporary files"
	@echo "  make cleanall  - Remove temporary files + .vagrant folder"
	@echo "  make help      - Show this help message"
	@echo ""
