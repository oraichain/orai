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
docker build -t <image-tag> -f orai/Dockerfile --build-arg WASMVM_VERSION=v1.2.4 --build-arg VERSION=v0.41.5 .

# prod
docker build -t <image-tag> -f orai/Dockerfile.prod --build-arg WASMVM_VERSION=v1.2.4 --build-arg VERSION=v0.41.5 .
```