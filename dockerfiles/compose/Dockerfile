ARG ALPINE_VERSION=3.16
ARG GO_VERSION=1.19
ARG XCPU_VERSION=v0.6.0
ARG BUILDPLATFORM=linux/amd64

FROM --platform=${BUILDPLATFORM} qmcgaw/xcputranslate:${XCPU_VERSION} AS xcputranslate

FROM --platform=${BUILDPLATFORM} golang:${GO_VERSION}-alpine${ALPINE_VERSION} AS build
ENV CGO_ENABLED=0
RUN apk add --no-cache git && \
    git config --global advice.detachedHead false
COPY --from=xcputranslate /xcputranslate /usr/local/bin/xcputranslate
WORKDIR /tmp/build

ARG VERSION=v2.14.2
WORKDIR /go/src/github.com/docker/cli
RUN git clone --depth 1 --branch ${VERSION} https://github.com/docker/compose.git .
RUN go mod download
ARG TARGETPLATFORM
RUN echo "Building for ${TARGETPLATFORM}..." && \
    GOARCH="$(xcputranslate translate -field arch -targetplatform ${TARGETPLATFORM})" \
    GOARM="$(xcputranslate translate -field arm -targetplatform ${TARGETPLATFORM})" \
    go build -trimpath -ldflags="-s -w \
    -X 'github.com/docker/compose/v2/internal.Version=${VERSION}' \
    " -o /tmp/bin ./cmd && \
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
    org.opencontainers.image.url="https://github.com/qdm12/binpot/blob/main/dockerfiles/compose" \
    org.opencontainers.image.documentation="https://github.com/qdm12/binpot" \
    org.opencontainers.image.source="https://github.com/qdm12/binpot.git" \
    org.opencontainers.image.title="Docker Compose v2 plugin binary" \
    org.opencontainers.image.description="Docker Compose v2 plugin binary"
COPY --from=build /tmp/bin /bin
