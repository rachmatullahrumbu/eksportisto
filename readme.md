# Eksportisto

Eksportisto (meaning 'exporter' in Esperanto) is a lightweight Celo blockchain parser we've built for internal and external use at cLabs. It will print all transactions (take a look in the [monitor directory](./monitor) to see exactly how it parses these) to standard out and additionally exposes Prometheus compatible metrics on port 8080.

Eksportisto uses SQLite to keep track of the last block parsed, so it is safe to start and stop without having to reparse the whole chain.

At cLabs we often rely on (Google's Operations (formerly Stackdriver))[https://cloud.google.com/products/operations] to collect these standard out logs and derive insights.

## How do I use it?

### Running a full node

We'd recommend running a Celo full node on the same network as Eksportisto. Taks a look at our [documentation for running a full node](https://docs.celo.org/getting-started/mainnet/running-a-full-node-in-mainnet) if you haven't already.

In addition to the steps in the above guide, you'll also need to make sure you run your full node with the following command line arguments:

- `--ws`
- `--wsapi eth,net,web3,debug`
- `--wsaddr 0.0.0.0`
- `--gcmode archive`

### Starting up Eksportisto

To start parsing blocks it should be as simple as running `go run main.go` from the root of this repository. We also maintain a Dockerfile if you want to run Eksportisto in a container.

The command line parameters most relevant to getting started quickly are:

- `-nodeUri (default ws://localhost:8546)` use this to point at your running full node
- `-datadir (default $HOME/.eksportisto)` where the Sqlite data directory will be stored. This is especially relevant if you choose to run Eksportisto in a Docker container and want to mount the same directory every time
- `-sensitiveAccounts` allows passing a JSON file of addresses->url entries. Whenever a transfer is initiated from one of these addresses a webhook will be sent with the payload of the transaction. It's important to note that these webhooks will only fire when in `tipMode`, or when Eksportisto has caught up to the tip of the chain and is reading blocks as they come.

More information can be found by running `go run main.go --help`.
