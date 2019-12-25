gefyra - a Redis Queue Bridge
=============================

* Forwards messages from one (or more) Redis queues to one (or more) other queues
* Simple configuration
* Lightweight, efficient, written in Go

## Quickstart

    docker run --rm -it -v /path/to/config/dir:/opt/gefyra/config outofcoffee/gefyra

See [example](./examples) configurations.

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
      SUBSCRIBE chan

Publish from upstream:

    docker run --rm -it --network gefyra_redis bitnami/redis:4.0 redis-cli \
      -h upstream \
      -a password \
      PUBLISH chan hello

---

_gefyra, γέφυρα, is Greek for 'bridge'._
