ARG ALPINE_VERSION=3.15
ARG GO_VERSION=1.18
ARG XCPU_VERSION=v0.6.0
ARG BUILDPLATFORM=linux/amd64

FROM --platform=${BUILDPLATFORM} qmcgaw/xcputranslate:${XCPU_VERSION} AS xcputranslate

FROM --platform=${BUILDPLATFORM} golang:${GO_VERSION}-alpine${ALPINE_VERSION} AS build
ENV CGO_ENABLED=0
RUN apk add --no-cache git && \
    git config --global advice.detachedHead false
COPY --from=xcputranslate /xcputranslate /usr/local/bin/xcputranslate
WORKDIR /tmp/build

ARG VERSION=v0.8.1
WORKDIR /go/src/github.com/docker/cli
RUN git clone --single-branch --branch ${VERSION} https://github.com/docker/buildx.git .
RUN go mod download
ARG TARGETPLATFORM
RUN echo "Building for ${TARGETPLATFORM}..." && \
    export GOARCH="$(xcputranslate translate -field arch -targetplatform ${TARGETPLATFORM})" && \
    export GOARM="$(xcputranslate translate -field arm -targetplatform ${TARGETPLATFORM})" && \
    go build -trimpath -ldflags="-s -w \
    -X 'github.com/docker/buildx/version.Version=${VERSION}' \
    -X 'github.com/docker/buildx/version.Revision=$(git rev-parse HEAD)' \
    -X 'github.com/docker/buildx/version.Package=github.com/docker/buildx' \
    " -o /tmp/bin ./cmd/buildx && \
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
    org.opencontainers.image.url="https://github.com/qdm12/binpot/blob/main/dockerfiles/buildx" \
    org.opencontainers.image.documentation="https://github.com/qdm12/binpot" \
    org.opencontainers.image.source="https://github.com/qdm12/binpot.git" \
    org.opencontainers.image.title="Docker Buildx plugin binary" \
    org.opencontainers.image.description="Docker Buildx plugin binary"
COPY --from=build /tmp/bin /bin
