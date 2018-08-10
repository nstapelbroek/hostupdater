# Hostupdater

[![Build Status](https://api.cirrus-ci.com/github/nstapelbroek/hostupdater.svg)](https://cirrus-ci.com/github/nstapelbroek/hostupdater)
[![Image details](https://images.microbadger.com/badges/image/nstapelbroek/hostupdater.svg)](https://hub.docker.com/r/nstapelbroek/hostupdater/tags/)
[![Image git ref](https://images.microbadger.com/badges/commit/nstapelbroek/hostupdater.svg)](https://microbadger.com/images/nstapelbroek/hostupdater)

> Use the hostupdater to patch your hostfile with the routing information collected form popular loadbalancers like traefik.

![traefik example](./docs/traefik-example.gif)

## Usage

### Traefik

The hostupdater can read the frontend rules from a [Traefik](https://traefik.io/) load balancer by using its API.
Make sure the Traefik API is enabled by using `--api` and point your hostupdater towards the endpoint.

An inline example:
```sh
# Start a traefik container that auto-registers new containers by listening to the docker socket
docker run -d -p 80:80 -p 8080:8080 --name traefik -v /var/run/docker.sock:/var/run/docker.sock:ro traefik --docker --api

# Start the hostupdater to persist routing information from traefik into your local hostfile every 7 seconds
export TRAEFIK_IP=$(docker inspect -f '{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' traefik)
docker run -d --name hostupdater -v /etc/hosts:/etc/hosts nstapelbroek/hostupdater traefik --interval 7 --address $TRAEFIK_IP

# To proof that it works, launch a ghost blog container
docker run -d --label traefik.frontend.rule=Host:ghost.local --name testdrive-ghost-blog ghost && echo "you can now visit http://ghost.local"
```

Hostupdater can really help when using multiple docker-compose stacks. Embed a hostupdater and traefik service to your docker-compose to set yourself free of all published port issues.

A docker-compose example:
```
version: '2'

services:
  traefik:
    image: traefik:alpine
    command: --api --docker
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock:ro

  hostupdater:
    image: docker.io/nstapelbroek/hostupdater:latest
    command: traefik --address traefik --wait 5
    depends_on:
      - "traefik"
    volumes:
      - /etc/hosts:/etc/hosts

  app:
    build:
      dockerfile: ./Dockerfile
      context: .
    volumes:
      - .:/www
    depends_on:
      - mariadb
    labels:
      - 'traefik.frontend.rule=Host:local.myproject.com'

  mariadb:
    image: mariadb:10.2
    environment:
      MYSQL_DATABASE: project
      MYSQL_USER: project_user
      MYSQL_PASSWORD: changemeplease

  mailhog:
    image: mailhog/mailhog:latest
    labels:
      - 'traefik.frontend.rule=Host:mailhog.local.myproject.com'
      - 'traefik.frontend.port=8025'
```

After a `docker-compose up` you can visit http://local.myproject.com and http://mailhog.local.myproject.com while the services can keep on using the integrated DNS resolver of docker to connect within the network.


