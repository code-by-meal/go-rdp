# Go-RDP
Simple Golang-based RDP-client. Structures was taken from official MS-Documentation and from repositories. Development is going at the moment. It can login to remote RDP server, to imitate mouse movement, to proccess bitmap stream to proccess image. Sounds cool! Do more research!

# Init auth flow 

```mermaid
sequenceDiagram
    autonumber
    participant RDP Client
    participant RDP Server

    RDP Client->>RDP Server: Negotiation Connect Request (RDP/PDU)
    RDP Server-->>RDP Client: Negotiation Confirm Response (RDP/PDU)

    RDP Client->>RDP Server: Client Data Request (RDP/PDU)
    RDP Server-->RDP Client: Server Data Response (RDP/PDU)

    RDP Client->>RDP Server: Erect Domain Request (MCS/T-1.25)
    RDP Client->>RDP Server: Attach User Request (MCS/T-1.25)
    RDP Server-->>RDP Client: Attach User Confirm (MCS/T-1.25)
    RDP Client->>RDP Server: Channel Join Request N-.. (MCS/T-1.25)
    RDP Server-->>RDP Client: Channel Join Confirm N-.. (MCS/T-1.25)

    RDP Client->>RDP Server: Security Exchange - OPTIONAL (RDP/PDU)

    RDP Client->>RDP Server: Client Info (RDP/PDU)
```

# Usefull shorts

- `TCP` Transmission Control Protocol
- `TPKT` TCP Packet
- `COPT` Connection-Oriented Transport Protocol (`X224/ISO 8073`)
- `MCS` Multipoint Communication Service (`T.125`)
- `BER` Basic Encoding Rules (`ASN.1`)
- `PER` Packed Encoding Rules 
- `GCC` Generic Conference Control (`T.124`)
- `RDP` Remote Desktop Protocol (`Top level abstraction of protocol.`)
- `PDU` Protocol Data Unit 

# TODO Authorization

    [+] Negotiation proccess.
        [+] RPD security (Without any secure encryption, plain RDP-stack protocols traffic)
        [-] CredSSP
        [-] TLS
        [-] NLA
    [+] Basic settings exchange. (Without checking of certificate from server.)
    [-] Channel connection.
    [-] Security commencement.
    [-] Secure settings exchange.
    [-] Licensing.
    [-] Capabilities exchange.
    [-] Connection finalization.
    [-] Data exchange.

