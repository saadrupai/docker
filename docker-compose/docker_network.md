## Docker Network
Docker networking is a system that enables Docker containers to communicate with each other, the Docker host, and external networks. It offers isolation and connectivity through various components and drivers that manage how containers interact in a networked environment

### Overview of Docker Network Components

- **Container Network Model (CNM):** The framework defining Docker's network design, consisting of three main building blocks:
    - **Sandboxes:** Isolated network stacks assigned to containers, containing interfaces like Ethernet, DNS, routing tables, and ports.
    - **Endpoints:** Virtual interfaces that connect containers to networks, providing communication channels.
    - **Networks:** Software bridges or switches that connect endpoints, allowing container communication within the same network while isolating them from others.
- **Libnetwork:** The implementation of CNM used by Docker, akin to how TCP/IP operates for the OSI model. It manages network creation, connection, and maintenance.
- **Network Drivers:** The core of Docker networking functionality, providing different network topologies and capabilities:
    - **bridge:** Default network driver that creates a software bridge for container communication on a single host. If more containers are on same bridge connected
  with same bridge network, they can communicate with each other witour NAT cause they are in same network. But if they want to snd packets accross containers in 
    different bridge networks, they need to use NAT. NAT is used to translate private IP addresses to public IP addresses and vice versa. Nat is done using
  the host's iptables rules.Nat masqurades the source IP address of outgoing packets with the host's IP address, allowing return traffic to be routed back to the host and then to the appropriate container.
    - **host:** Shares the host’s network stack directly with containers without isolation. So the containers use the host's IP address and port space.
    - **none:** Isolates a container completely from any network connectivity. Only loopback interface is available.So they can only communicate with themselves.
    - **overlay:** Connects containers across multiple Docker hosts for distributed systems and Docker Swarm.
    - **ipvlan:** Provides granular control over container IP addressing, useful for VLAN setups. But unlike macvlan, ipvlan does not assign a unique MAC address to each container.
  Rather more conainers on same ipvlan network can share same MAC address. This is useful in scenarios where MAC address conservation is important, such as in large-scale deployments or when integrating with existing network infrastructures that have MAC address limitations.
  But this has some limitations like it does not support communication between containers on the same host using the ipvlan driver. This is because the host's network stack does not recognize the shared MAC address as belonging to different containers.
  ALso DHCP(working is mentioned below) servers finds it difficult to assign IP addresses to multiple containers sharing the same MAC address.
    - **macvlan:** Assigns a unique MAC address to a container, making it appear as a physical network device. This allows direct communication with the external network, bypassing the host's network stack.
  But macvlan has some limitations like it does not support communication between containers on the same host using the macvlan driver. This is because each container has its own unique MAC address, and the host's network stack does not recognize them as belonging to the same host.
  Also if there are limitations on MAC addresses in the network, it can be a problem when there are many containers.

**DHCP and IP Addressing:**
DORA Process:
- Docker uses a DHCP-like process to assign IP addresses to containers. This process is often referred to as DORA (Discover, Offer, Request, Acknowledge):
- Discover: When a container starts, it sends a DHCPDISCOVER message to find available IP addresses.
- Offer: The Docker internal DHCP server responds with a DHCPOFFER message, providing an available IP address.
- Request: The container sends a DHCPREQUEST message to request the offered IP address.
- Acknowledge: The DHCP server sends a DHCPACK message to confirm the IP address assignment.
This process ensures that each container receives a unique IP address within the network's subnet.

- Docker uses an internal DHCP server to assign IP addresses to containers within a network.
- IP addresses are allocated from a predefined subnet associated with the network.
- Users can specify custom subnets and IP ranges when creating user-defined networks.

### Key Network Concepts

- Containers can be connected to multiple networks simultaneously.
- Default gateway selection can be controlled to specify network priority.
- User-defined networks allow custom networking between containers with hostname and IP-based communication.
- Network isolation ensures that containers on different networks cannot communicate unless explicitly routed.

This overview captures the essential Docker networking components and how they interact to provide flexible and isolated networking for containerized environments.



# Docker Overlay Network Deep Dive

## What is an Overlay Network?

An **overlay network** is a virtual network that operates on top of physical host networks, enabling containers running on different Docker hosts to communicate securely as if they were part of the same local network. It abstracts physical network boundaries and uses VXLAN tunneling to encapsulate container traffic for multi-host communication.

Overlay networks are critical for distributed applications and Docker Swarm services, allowing containers across hosts to communicate seamlessly without exposing the underlying host networks directly.

---

## How Docker Overlay Network Works Under the Hood

- **Network Namespaces:** Each container runs in an isolated network namespace with its own interfaces, routing tables, and DNS settings.
- **VXLAN(Virtual Extensible Lan) Encapsulation:** Docker encapsulates Layer 2 Ethernet frames inside Layer 4 UDP packets (using VXLAN) enabling these frames to be transported across disparate Layer 3 IP networks.
- **VXLAN Tunnel Endpoints (VTEPs):** Each Docker host runs a VTEP(VXLAN Tunnel Endpoint) interface. It encapsulates outgoing traffic with VXLAN headers and UDP/IP headers, and decapsulates incoming VXLAN packets.
- **VXLAN Network Identifier (VNI(VXLAN Network Identifier)):** A unique 24-bit value identifying each overlay network, allowing up to 16 million isolated virtual networks on the same physical infrastructure.
- **Distributed KV Store:** Docker Swarm managers use the Raft consensus algorithm to maintain a distributed key-value store that synchronizes network state (endpoints, IP assignments) across all nodes.
- **Service Discovery & DNS:** Built into overlay networks, enabling containers to discover and communicate with services by name across hosts.

---

## VXLAN Packet Forwarding Chain

1. **Source Container Sends Packet:** Container generates a Layer 2 Ethernet frame.
2. **Local Linux Bridge:** Frame is forwarded to the host's Linux bridge.
3. **VXLAN Encapsulation:** The local VTEP encapsulates the Ethernet frame inside a VXLAN (UDP) packet, adding VNI.
4. **Destination VTEP IP Resolution:** Source VTEP determines the destination host VTEP IP via MAC-to-VTEP mappings learned dynamically.
5. **Packet Routing over IP Network:** Encapsulated VXLAN packet is routed over the physical Layer 3 network to the destination VTEP’s IP.
6. **Decapsulation at Destination VTEP:** VXLAN and UDP/IP headers removed to reveal original Ethernet frame.
7. **Local Bridge Delivery:** Destination host’s Linux bridge forwards the frame to the target container’s virtual interface.
8. **Container Receives Packet:** Destination container processes the Ethernet frame and payload.

---

## Why VXLAN Uses UDP in Encapsulation

- **Transport Flexibility:** UDP allows VXLAN to carry Layer 2 frames over any IP Layer 3 network, without requiring direct Layer 2 adjacency.
- **Connectionless Protocol:** UDP's stateless nature avoids session overhead, minimizing latency and improving efficiency.
- **Multicast and Broadcast Support:** VXLAN uses UDP multicast or control plane protocols for flooding unknown unicast or broadcast frames in the overlay.
- **Scalability:** UDP encapsulation allows massive scale (up to 16 million VNIs) over routed networks.
- **Firewall and Router Compatibility:** UDP packets traverse existing IP infrastructure more easily than other protocols.

---

## NAT and VXLAN Overlay Networks

- **No NAT within Overlay:** Since VXLAN extends a virtual Layer 2 network, containers across hosts appear on the same subnet, communicating directly by IP without NAT.
- **When is NAT Required?** When containers communicate outside the VXLAN overlay (e.g., internet access) or when subnet overlaps occur, NAT or masquerading might be configured on hosts.
- **Bridge vs Overlay NAT:** Docker applies NAT for outbound internet from bridge networks but typically not across overlay VXLAN networks.

---

## Example Scenario: Single VXLAN Overlay Network

- Docker Swarm cluster with 3 hosts (A, B, C), each hosting 2 containers.
- One overlay network defined with VNI 5000 common to all containers.
- Container 1 on Host A communicating with Container 2 on Host C:
    - Packet flows through local bridge, then encapsulated in VXLAN with VNI 5000.
    - Routed over IP to Host C's VTEP, decapsulated.
    - Delivered via Host C’s Linux bridge to Container 2.

---

## Example Scenario: Multiple VXLAN Overlay Networks

- Same cluster with 2 overlay networks: VNI 5000 and VNI 6000.
- Containers assigned to networks based on their service requirements.
- Traffic isolation ensured by unique VNIs.
- Cross-network container communication requires multi-network attachment, routing, or NAT.

---

## Achieving Isolated Multi-Host Container Communication Without External Access

- Create an **internal overlay network** with `internal: true`.
- Containers connected to this network communicate with each other across hosts.
- Network is isolated from the host’s physical network and the internet, blocking external access.
- Useful for secure inter-service communication in distributed apps with no external exposure.

---

# Summary

Docker Overlay Networks provide seamless, secure, and scalable Layer 2 virtual networking atop Layer 3 infrastructure using VXLAN tunnels and VNIs. They enable multi-host container communication essential for orchestrated environments like Docker Swarm, without exposing host networks or requiring NAT within the overlay. Features like distributed state synchronization, service discovery, and DNS make container-to-container communication effortless and reliable.

---

