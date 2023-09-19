# syntax=docker/dockerfile:1
ARG ALPINE_VERSION=3.18
ARG GO_VERSION=1.21

FROM --platform=$BUILDPLATFORM golang:${GO_VERSION}-alpine${ALPINE_VERSION} AS builder

RUN apk update && apk add --no-cache ca-certificates git

WORKDIR /build
COPY [ "go.mod", "go.sum", "./" ]
RUN go mod download
RUN go mod verify

COPY [ ".", "." ]
RUN go test ./...
ARG TARGETOS TARGETARCH
RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build -ldflags='-w -s' -o 'dist/terraform-version-inspect'

FROM scratch
COPY --from=builder [ "/build/dist/terraform-version-inspect", "/terraform-version-inspect" ]
COPY --from=builder [ "/etc/ssl/certs/ca-certificates.crt", "/etc/ssl/certs/" ]
WORKDIR /mnt
ENTRYPOINT [ "/terraform-version-inspect" ]
CMD [ "--help" ]
