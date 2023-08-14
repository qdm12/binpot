ARG ALPINE_VERSION=3.18
ARG GO_VERSION=1.21
ARG XCPU_VERSION=v0.6.0
ARG BUILDPLATFORM=linux/amd64

FROM --platform=${BUILDPLATFORM} qmcgaw/xcputranslate:${XCPU_VERSION} AS xcputranslate

FROM --platform=${BUILDPLATFORM} golang:${GO_VERSION}-alpine${ALPINE_VERSION} AS build
ENV CGO_ENABLED=0
RUN apk add --no-cache git && \
    git config --global advice.detachedHead false
COPY --from=xcputranslate /xcputranslate /usr/local/bin/xcputranslate
WORKDIR /tmp/build

ARG VERSION=v1.27.4
RUN git clone --depth 1 --branch ${VERSION} https://github.com/kubernetes/kubernetes.git .
RUN go mod download
ARG TARGETPLATFORM
RUN echo "Building for ${TARGETPLATFORM}..." && \
    SOURCE_DATE_EPOCH="$(git show -s --format=format:%ct HEAD)" && \
    BUILD_DATE="$(date ${SOURCE_DATE_EPOCH:+"--date=@${SOURCE_DATE_EPOCH}"} -u +'%Y-%m-%dT%H:%M:%SZ')" && \
    GITCOMMIT="$(git rev-parse HEAD)" && \
    GOARCH="$(xcputranslate translate -field arch -targetplatform ${TARGETPLATFORM})" \
    GOARM="$(xcputranslate translate -field arm -targetplatform ${TARGETPLATFORM})" \
    # Ldflags: see https://github.com/kubernetes/kubernetes/blob/ea0764452222146c47ec826977f49d7001b0ea8c/hack/lib/version.sh#L151
    go build -trimpath -ldflags="-s -w \
    -X 'k8s.io/client-go/pkg/version.buildDate=${BUILD_DATE}' \
    -X 'k8s.io/client-go/pkg/version.gitCommit=${GITCOMMIT}' \
    -X 'k8s.io/client-go/pkg/version.gitTreeState=clean' \
    -X 'k8s.io/client-go/pkg/version.gitVersion=${VERSION}' \
    -X 'k8s.io/client-go/pkg/version.gitMajor=1' \
    -X 'k8s.io/client-go/pkg/version.gitMinor=29' \
    -X 'k8s.io/component-base/version.buildDate=${BUILD_DATE}' \
    -X 'k8s.io/component-base/version.gitCommit=${GITCOMMIT}' \
    -X 'k8s.io/component-base/version.gitTreeState=clean' \
    -X 'k8s.io/component-base/version.gitVersion=${VERSION}' \
    -X 'k8s.io/component-base/version.gitMajor=1' \
    -X 'k8s.io/component-base/version.gitMinor=21' \
    " -o /tmp/bin cmd/kubectl/kubectl.go && \
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
    org.opencontainers.image.url="https://github.com/qdm12/binpot/blob/main/dockerfiles/kubectl" \
    org.opencontainers.image.documentation="https://github.com/qdm12/binpot" \
    org.opencontainers.image.source="https://github.com/qdm12/binpot.git" \
    org.opencontainers.image.title="Kubectl binary" \
    org.opencontainers.image.description="Kubectl binary"
COPY --from=build /tmp/bin /bin
