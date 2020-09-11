# P2P Chat

## Frontend

### Run

```sh
cd frontend
REACT_APP_SERVER="ws://3.106.131.208/connect" yarn start
```

Where `ws://3.106.131.208/connect` can be changed to represent the server's public ip
or include localhost if you wish to run the peer yourself.

### Build

```sh
cd frontend
REACT_APP_SERVER="ws://3.106.131.208/connect" yarn build 
```

## Backend

### Setup

```sh
go build
```

### Usage

```sh
./chat -nickname Hiru \
-bootstrap /ip4/3.106.131.208/tcp/35491/p2p/QmdKktf8LkcqoWGnoFmg9yz8Wyr2HdhRYhKssnxq8gtrqM \
-bootstrap /ip4/52.65.56.175/tcp/35431/p2p/QmbRQw7pHuPpGnHqkzqZDAbUnqaiagABHfaiGZaNZW6w62 \
-bootstrap /ip4/3.25.68.133/tcp/44091/p2p/QmRt9CfTKBJKXdTPn3SGXNa2meVrgmFnMmajpgh43YaWn2 \
-bootstrap /ip4/52.64.136.38/tcp/44477/p2p/QmcoUbLKsgMmx2kerWxWW54sNL89kRMYHRqJGXN6wG2zhK
```
