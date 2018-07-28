.PHONY: build, run

PROJECTNAME=nstapelbroek/hostupdater
TAGNAME=UNDEF
TAGNAME_CLEAN:=$(subst /,-,$(TAGNAME))
PROJECT_PACKAGES=$(shell go list ./... | grep -v vendor)

build:
	if [ "$(TAGNAME)" = "UNDEF" ]; then echo "please provide a valid TAGNAME" && exit 1; fi
	CGO_ENABLED=0 GOOS=linux go build  -ldflags '-w -s' -a -installsuffix cgo -o hostupdater .
	docker build --tag $(PROJECTNAME):$(TAGNAME_CLEAN) --pull .
	rm hostupdater

run:
	if [ "$(TAGNAME)" = "UNDEF" ]; then echo "please provide a valid TAGNAME" && exit 1; fi
	docker run --rm --name hostupdater-run -v /etc/hosts:/etc/hosts $(PROJECTNAME):$(TAGNAME_CLEAN) traefik --address traefik