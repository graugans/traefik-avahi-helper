<p align="center">
  <img src="https://img.shields.io/github/v/tag/graugans/traefik-avahi-helper" alt="GitHub tag (with filter)">
  <img src="https://img.shields.io/badge/License-Apache_2.0-blue.svg" alt="The Apache-2 License">
  <img src="https://github.com/graugans/traefik-avahi-helper/actions/workflows/commit.yml/badge.svg" alt="Build Status">
</p>

# Introduction

When using the reverse proxy [Træfik](https://traefik.io) in a local setup one may have the need to create a custom Avahi CNAME like `traefik.local` automatically when a container is created and routed by traefik. This lightweight helper container provides a Go binary which handles the communication to both the Docker daemon and the avahi daemon.

## Motivation

At home I assume most of the people have some sort of a broadband router which handles DHCP and DNS. If you are like me and do not want to mess around with DNS and run your own DNS server inside your home network than [mDNS or zeroconf](https://en.wikipedia.org/wiki/Zero-configuration_networking#Avahi) may come handy.

When working with the reverse Proxy [Træfik](https://traefik.io) most of the configuration options can be performed by Docker Labels, why not automatically register the mDNS entries based on those labels.

There is a great article by [Andrew Dupont](https://andrewdupont.net/2022/01/27/using-mdns-aliases-within-your-home-network/) about why and how to announce mDNS addresses.

## Alternatives

There are a couple of alternatives to this approach

- **[go-avahi-cname](https://github.com/grishy/go-avahi-cname)**, the idea of this project is to register specific CNAMEs or subdomains based on a **F**ull **Q**ualified **D**omain **N**ame (fqdn). This comes handy when you are looking for `*.<anyname>.local`. In this case no access to the Docker socket is required only access to the DBus socket.

- **[hardillb/traefik-avahi-helper](https://github.com/hardillb/traefik-avahi-helper)**, this is basically the same approach as this project. It uses a NodeJS and Python mix to achieve the goal.


## Project status

This is a WIP so please do not expect a working solution.

## LICENSE
This software is released under the Apache-2 license. See the [LICENSE](LICENSE) file for details. Some of the resources are released under the MIT License. Please check the related files.