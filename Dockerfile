FROM alpine:3.8

RUN apk add --update --no-cache ca-certificates

ADD hostupdater /

ENTRYPOINT ["/hostupdater"]