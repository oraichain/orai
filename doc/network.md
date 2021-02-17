## Setup nodes and validators for network

**Validator node configuration**
```bash
pex = false
persistent_peers = list of sentry nodes
private_peer_ids = 
addr_book_strict = false
```

**Sentry Node Configuration**
```bash
pex	= true
unconditional-peer-ids = validator node id
persistent_peers	= validator node, optionally other sentry nodes
private_peer_ids	= validator node id
addr_book_strict = false
external_address = public ip
start --rpc.laddr tcp://0.0.0.0:26657
```

**Full Node Configuration**
```bash
seeds = sentry nodes
addr_book_strict = true
```