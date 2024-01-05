[![License](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

# Overview

When using the reverse proxy [Træfik](https://traefik.io) in a local setup one may have the need to create custom Avahi cnames like `traefik.local` automatically when a container is created and routed by traefik. This helper container provides a Go binary which handles the communication to both the Docker daemon and the avahi daemon. 


## Project status

This is a WIP so please do not expect a working solution.

## LICENSE
This software is released under the Apache-2 license. See the [LICENSE](LICENSE) file for details.