#!/usr/bin/env bash
set -e

MANAGER_FILE="Vagrantfile.manager"
WORKERS_FILE="Vagrantfile.worker"

echo "🛑 Step 1: Stopping Docker Swarm workers..."
VAGRANT_VAGRANTFILE=$WORKERS_FILE vagrant halt || true

echo "🛑 Step 2: Stopping Docker Swarm manager..."
VAGRANT_VAGRANTFILE=$MANAGER_FILE vagrant halt || true

echo "🗑 Step 3: Destroying VMs..."
VAGRANT_VAGRANTFILE=$WORKERS_FILE vagrant destroy -f || true
VAGRANT_VAGRANTFILE=$MANAGER_FILE vagrant destroy -f || true

echo "🧹 Step 4: Cleaning up temporary files..."
rm -f worker_token.sh

# -----------------------------
# Optional: Force cleanup of leftover hypervisor resources
# -----------------------------
echo "⚠️ Step 5: Force cleaning VirtualBox VMs if still present..."
for vm in vagrant_manager vagrant_worker1 vagrant_worker2; do
    VBoxManage list vms | grep -q "$vm" && VBoxManage unregistervm "$vm" --delete || true
done

echo "⚠️ Step 6: Force cleaning libvirt VMs and volumes if still present..."
for vm in vagrant_manager vagrant_worker1 vagrant_worker2; do
    sudo virsh list --all --name | grep -q "$vm" && sudo virsh destroy "$vm" || true
    sudo virsh list --all --name | grep -q "$vm" && sudo virsh undefine "$vm" --remove-all-storage || true
    sudo virsh vol-list --pool default | grep -q "$vm.img" && sudo virsh vol-delete --pool default "$vm.img" || true
done
sudo systemctl restart libvirtd || true

echo "✅ Cluster fully cleaned!"
