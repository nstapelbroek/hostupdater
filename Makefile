.PHONY: build

PROJECTNAME=nstapelbroek/hostupdater
TAGNAME=UNDEF
TAGNAME_CLEAN:=$(subst /,-,$(TAGNAME))

build:
	if [ "$(TAGNAME)" = "UNDEF" ]; then echo "please provide a valid TAGNAME" && exit 1; fi
	docker build --tag $(PROJECTNAME):$(TAGNAME_CLEAN) --pull .
