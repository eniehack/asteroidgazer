FROM golang:1.14.0-alpine AS build

WORKDIR /opt/asteroidgazer

COPY . /opt/asteroidgazer
RUN apk add -U --no-cache git && \
    go build -o bin/asteroidgazer ./cmd/http

FROM alpine

LABEL maintainer="eniehack@outlook.jp" \
    org.opencontainers.image.title="asteroidgazer" \
    org.opencontainers.image.description="Asteroidgazer - ActivityPub relay also be able to use as search engine server." \
    org.opancontainers.image.licenses="AGPL-3.0"

COPY --from=build /opt/asteroidgazer/bin/asteroidgazer /usr/local/bin/asteroidgazer
RUN chmod +x /usr/local/bin/asteroidgazer && \
    apk add -U --no-cache ca-certificates
ENTRYPOINT ["/usr/local/bin/asteroidgazer"]
