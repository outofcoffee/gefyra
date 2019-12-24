Redis Queue Bridge
==================

Build and run:

    docker-compose up --build

## Test

Subscribe to downstream:

    docker run --rm -it --network redis-queue-bridge_redis bitnami/redis:4.0 redis-cli \
      -h downstream \
      -a password \
      SUBSCRIBE chan

Publish from upstream:

    docker run --rm -it --network redis-queue-bridge_redis bitnami/redis:4.0 redis-cli \
      -h upstream \
      -a password \
      PUBLISH chan hello
