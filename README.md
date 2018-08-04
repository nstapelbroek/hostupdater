# Hostupdater

[![Build Status](https://api.cirrus-ci.com/github/nstapelbroek/hostupdater.svg)](https://cirrus-ci.com/github/nstapelbroek/hostupdater)
[![Image details](https://images.microbadger.com/badges/image/nstapelbroek/hostupdater.svg)](https://hub.docker.com/r/nstapelbroek/hostupdater/tags/)
[![Image git ref](https://images.microbadger.com/badges/commit/nstapelbroek/hostupdater.svg)](https://microbadger.com/images/nstapelbroek/hostupdater)

> Use the hostupdater to patch your hostfile with the routing information collected form popular local development solutions like minicube, docker-compose and traefik.

![traefik example](./docs/traefik-example.gif)

## Usage

### Traefik

The hostupdater can read the frontend rules from a [Traefik](https://traefik.io/) load balancer by using its API.
Make sure the Traefik API is enabled by using `--api` and point your hostupdater towards the endpoint.

An example:
```sh
# Start a traefik container that auto-registers new containers by listening to the docker socket
docker run -d -p 80:80 -p 8080:8080 --name traefik -v /var/run/docker.sock:/var/run/docker.sock:ro traefik --docker --api

# Start the hostupdater to persist routing information from traefik into your local hostfile every 7 seconds
export TRAEFIK_IP=$(docker inspect -f '{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' traefik)
docker run -d --name hostupdater -v /etc/hosts:/etc/hosts nstapelbroek/hostupdater traefik --interval 7 --address $TRAEFIK_IP

# To proof that it works, launch a ghost blog container
docker run -d --label traefik.frontend.rule=Host:ghost.local --name testdrive-ghost-blog ghost && echo "you can now visit http://ghost.local"
```
