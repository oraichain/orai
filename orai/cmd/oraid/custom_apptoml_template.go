package main

const CustomAppTomlTemplateWithEvm = `# This is a TOML config file.
# For more information, see https://github.com/toml-lang/toml

###############################################################################
###                           Base Configuration                            ###
###############################################################################

# The minimum gas prices a validator is willing to accept for processing a
# transaction. A transaction's fees must meet the minimum of any denomination
# specified in this config (e.g. 0.25token1;0.0001token2).
minimum-gas-prices = "{{ .BaseConfig.MinGasPrices }}"

# default: the last 100 states are kept in addition to every 500th state; pruning at 10 block intervals
# nothing: all historic states will be saved, nothing will be deleted (i.e. archiving node)
# everything: all saved states will be deleted, storing only the current and previous state; pruning at 10 block intervals
# custom: allow pruning options to be manually specified through 'pruning-keep-recent', 'pruning-keep-every', and 'pruning-interval'
pruning = "{{ .BaseConfig.Pruning }}"

# These are applied if and only if the pruning strategy is custom.
pruning-keep-recent = "{{ .BaseConfig.PruningKeepRecent }}"
pruning-keep-every = "{{ .BaseConfig.PruningKeepEvery }}"
pruning-interval = "{{ .BaseConfig.PruningInterval }}"

# HaltHeight contains a non-zero block height at which a node will gracefully
# halt and shutdown that can be used to assist upgrades and testing.
#
# Note: Commitment of state will be attempted on the corresponding block.
halt-height = {{ .BaseConfig.HaltHeight }}

# HaltTime contains a non-zero minimum block time (in Unix seconds) at which
# a node will gracefully halt and shutdown that can be used to assist upgrades
# and testing.
#
# Note: Commitment of state will be attempted on the corresponding block.
halt-time = {{ .BaseConfig.HaltTime }}

# MinRetainBlocks defines the minimum block height offset from the current
# block being committed, such that all blocks past this offset are pruned
# from Tendermint. It is used as part of the process of determining the
# ResponseCommit.RetainHeight value during ABCI Commit. A value of 0 indicates
# that no blocks should be pruned.
#
# This configuration value is only responsible for pruning Tendermint blocks.
# It has no bearing on application state pruning which is determined by the
# "pruning-*" configurations.
#
# Note: Tendermint block pruning is dependant on this parameter in conunction
# with the unbonding (safety threshold) period, state pruning and state sync
# snapshot parameters to determine the correct minimum value of
# ResponseCommit.RetainHeight.
min-retain-blocks = {{ .BaseConfig.MinRetainBlocks }}

# InterBlockCache enables inter-block caching.
inter-block-cache = {{ .BaseConfig.InterBlockCache }}

# IndexEvents defines the set of events in the form {eventType}.{attributeKey},
# which informs Tendermint what to index. If empty, all events will be indexed.
#
# Example:
# ["message.sender", "message.recipient"]
index-events = {{ .BaseConfig.IndexEvents }}

# IavlCacheSize set the size of the iavl tree cache. 
# Default cache size is 50mb.
iavl-cache-size = {{ .BaseConfig.IAVLCacheSize }}

# IAVLDisableFastNode enables or disables the fast node feature of IAVL. 
# Default is true.
iavl-disable-fastnode = {{ .BaseConfig.IAVLDisableFastNode }}

###############################################################################
###                         Telemetry Configuration                         ###
###############################################################################

[telemetry]

# Prefixed with keys to separate services.
service-name = "{{ .Telemetry.ServiceName }}"

# Enabled enables the application telemetry functionality. When enabled,
# an in-memory sink is also enabled by default. Operators may also enabled
# other sinks such as Prometheus.
enabled = {{ .Telemetry.Enabled }}

# Enable prefixing gauge values with hostname.
enable-hostname = {{ .Telemetry.EnableHostname }}

# Enable adding hostname to labels.
enable-hostname-label = {{ .Telemetry.EnableHostnameLabel }}

# Enable adding service to labels.
enable-service-label = {{ .Telemetry.EnableServiceLabel }}

# PrometheusRetentionTime, when positive, enables a Prometheus metrics sink.
prometheus-retention-time = {{ .Telemetry.PrometheusRetentionTime }}

# GlobalLabels defines a global set of name/value label tuples applied to all
# metrics emitted using the wrapper functions defined in telemetry package.
#
# Example:
# [["chain_id", "cosmoshub-1"]]
global-labels = [{{ range $k, $v := .Telemetry.GlobalLabels }}
  ["{{index $v 0 }}", "{{ index $v 1}}"],{{ end }}
]

###############################################################################
###                           API Configuration                             ###
###############################################################################

[api]

# Enable defines if the API server should be enabled.
enable = {{ .API.Enable }}

# Swagger defines if swagger documentation should automatically be registered.
swagger = {{ .API.Swagger }}

# Address defines the API server to listen on.
address = "{{ .API.Address }}"

# MaxOpenConnections defines the number of maximum open connections.
max-open-connections = {{ .API.MaxOpenConnections }}

# RPCReadTimeout defines the Tendermint RPC read timeout (in seconds).
rpc-read-timeout = {{ .API.RPCReadTimeout }}

# RPCWriteTimeout defines the Tendermint RPC write timeout (in seconds).
rpc-write-timeout = {{ .API.RPCWriteTimeout }}

# RPCMaxBodyBytes defines the Tendermint maximum response body (in bytes).
rpc-max-body-bytes = {{ .API.RPCMaxBodyBytes }}

# EnableUnsafeCORS defines if CORS should be enabled (unsafe - use it at your own risk).
enabled-unsafe-cors = {{ .API.EnableUnsafeCORS }}

###############################################################################
###                           Rosetta Configuration                         ###
###############################################################################

[rosetta]

# Enable defines if the Rosetta API server should be enabled.
enable = {{ .Rosetta.Enable }}

# Address defines the Rosetta API server to listen on.
address = "{{ .Rosetta.Address }}"

# Network defines the name of the blockchain that will be returned by Rosetta.
blockchain = "{{ .Rosetta.Blockchain }}"

# Network defines the name of the network that will be returned by Rosetta.
network = "{{ .Rosetta.Network }}"

# Retries defines the number of retries when connecting to the node before failing.
retries = {{ .Rosetta.Retries }}

# Offline defines if Rosetta server should run in offline mode.
offline = {{ .Rosetta.Offline }}

###############################################################################
###                           gRPC Configuration                            ###
###############################################################################

[grpc]

# Enable defines if the gRPC server should be enabled.
enable = {{ .GRPC.Enable }}

# Address defines the gRPC server address to bind to.
address = "{{ .GRPC.Address }}"

###############################################################################
###                        gRPC Web Configuration                           ###
###############################################################################

[grpc-web]

# GRPCWebEnable defines if the gRPC-web should be enabled.
# NOTE: gRPC must also be enabled, otherwise, this configuration is a no-op.
enable = {{ .GRPCWeb.Enable }}

# Address defines the gRPC-web server address to bind to.
address = "{{ .GRPCWeb.Address }}"

# EnableUnsafeCORS defines if CORS should be enabled (unsafe - use it at your own risk).
enable-unsafe-cors = {{ .GRPCWeb.EnableUnsafeCORS }}

###############################################################################
###                        State Sync Configuration                         ###
###############################################################################

# State sync snapshots allow other nodes to rapidly join the network without replaying historical
# blocks, instead downloading and applying a snapshot of the application state at a given height.
[state-sync]

# snapshot-interval specifies the block interval at which local state sync snapshots are
# taken (0 to disable). Must be a multiple of pruning-keep-every.
snapshot-interval = {{ .StateSync.SnapshotInterval }}

# snapshot-keep-recent specifies the number of recent snapshots to keep and serve (0 to keep all).
snapshot-keep-recent = {{ .StateSync.SnapshotKeepRecent }}

###############################################################################
###                         Store / State Streaming                         ###
###############################################################################

[store]
streamers = [{{ range .Store.Streamers }}{{ printf "%q, " . }}{{end}}]

[streamers]
[streamers.file]
keys = [{{ range .Streamers.File.Keys }}{{ printf "%q, " . }}{{end}}]
write_dir = "{{ .Streamers.File.WriteDir }}"
prefix = "{{ .Streamers.File.Prefix }}"

# output-metadata specifies if output the metadata file which includes the abci request/responses 
# during processing the block.
output-metadata = "{{ .Streamers.File.OutputMetadata }}"

# stop-node-on-error specifies if propagate the file streamer errors to consensus state machine.
stop-node-on-error = "{{ .Streamers.File.StopNodeOnError }}"

# fsync specifies if call fsync after writing the files.
fsync = "{{ .Streamers.File.Fsync }}"

###############################################################################
###                             EVM Configuration                           ###
###############################################################################

[evm]

# Tracer defines the 'vm.Tracer' type that the EVM will use when the node is run in
# debug mode. To enable tracing use the '--evm.tracer' flag when starting your node.
# Valid types are: json|struct|access_list|markdown
tracer = ""

# MaxTxGasWanted defines the gas wanted for each eth tx returned in ante handler in check tx mode.
max-tx-gas-wanted = 500000

###############################################################################
###                           JSON RPC Configuration                        ###
###############################################################################

[json-rpc]

# Enable defines if the gRPC server should be enabled.
enable = true

# Address defines the EVM RPC HTTP server address to bind to.
address = "127.0.0.1:8545"

# Address defines the EVM WebSocket server address to bind to.
ws-address = "127.0.0.1:8546"

# API defines a list of JSON-RPC namespaces that should be enabled
# Example: "eth,txpool,personal,net,debug,web3"
api = "eth,net,web3"

# GasCap sets a cap on gas that can be used in eth_call/estimateGas (0=infinite). Default: 25,000,000.
gas-cap = 25000000

# EVMTimeout is the global timeout for eth_call. Default: 5s.
evm-timeout = "5s"

# TxFeeCap is the global tx-fee cap for send transaction. Default: 1eth.
txfee-cap = 1

# FilterCap sets the global cap for total number of filters that can be created
filter-cap = 200

# FeeHistoryCap sets the global cap for total number of blocks that can be fetched
feehistory-cap = 100

# LogsCap defines the max number of results can be returned from single 'eth_getLogs' query.
logs-cap = 10000

# BlockRangeCap defines the max block range allowed for 'eth_getLogs' query.
block-range-cap = 10000

# HTTPTimeout is the read/write timeout of http json-rpc server.
http-timeout = "30s"

# HTTPIdleTimeout is the idle timeout of http json-rpc server.
http-idle-timeout = "2m0s"

###############################################################################
###                             TLS Configuration                           ###
###############################################################################

[tls]

# Certificate path defines the cert.pem file path for the TLS configuration.
certificate-path = ""

# Key path defines the key.pem file path for the TLS configuration.
key-path = ""
`