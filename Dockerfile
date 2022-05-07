# Build
FROM golang:1.17.5-alpine3.15 AS base
COPY . /go/src/github.com/aschbacd/strava-export
WORKDIR /go/src/github.com/aschbacd/strava-export
RUN go build -a -tags netgo -ldflags '-w' -o /go/bin/strava-export /go/src/github.com/aschbacd/strava-export

# Package
FROM alpine:3.15.4
RUN apk update && apk add ca-certificates

COPY --from=base /go/bin/strava-export /usr/share/strava-export/strava-export
COPY ./assets /usr/share/strava-export/assets
COPY ./views /usr/share/strava-export/views

WORKDIR /usr/share/strava-export
ENTRYPOINT ["./strava-export"]
