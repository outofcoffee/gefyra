# Building for multiple architectures

## Default (amd64) build

    docker-compose build

## Other architectures

Set the `GOARCH` build argument and, optionally, the `IMAGE_TAG` environment variable.

For example, for arm32v7, set as follows:

    IMAGE_TAG=latest-arm32v7 docker-compose build --build-arg GOARCH=arm

For aarch64, set as follows:

    IMAGE_TAG=latest-aarch64 docker-compose build --build-arg GOARCH=arm64
