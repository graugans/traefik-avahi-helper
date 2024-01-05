# Introduction

This file is based on the [basic-example](https://doc.traefik.io/traefik/user-guides/docker-compose/basic-example/) from the Traefic user guide.

## TODO

- Replace `whoami.local` by your own domain within the traefik.`http.routers.whoami.rule` label of the whoami service.
- Run `docker-compose up -d` within the folder where you created the previous file.
- Wait a bit and visit `http://your_own_host.local` to confirm everything went fine. You should see the output of