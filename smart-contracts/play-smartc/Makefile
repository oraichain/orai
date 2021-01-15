

build:
	RUSTFLAGS='-C link-arg=-s' cargo wasm
	du -h ./target/wasm32-unknown-unknown/release/play_smartc.wasm | awk '{print $$1}'
