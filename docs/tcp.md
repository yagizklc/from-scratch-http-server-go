# TCP

## TCP relies on IP

TCP is a **reliable** protocol built on top of **unreliable** IP.
Reliable := Guaranteed delivery, in order, and without errors.
Transport Layer

## IP Route Packets Between Systems (Postal Service)

### Packet

A packet is a chunk of data broken up for sending over a network. It consists of two sections:

- **Header**: Includes the source and destination addresses.
- **Data**: The payload of the packet.

### Characteristics of Packets

- Packets are **NOT GUARANTEED** to arrive at their destination.
  - Best effort delivery.
  - Sometimes lost in transit.
  - Arrival time and order are not guaranteed.

---

## TCP Guarantees

1. **Reliable Delivery of Packets**:
   1. TCP ensures that no packets are lost in transit.
   2. It does this by asking the receiver to acknowledge all sent packets, and re-transmitting any packets if an acknowledgement isn't received.
2. **In-Order Delivery of Packets**:
   1. In addition to guaranteeing packets reach their destination, TCP also guarantees that the packets are delivered in order.
   2. It does this by labelling each packet with a sequence number.
   3. The receiver tracks these numbers and reorders out-of-sequence packets. If a packet is missing, the receiver waits for it to be re-transmitted.

- if your browser opens multiple connections to Google's server, only the "source port number" will change, the rest will remain the same.
