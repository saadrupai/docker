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




# **Docker Overlay Networking (with VXLAN)**

## **What is an Overlay Network?**

-   An **overlay network** is a virtual network built on top of physical
    > host networks.

-   It allows containers across **different Docker hosts** to
    > communicate securely **as if they were on the same local
    > network**.

-   It works by encapsulating traffic using **VXLAN (Virtual Extensible
    > LAN) tunneling**, creating secure tunnels between hosts.

-   Overlay networks are critical for **multi-host container
    > communication**, especially for **Docker Swarm services**.

-   Benefits:

    -   Multi-host connectivity

    -   Built-in **service discovery & DNS\
        > **

    -   Secure, isolated communication

    -   Scales across data centers

## **How Docker Overlay Networks Work (Under the Hood)**

1.  **Network namespaces** isolate each container's networking stack.

2.  Docker creates a **Linux bridge** on each host to connect
    > containers.

3.  **VXLAN encapsulation** wraps Layer 2 Ethernet frames into Layer 4
    > UDP packets.

4.  **VXLAN tunnels** connect hosts in the Swarm, enabling
    > container-to-container traffic across machines.

5.  A **distributed key-value store** (via Swarm's Raft consensus)
    > synchronizes state across nodes:

    -   Endpoint info

    -   Overlay subnet IP assignments

    -   VNI mappings

6.  Each container receives an IP address from the **overlay subnet**.

7.  Built-in **DNS service discovery** allows containers to resolve each
    > other by name.

**Path Example:**

```
Container -\> Linux Bridge -\> VXLAN Interface -\> Physical Host Network
```

## **VXLAN (Virtual Extensible LAN)**

### **How VXLAN Works**

-   VXLAN enables creation of **virtual Layer 2 networks over Layer 3 IP
    > networks**.

-   Ethernet frames (L2) are encapsulated into UDP packets (L4).

-   **Key components:\
    > **

    -   **VTEP (VXLAN Tunnel Endpoint):** Encapsulates/decapsulates
        > VXLAN packets. Usually one per host.

    -   **VNI (VXLAN Network Identifier):** Identifies which overlay
        > network a packet belongs to (like VLAN ID). Supports up to
        > **16M isolated networks**.

### **Encapsulation/Decapsulation Process**

1.  **Container sends Ethernet frame** → goes to host's Linux bridge.

2.  **Linux bridge forwards frame** → to host's VXLAN interface (VTEP).

3.  **Encapsulation:\
    > **

    -   Ethernet frame wrapped inside UDP packet.

    -   VXLAN header added, including **VNI**.

4.  **VTEP determines destination:\
    > **

    -   Learns remote VTEP IP via ARP/MAC learning or control-plane
        > sync.

    -   Maps container MAC → VTEP IP.

5.  **Send packet across physical network.\
    > **

6.  **Destination VTEP decapsulates** UDP/VXLAN headers → retrieves
    > original frame.

7.  **Local Linux bridge forwards** frame to destination container.

**Flow Visualization:**

```
→ Source Container

→ Source Bridge

→ Source VTEP (encapsulation, VNI)

→ Physical IP Network

→ Destination VTEP (decapsulation)

→ Destination Bridge

→ Destination Container
```

## **Why VXLAN Uses UDP**

-   **Encapsulation & Flexibility:** Packets can travel over any IP
    > network.

-   **Connectionless:** UDP requires no session, reducing overhead.

-   **Compatibility:** UDP (default port 4789) traverses
    > routers/firewalls easily.

-   **Scalability:** Supports millions of VNIs across data centers.

-   **Multicast support:** Efficient flooding of unknown/broadcast
    > traffic (BUM).

### **Broadcast, Unknown Unicast, and Multicast (BUM) Handling**

-   **Flooding:** VTEPs replicate packets to all other VTEPs in the same
    > VNI.

-   **Known unicast:** Sent directly to destination VTEP (no flooding).

-   **Efficiency:** Reduces broadcast storms by learning MAC-to-VTEP
    > mappings.

## **NAT in Overlay Networks**

-   **Not required** for container-to-container traffic within the
    > overlay.

-   Containers share the **same subnet across hosts** → VXLAN tunnels
    > carry packets without rewriting IPs.

### **When NAT *is* needed:**

-   External internet access.

-   Overlapping subnets.

-   Security/load-balancing policies.

**Bridge vs Overlay:**

-   **Bridge Network:** NAT applied for external communication.

-   **Overlay Network:** No NAT needed inside overlay; NAT only applied
    > for outside traffic.

## **Example: Docker Swarm Overlay with VXLAN**

### **Scenario**

-   Swarm cluster: **3 hosts (A, B, C)**.

-   Each host runs 2 containers.

-   All attached to one overlay network (e.g., **VNI 5000**).

### **Packet Flow (Container on Host A → Container on Host C)**

1.  All containers get IPs from overlay subnet (e.g., 10.0.0.0/24).

    -   Host A, Container 1: 10.0.0.2

    -   Host C, Container 2: 10.0.0.6

2.  Container 1 sends Ethernet frame → Host A's Linux bridge.

3.  Bridge forwards → Host A's VXLAN VTEP.

4.  Encapsulation:

    -   Ethernet frame → UDP + VXLAN header (VNI = 5000).

5.  Source VTEP maps Container 2's MAC → Host C's VTEP IP.

6.  VXLAN packet routed over physical Layer 3 network.

7.  Host C's VTEP decapsulates → retrieves Ethernet frame.

8.  Local bridge delivers frame to Container 2.

## **Key Takeaways**

-   Overlay networks allow **cross-host container communication**
    > without exposing physical networks.

-   VXLAN encapsulates Ethernet frames into UDP packets, carried across
    > IP underlay.

-   **VNI & VTEPs** ensure correct isolation and delivery.

-   **NAT is unnecessary** within overlay; only used for external
    > communication.

-   **Docker Swarm** manages synchronization of overlay state via Raft.