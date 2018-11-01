# Build environment
FROM golang:1.11 AS build-env
# GOPATH is /go
WORKDIR  /go/src/github.com/nstapelbroek/hostupdater 

COPY . .
RUN go get -u github.com/golang/dep/cmd/dep
RUN dep ensure
RUN CGO_ENABLED=0 GOOS=linux go build  -ldflags '-w -s' -a -installsuffix cgo -o hostupdater .

# Final container
FROM alpine:3.8

ARG VCS_REF
LABEL org.label-schema.vcs-ref=$VCS_REF \
      org.label-schema.vcs-url="https://github.com/nstapelbroek/hostupdater"

RUN apk add --update --no-cache ca-certificates
COPY --from=build-env /go/src/github.com/nstapelbroek/hostupdater/hostupdater /

ENTRYPOINT ["/hostupdater"]