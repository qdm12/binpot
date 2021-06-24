# Unavailable

Unavailable is a tiny, simple Go program to print out a program is unavailable for a certain platform.

For example:

```log
helm v3.6.1 is unavailable for linux/riscv64
```

It is designed to be built in a Dockerfile as an alternative to a program if this program is not supported for a certain architecture.

To reduce the complexity of the Dockerfile if/else block in a `RUN` instruction, the [`build.sh`](build.sh) script can be pipe to `sh` with:

```sh
wget -qO- https://raw.githubusercontent.com/qdm12/binpot/main/unavailable/build.sh | \
  sh -s -- "programName v9.9.9" "linux/riscv64"
```

The script requires:

- `xcputranslate` (see one of the [Dockerfiles](../dockerfiles) on how-to)
- `wget`
- `go`
- Argument 1 as the program name
- `${TARGETPLATFORM}` to be set

It will:

1. Clear the current working directory
1. Download the [main.go](main.go) file
1. Cross build the binary, setting the program name and platform in the binary with `-ldflags`
1. Write the binary to `/tmp/bin` with `500` permissions.

## Concrete example

You can see its usage in some of [the Dockerfiles](../dockerfiles) where not all architectures are supported by the program, such as for [dlv](../dockerfiles/dlv/Dockerfile).

Relevant section for `dlv`:

```Dockerfile
ARG TARGETPLATFORM
RUN export GOARCH="$(xcputranslate translate -field arch -targetplatform ${TARGETPLATFORM})" && \
    export GOARM="$(xcputranslate translate -field arm -targetplatform ${TARGETPLATFORM})" && \
    if [ "${TARGETPLATFORM}" = "linux/amd64" ] || [ "${TARGETPLATFORM}" = "linux/arm64" ]; then \
        echo "Building for ${TARGETPLATFORM}..." && \
        COMMIT="$(git rev-parse --short HEAD)" && \
        DATE="$(date +%Y-%m-%dT%T%z)" && \
        go build -trimpath -ldflags="-s -w" -o /tmp/bin ./cmd/dlv && \
        chmod 500 /tmp/bin && \
        exit 0; \
    fi && \
    wget -qO- https://raw.githubusercontent.com/qdm12/binpot/main/unavailable/build.sh | \
       sh -s -- "dlv ${VERSION}" "${TARGETPLATFORM}"
```

If you build it for `linux/amd64`, the final `/tmp/bin` program will be the actual `dlv` binary.

Otherwise, if you build it for `linux/arm/v7` for example, it will be this unavailable program showing: `dlv v1.6.1 is unavailable for linux/arm/v7`.
