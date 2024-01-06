# syntax=docker/dockerfile:1

# Copyright (c) 2024 Christian Ege <ch@ege.io>
# Copyright (c) 2023 Sergey G. <mailgrishy@gmail.com>
# SPDX-License-Identifier: MIT

# The first build stage to build the executable
FROM docker.io/golang:1.21.5-bookworm as builder
WORKDIR /build

# Download all dependencies
COPY go.mod go.sum ./
RUN go mod download

# We have to disable CGO to not depend on the glibc
# Otherwise we are not able to run on the `scratch` image
ENV CGO_ENABLED=0
COPY . .
# Build the binary
RUN go build -ldflags="-w -s" -o /traefik-avahi-helper


# The second stage use the image in a bare minimum image
FROM scratch
LABEL maintainer="Christian Ege ch@ege.io"
COPY --from=builder /traefik-avahi-helper /traefik-avahi-helper

# Expose the mDNS service port
EXPOSE 5353/udp

# Our executable
ENTRYPOINT ["/traefik-avahi-helper"]