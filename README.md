gefyra - a redis bridge for channels and lists
============================================

* Forwards messages from one (or more) channels/lists to one (or more) other channels/lists
* Channels/lists can be on same or different Redis servers
* Simple configuration
* Lightweight, efficient, written in Go

## Quickstart

You can install the native binary for your operating system or use the Docker container.

> See [example](./examples) configurations.

### Option 1: Install using Homebrew

Install:

    brew tap outofcoffee/gefyra
    brew install gefyra

Run:

    export BRIDGE_CONFIG=./examples/channel_fan_out.yaml gefyra

### Option 2: Run using Docker

    docker run --rm -it -v /path/to/config/dir:/opt/gefyra/config outofcoffee/gefyra

## Supported modes

The following objects can be bridged:

| Upstream object | Downstream object | Same server? | Different server? |
|-----------------|-------------------|--------------|-------------------|
| pubsub channel  | pubsub channel    | Yes          | Yes               |
| pubsub channel  | list              | Yes          | Yes               |
| list            | pubsub channel    | Yes          | Yes               |
| list            | list              | Yes          | Yes               |

---

## Build and run

Build and run:

    docker-compose up --build

> For other architectures, see the [multiarchitecture](./docs/multiarch.md) documentation.

You should see something similar to the following:

```bash
loading config file: /etc/gefyra/config/one_up_one_down.yaml
loaded 1 upstream(s) and 1 downstream(s)
connecting to upstreams
initialising redis connection to upstream:6379
checking redis connection [attempt #1]...
redis connected at upstream:6379
connecting to downstreams
initialising redis connection to downstream:6379
checking redis connection [attempt #1]...
redis connected at downstream:6379
starting bridge upstream:6379/chan->downstream:6379/chan
```

## Test

If you've started the example using Docker Compose (above), then you can test by sending a message to the upstream and seeing it propagate downstream.

Subscribe to downstream:

    docker run --rm -it --network gefyra_redis bitnami/redis:4.0 redis-cli \
      -h downstream \
      -a password \
      SUBSCRIBE output-chan

Publish from upstream:

    docker run --rm -it --network gefyra_redis bitnami/redis:4.0 redis-cli \
      -h upstream \
      -a password \
      PUBLISH input-chan hello

---

_gefyra, γέφυρα, is Greek for 'bridge'._
