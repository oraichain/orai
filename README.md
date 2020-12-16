<br/>
<p align="center">
<a href="https://orai.io/" target="_blank" color="#0d2990">
  <img src="https://orai.io/images/favicon.png" width="225" alt="Oraichain logo">
</a>
</p>
<br/>

## Monolithic repository of everything Oraichain

## Starterkit

```bash
docker-compose up -d
docker-compose exec orai bash

# using fn to run on multiple terminal
fn oraid # start node
fn websocketInit # start websocket
fn restServer # start rest server

# use fn to run on one terminal
fn init && fn start

Note: fn init should only be called once. After that, when running, you can use fn start to keep the data. fn init will clean everything and start over.

# steps to broadcast a price request
fn initScript <account-name> # initiate a complete price script
fn unsign <account-name>
fn sign <account-name>

# start after init for development to automatically rebuild the binaries
make watch-oraid

# rebuild the binary manually after editing the source code
make all

Note: After `make all`, you need to re-run the fn commands to apply new binaries.

```
