# P2P Chat

## Frontend

The frontend is also hosted at [chat.hiru.dev](https://chat.hiru.dev).

### Run

```sh
cd frontend
yarn start
```

Then go into settings and set the peer you wish to connect to!

### Build

```sh
cd frontend
yarn build 
```

## Backend

### Setup

```sh
go build
```

### Usage

```sh
./chat -wsport 8000 \
-bootstrap /ip4/13.236.84.197/tcp/39589/p2p/Qmc9NUmrtY8ZdW9Tzd45obJ967n1HKtmEfKhdGtXEcWRzo \
-bootstrap /ip4/3.106.59.113/tcp/37063/p2p/QmafbLCxsoWB7Zk4z3rXMetpzXsKZhEA5PdqgQpDx8iKkz \
-bootstrap /ip4/3.106.53.149/tcp/34227/p2p/QmcfHPvy2VTPQsqstoDmtc5T6drW9cL9KtwxJrRWatY8y8 \
-bootstrap /ip4/54.252.218.141/tcp/39825/p2p/QmbffA9nGngJB6k6FVaSRUL5uGqHCCWL9DwoNsSdEGKR88
```
