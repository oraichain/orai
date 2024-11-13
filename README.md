# Oraichain

![Banner!](./doc/logo-full-h-light.png#gh-light-mode-only)
![Banner!](./doc/logo-full-h-dark.png#gh-dark-mode-only)

[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

> [!WARNING] 
> Oraichain Labs, the team behind Oraichain and this repo orai, has fully shifted its focus to [Our Wasmd fork repository](https://github.com/oraichain/wasmd), where we continue to build interoperability solutions for the future! 
>
> Please read our [Upgrade guideline](https://github.com/oraichain/wasmd/blob/release/v0.50.x/UPGRADE_V050.md) to understand the reasons for the change.
>
> As a result, this repository is no longer maintained, and all relevant code has been migrated.
>
> The last production version of this repo is [v0.42.4](https://github.com/oraichain/orai/releases/tag/v0.42.4).
>

Oraichain is the world’s first AI-powered oracle and ecosystem for blockchains. 

Beyond data oracles, Oraichain aims to become the first AI Layer 1 in the Blockchain sphere with a complete AI ecosystem, serving as a foundational layer for the creation of a new generation of smart contracts and Dapps. With AI as the cornerstone, Oraichain has developed many essential and innovative products and services including AI price feeds, fully on-chain VRF, Data Hub, AI Marketplace with 100+ AI APIs, AI-based NFT generation and NFT copyright protection, Royalty Protocol, AI-powered Yield Aggregator Platform, and Cosmwasm IDE.

This repository contains the source code & how to build the Oraichain mainnet, a Cosmos-based blockchain network that levarages the [CosmWasm](https://github.com/CosmWasm/cosmwasm) technology to integrate AI into the ecosystem.

## Getting Started

[These instructions](https://docs.orai.io/developers/tutorials/getting-setup) will get you a copy of the project up and running on your local machine for development and testing purposes. See [deployment](#deployment) for notes on how to deploy the project on a live system.

### Prerequisites

* If you want to build the binary using Docker (recommended), then you only need Docker.

* If you want to build the binary from scratch, you will need:

    - Go 1.15+

    - Make

    - Wasmvm library: https://github.com/CosmWasm/wasmvm/releases/download/v0.13.0/libwasmvm_muslc.a (you can download it and put in /lib/libwasmvm_muslc.a). The file is used by CosmWasm when building the binary

### Hardware requirements

[Please visit the official hardware requirement for Oraichain mainnet here](https://docs.orai.io/developers/networks/mainnet#node-hardwarde-specification)

### Installing

* **Install Golang**

[Please visit the official Golang website to download & install Go](https://go.dev/doc/install)

* **Install make**

Normally, for Linux-based machines, you already have Make installed by default.

* **Install libwasmvm**

the wasmd module of CosmWasm uses a wasm vm library, which should be included when building the chain binary. Hence, we need to download and place it in a specific location.

For Linux based machines, please run the following command:

```bash
sudo wget https://github.com/CosmWasm/wasmvm/releases/download/v0.13.0/libwasmvm_muslc.a -O /lib/libwasmvm_muslc.a
```

* **Download Go dependencies**

`go get ./...`

* **Build the binary**

`make install`

* **Verify the binary version**

`oraid version`

## Protobuf & protobuf swagger generation

* [Install Docker](https://docs.docker.com/engine/install)

* Start the proto docker: `docker-compose up -d proto`

* Install neccessary tools: `docker-compose exec proto ash -c 'apk add build-base bash && go get ./...'`

* Gen protobuf: `docker-compose exec proto ash -c 'make proto-gen'`

## Contributing

Please read [CONTRIBUTING.md](./CONTRIBUTING.md) for details on our code of conduct, and the process for submitting pull requests to us.

## Versioning

We use [SemVer](http://semver.org/) for versioning. For the versions available, see the [tags on this repository](https://github.com/oraichain/orai/tags). 

## Authors

* [Duc Le Pham](https://github.com/ducphamle2)

See also the list of [contributors](https://github.com/oraichain/orai/contributors) who participated in this project.

## License

This project is licensed under the Apache 2.0 license - see the [LICENSE](LICENSE) file for details.

<!-- ## Acknowledgments

* Hat tip to anyone whose code was used
* Inspiration
* etc -->

<!-- ## Run test
`make test-method PACKAGE=github.com/oraichain/orai/x/airequest/keeper METHOD=TestCalucateMol` -->

<!-- ## Create swagger documentation

```bash
# go to proto
docker-compose exec proto ash
make proto-swagger
# then create static file
go install github.com/rakyll/statik
statik -src doc/swagger-ui/ -dest doc -f
```

## Non-docker build

make build LEDGER_ENABLED=false GOMOD_FLAGS= VERSION=0.41.0
