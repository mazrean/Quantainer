# syntax = docker/dockerfile:1.3.0

FROM golang:1.17.5-alpine AS build

RUN apk add --update --no-cache git

WORKDIR /go/src/github.com/mazrean/quauntainer

RUN --mount=type=cache,target=/root/.cache/go-build \
  go install github.com/golang/mock/mockgen@v1.6.0
RUN --mount=type=cache,target=/root/.cache/go-build \
  go install github.com/google/wire/cmd/wire@v0.5.0

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod/cache \
  go mod download

COPY . .
RUN go generate ./...
RUN --mount=type=cache,target=/root/.cache/go-build \
  && go build -o quantainer -ldflags "-s -w"

FROM alpine:3.15.0

WORKDIR /go/src/github.com/mazrean/quauntainer

RUN apk --update --no-cache add tzdata \
  && cp /usr/share/zoneinfo/Asia/Tokyo /etc/localtime \
  && apk del tzdata \
  && mkdir -p /usr/share/zoneinfo/Asia \
  && ln -s /etc/localtime /usr/share/zoneinfo/Asia/Tokyo
RUN apk --update --no-cache add ca-certificates \
  && update-ca-certificates \
  && rm -rf /usr/share/ca-certificates

COPY --from=build /go/src/github.com/mazrean/quauntainer/quantainer ./quantainer

ENTRYPOINT ["./quantainer"]
