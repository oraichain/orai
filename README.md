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

# steps to broadcast a price request
fn initScript <account-name> # initiate a complete price script
fn unsign <account-name>
fn sign <account-name>
fn broadcast

# start after init for development
make watch-oraid

```
