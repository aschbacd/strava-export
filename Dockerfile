# Build
FROM golang:1.17.5-alpine3.15 AS base
COPY . /go/src/github.com/aschbacd/strava-export
WORKDIR /go/src/github.com/aschbacd/strava-export
RUN go build -a -tags netgo -ldflags '-w' -o /go/bin/strava-export /go/src/github.com/aschbacd/strava-export

# Package
FROM scratch
COPY --from=base /go/bin/strava-export /strava-export
ENTRYPOINT ["/strava-export"]
