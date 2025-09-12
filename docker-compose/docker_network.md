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

