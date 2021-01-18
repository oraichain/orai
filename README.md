# Installation

```bash
docker-compose up -d && docker-compose exec orai bash
wget https://github.com/CosmWasm/wasmvm/releases/download/v0.13.0/libwasmvm_muslc.a -O /lib/libwasmvm_muslc.a
make build-oraid

# setup blockchain and run
./docker/setup_wasmd.sh && ./docker/run_wasmd.sh 
```