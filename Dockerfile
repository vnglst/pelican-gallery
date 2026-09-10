# syntax=docker/dockerfile:1.7

FROM --platform=$BUILDPLATFORM alpine:3.22 AS tailwind

ARG BUILDARCH
ARG TAILWIND_VERSION=v4.3.3

RUN apk add --no-cache curl \
    && case "$BUILDARCH" in \
        amd64) tailwind_arch="x64" ;; \
        arm64) tailwind_arch="arm64" ;; \
        *) echo "Unsupported build architecture: $BUILDARCH" >&2; exit 1 ;; \
    esac \
    && curl --fail --location --retry 3 \
        "https://github.com/tailwindlabs/tailwindcss/releases/download/${TAILWIND_VERSION}/tailwindcss-linux-${tailwind_arch}" \
        --output /tailwindcss \
    && chmod +x /tailwindcss

FROM --platform=$BUILDPLATFORM golang:1.24-bookworm AS build

ARG TARGETOS
ARG TARGETARCH

WORKDIR /src

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

COPY --from=tailwind /tailwindcss /usr/local/bin/tailwindcss
COPY static ./static
COPY templates ./templates
RUN tailwindcss -i ./static/css/input.css -o ./static/css/output.css --minify

COPY config ./config
COPY internal ./internal
COPY main.go ./
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build -trimpath -ldflags="-s -w" -o /server ./main.go

FROM alpine:3.22 AS runtime

RUN apk add --no-cache ca-certificates curl

WORKDIR /app

COPY --from=build /server ./server
COPY --from=build /src/static/css/output.css ./static/css/output.css
COPY config ./config
COPY artworks.db ./artworks.db

ENV GO_ENV=production
EXPOSE 8080

CMD ["./server"]
