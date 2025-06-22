ARG ALPINE_VERSION=3.22
ARG GO_VERSION=1.24
ARG XCPU_VERSION=v0.6.0
ARG BUILDPLATFORM=linux/amd64

FROM --platform=${BUILDPLATFORM} qmcgaw/xcputranslate:${XCPU_VERSION} AS xcputranslate

FROM --platform=${BUILDPLATFORM} golang:${GO_VERSION}-alpine${ALPINE_VERSION} AS build
ENV CGO_ENABLED=0
RUN apk add --no-cache git && \
    git config --global advice.detachedHead false
COPY --from=xcputranslate /xcputranslate /usr/local/bin/xcputranslate
WORKDIR /tmp/build

ARG VERSION=v2.1.6
RUN git clone --depth 1 --branch ${VERSION} https://github.com/golangci/golangci-lint.git .
RUN go mod download
ARG TARGETPLATFORM
ARG GO_VERSION=1.24
RUN echo "Building for ${TARGETPLATFORM}..." && \
    GOARCH="$(xcputranslate translate -field arch -targetplatform ${TARGETPLATFORM})" \
    GOARM="$(xcputranslate translate -field arm -targetplatform ${TARGETPLATFORM})" \
    go build -trimpath -ldflags="-s -w \
    -X 'main.goVersion=${GO_VERSION}' \
    -X 'main.version=${VERSION}' \
    -X 'main.commit=$(git rev-parse --short HEAD)' \
    -X 'main.date=$(date +%Y-%m-%dT%T%z)' \
    " -o /tmp/bin ./cmd/golangci-lint && \
    chmod 500 /tmp/bin

FROM scratch
ENTRYPOINT [ "/bin" ]
ARG VERSION
ARG CREATED
ARG REVISION
LABEL \
    org.opencontainers.image.authors="quentin.mcgaw@gmail.com" \
    org.opencontainers.image.version=${VERSION} \
    org.opencontainers.image.created=${CREATED} \
    org.opencontainers.image.revision=${REVISION} \
    org.opencontainers.image.url="https://github.com/qdm12/binpot/blob/main/dockerfiles/golangci-lint" \
    org.opencontainers.image.documentation="https://github.com/qdm12/binpot" \
    org.opencontainers.image.source="https://github.com/qdm12/binpot.git" \
    org.opencontainers.image.title="Golangci-lint binary" \
    org.opencontainers.image.description="Golangci-lint binary"
COPY --from=build /tmp/bin /bin
