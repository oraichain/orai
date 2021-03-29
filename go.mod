module github.com/oraichain/orai

go 1.15

require (
	github.com/CosmWasm/wasmd v0.15.0
	github.com/cosmos/cosmos-sdk v0.42.3
	github.com/cosmos/go-bip39 v1.0.0
	github.com/cosmtrek/air v1.21.2 // indirect
	github.com/creack/pty v1.1.11 // indirect
	github.com/fatih/color v1.10.0 // indirect
	github.com/gogo/protobuf v1.3.3
	github.com/golang/protobuf v1.4.3
	github.com/gorilla/mux v1.8.0
	github.com/grpc-ecosystem/grpc-gateway v1.16.0
	github.com/imdario/mergo v0.3.11 // indirect
	github.com/kyokomi/emoji v2.2.4+incompatible
	github.com/oasisprotocol/oasis-core/go v0.2012.4
	github.com/pelletier/go-toml v1.8.1 // indirect
	github.com/pkg/errors v0.9.1
	github.com/pwaller/goupx v0.0.0-20160623083017-1d58e01d5ce2 // indirect
	github.com/rakyll/statik v0.1.7
	github.com/rs/zerolog v1.20.0
	github.com/segmentio/ksuid v1.0.3
	github.com/snikch/goodman v0.0.0-20171125024755-10e37e294daa
	github.com/spf13/cast v1.3.1
	github.com/spf13/cobra v1.1.1
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.7.1
	github.com/stretchr/testify v1.7.0
	github.com/tendermint/tendermint v0.34.8
	github.com/tendermint/tm-db v0.6.4
	golang.org/x/sync v0.0.0-20201207232520-09787c993a3a
	golang.org/x/sys v0.0.0-20210220050731-9a76102bfb43 // indirect
	google.golang.org/genproto v0.0.0-20210114201628-6edceaf6022f
	google.golang.org/grpc v1.35.0
)

replace google.golang.org/grpc => google.golang.org/grpc v1.33.2

replace github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.3-alpha.regen.1
