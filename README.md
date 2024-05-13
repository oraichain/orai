# Orai monorepo

Cosmos based blockchain integrated with Smart Contracts [Orai](https://orai.io).

## Reporitories

| Name                               | Description                                                                           |
| ---------------------------------- | ------------------------------------------------------------------------------------- |
| [`orai`](orai)                     | The world’s first AI-powered oracle and ecosystem for blockchains                     |
| [`oraivisor`](oraivisor)           | A small process manager around Oraichain binaries that monitors the governance module |
| [`interchaintest`](interchaintest) | Docker containers for hooks testing of IBC-compatible blockchains                     |

## Docker Build

```bash
# dev
docker build -t <image-tag> -f orai/Dockerfile --build-arg WASMVM_VERSION=v1.5.2 --build-arg VERSION=v0.42.0 .

# prod
docker build -t <image-tag> -f orai/Dockerfile.prod --build-arg WASMVM_VERSION=v1.5.2 --build-arg VERSION=v0.42.0 .
```

## Upgrade command

```bash
oraid tx gov submit-proposal software-upgrade "v0.42.0" --title "upgrade Oraichain network to v0.42.0" --description "Please visit https://github.com/oraichain/orai/releases/tag/v0.42.0 to view the CHANGELOG for this upgrade" --from wallet --upgrade-height 21627705 --upgrade-info "https://github.com/oraichain/orai/releases/download/v0.42.0/manifest.json" --deposit 10000000orai --chain-id Oraichain -y -b block --gas-prices 0.001orai --gas 20000000 --node https://rpc.orai.io:443
```