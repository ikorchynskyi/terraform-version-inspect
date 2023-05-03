FROM --platform=$BUILDPLATFORM golang:1.20-alpine3.17 as builder

RUN apk update && apk add --no-cache ca-certificates git

ARG TARGETOS TARGETARCH

WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
RUN go mod verify

COPY . .
RUN go test ./...
RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build -ldflags='-w -s' -o 'dist/terraform-version-inspect'

FROM scratch
COPY --from=builder /build/dist/terraform-version-inspect /terraform-version-inspect
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
WORKDIR /mnt
ENTRYPOINT [ "/terraform-version-inspect" ]
CMD [ "--help" ]
