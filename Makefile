.PHONY: build

PROJECTNAME=nstapelbroek/hostupdater
TAGNAME=UNDEF
TAGNAME_CLEAN:=$(subst /,-,$(TAGNAME))
PROJECT_PACKAGES=$(shell go list ./... | grep -v vendor)

build:
	if [ "$(TAGNAME)" = "UNDEF" ]; then echo "please provide a valid TAGNAME" && exit 1; fi
	CGO_ENABLED=0 GOOS=linux go build  -ldflags '-w -s' -a -installsuffix cgo -o hostupdater .
	docker build --tag $(PROJECTNAME):$(TAGNAME_CLEAN) --pull .
	rm hostupdater