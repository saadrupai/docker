## **Docker :**

## **Problem: Dockerfile failing to build image**

Cause:

- Incorrect initial base image or syntax error.
- Multi-stage build setup or starting CMD/ENTRYPOINT was wrong.

Solution:

- Verified syntax and ensured the Dockerfile uses correct multi-stage builds.
- Checked that the final stage includes a valid CMD or ENTRYPOINT as per requirements and the intended image purpose.

## **Docker Compose**

## **Problem: docker-compose not starting containers**

Cause:

- Dockerfile issues causing failed builds.
- Incorrect or missing environment variables, volumes, or network definitions.
- Invalid depends_on usage or order of containers.

Solution:

- Fixed Dockerfile according to best practices and validated all docker-compose.yml paths, depends_on, environment sections, and network configuration.
- Checked logs to identify specific startup failures.

## **Docker Swarm**

## **Problem: Services/containers not showing as expected**

Cause:

- On non-manager nodes: only containers running on that node show in docker ps; Docker Swarm schedules containers across the cluster but only the manager lists all services.

Solution:

- Used docker service ls and docker stack services &lt;stack&gt; on the manager VM to monitor cluster-wide state.
- Understood node-local docker ps only lists containers present on the given host.

## **Vagrant / Vagrantfile**

## **Problem: Vagrantfile not creating VMs**

Cause:

- Default VirtualBox provider had issues on the host setup.
- VM creation failed when using certain IP ranges.
- Provisioned IP (192.168.56.100 or similar) resulted in unreachable nodes depending on how the host-only network is set up.

Solution:

- Switched to libvirt provider for better compatibility and performance on Linux systems.
- Adjusted Vagrantfile to explicitly use libvirt blocks.
- Changed VM private IP addresses to 192.168.100.x network, which worked reliably with the chosen provider.
- The original 192.168.56.100 did not work due to possible IP conflicts or local network limitations; the 192.168.100.101 subnet resolved these connectivity issues.

Why use libvirt?

- libvirt is more robust on Linux than VirtualBox, providing better integration, resource control, and fewer driver/compatibility issues. Especially useful for running multiple VMs for Docker Swarm/multihost scenarios.

## **Additional Issues & Solutions (from project context)**

## **Problem: Replica not syncing to master (Postgres)**

**Cause:**

- **Replica Postgres containers couldn't connect to master (networking or creds).
    Solution:**
- **Verified overlay network setup and environment variables for replication.**
- **Used logs and Postgres replication status queries to diagnose issues.**

## **Problem: Pgpool-II connection/auth errors**

**Cause:**

- **Mismatch in Pgpool-II backend/user credentials vs. Postgres service settings.
    Solution:**
- **Double-checked PGPOOL_\* environment variables in compose to ensure they match master DB settings.**

## **Problem: Volumes and state persistence across redeploys**

**Cause:**

- **Node-local volumes sometimes did not initialize as expected after hard-crash or reschedule.
    Solution:**
- **Manually cleared or inspected named volumes (docker volume ls, docker volume inspect, docker volume rm ...) to force proper initialization.**
- **PGPOOL admin username and password with db name should be added**

## **Networking/IP Addressing Lessons**

- **IP 192.168.56.100 (default Vagrant/VirtualBox range) sometimes led to routing/network problems on certain Linux hosts.**
- **Using the custom subnet 192.168.100.x (192.168.100.101, etc.) in the Vagrantfile was reliable and avoided collision/conflict with host or VM DHCP settings.**

## **Summary Table**

| **Tool/Area** | **Problem** | **Solution/Diagnosis** |
| --- | --- | --- |
| **Dockerfile** | **Build failed** | **Fixed multi-stage, ENTRYPOINT, and syntax** |
| --- | --- | --- |
| **Docker Compose** | **Services not starting** | **Checked Dockerfile, env vars, network, and paths** |
| --- | --- | --- |
| **Docker Swarm** | **Not all services show on workers** | **Used manager for service listing, understood Swarm roles** |
| --- | --- | --- |
| **Vagrantfile** | **VMs not starting; IP not working** | **Switched to libvirt, fixed IP range to 192.168.100.x** |
| --- | --- | --- |
| **Pgpool/Replica** | **Auth/replication issues** | **Fixed env, checked network connectivity, DB user setup** |
| --- | --- | --- |
| **Volumes/State** | **Replicas not reinitializing** | **Cleared named volumes for clean sync/reinit** |
| --- | --- | --- |
