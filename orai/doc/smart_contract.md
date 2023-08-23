# How to simulate the oracle script smart contract

## NOTE: You need to be in the root directory of the repository to successfully follow the below steps. 

### 1. Build the smart contract

Start the rust container

```bash
docker-compose up -d rust && docker-compose exec rust bash
```

```bash
cd oscript-price/ && optimize.sh . && exit
```

### 2. Start and enter the docker container

```bash
docker-compose up -d simulate && docker-compose exec simulate bash
```

### 3. Simulate the smart contract

```bash
cosmwasm-simulate oscript-price/artifacts/oscript_price.wasm
```

### 4. Type "query" in the command line

The terminal will show something similar to:

```bash
Auto loading json schema from oscript-price/artifacts/schema
Input call type(init | handle | query):
```

Please type **query** to continue

### 5. Type the input for the query

The terminal will ask you to type an input json string, similar to this below:

```bash
Input json string:
```

Type the following input to query:

```bash
{"aggregate":{"results":["100.5","200.1","300.1"]}}
```

The **results** field requires an array of float numbers to aggregate the results

After typing, the result should look like:

```bash
executing func [query] , params is {"aggregate":{"results":["100.5","200.1","300.1"]}}
query data = "200.23"
Gas used   : 34661
***************************call finished***************************
```

The simulation process is done