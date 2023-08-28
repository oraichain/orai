## Setup nodes and validators for network

**Validator node configuration**

```bash
pex = false
persistent_peers = list of sentry nodes, optionally other vpn validators
private_peer_ids =
addr_book_strict = false

oraid start --p2p.pex false --p2p.persistent_peers ""
```

**Sentry Node Configuration**

```bash
pex	= true
unconditional_peer_ids = validator node id
persistent_peers	= validator node, optionally other vpn sentry nodes
private_peer_ids	= validator node id
addr_book_strict = false
external_address = public ip
oraid start --rpc.laddr tcp://0.0.0.0:26657 --p2p.pex true --p2p.persistent_peers "" --p2p.unconditional_peer_ids "" --p2p.private_peer_ids ""
```

**Full Node Configuration**

```bash
persistent_peers = sentry nodes
addr_book_strict = true
```
