# First Rust based program of this repository: all the others are Go programs.
# Cross compiling Rust needs a C cross compiler too, because rumdl links
# jemalloc. cargo-zigbuild drives the build with zig (the Alpine package) as C
# cross compiler and linker, and every architecture below gives a fully static
# musl binary, like the Go based programs do.
#
# See the notes on Rust programs in the README for the file descriptor limit
# needed by zig and the Arm 32 bits cc wrapper, and the note on rumdl there for
# the linux/s390x substitute program.
#
# Alpine 3.24 rather than 3.23 like the other programs: it carries zig 0.16, with
# which all the platforms below have been validated (3.23 carries zig 0.15).
ARG ALPINE_VERSION=3.24
ARG RUST_VERSION=1.96.0
ARG ZIGBUILD_VERSION=v0.23.4
ARG XCPU_VERSION=v0.9.0
ARG BUILDPLATFORM=linux/amd64

# Only used to build the substitute program of this repository on the platforms
# rumdl cannot be built statically on.
FROM --platform=${BUILDPLATFORM} ghcr.io/qdm12/xcputranslate:${XCPU_VERSION} AS xcputranslate

FROM --platform=${BUILDPLATFORM} alpine:${ALPINE_VERSION} AS zigbuild
ARG ZIGBUILD_VERSION
RUN apk add --no-cache xz && \
    wget -q -O /tmp/cargo-zigbuild.tar.xz \
    "https://github.com/rust-cross/cargo-zigbuild/releases/download/${ZIGBUILD_VERSION}/cargo-zigbuild-x86_64-unknown-linux-musl.tar.xz" && \
    tar -xJf /tmp/cargo-zigbuild.tar.xz -C /usr/local/bin --strip-components=1 \
    "cargo-zigbuild-x86_64-unknown-linux-musl/cargo-zigbuild" && \
    cargo-zigbuild --version

FROM --platform=${BUILDPLATFORM} rust:${RUST_VERSION}-alpine${ALPINE_VERSION} AS build
# zig: the C cross compiler and linker; make: drives the jemalloc build;
# git: fetches the sources.
RUN apk add --no-cache git zig make && \
    git config --global advice.detachedHead false
COPY --from=zigbuild /usr/local/bin/cargo-zigbuild /usr/local/bin/cargo-zigbuild
COPY --from=xcputranslate /xcputranslate /usr/local/bin/xcputranslate
# Only used for the 32 bits Arm targets, see its header comment.
COPY zig-cc-wrapper.sh /usr/local/bin/zig-cc-wrapper
RUN chmod 500 /usr/local/bin/zig-cc-wrapper
WORKDIR /tmp/build

ARG VERSION=v0.2.79
RUN git clone --depth 1 --branch ${VERSION} https://github.com/rvben/rumdl.git .

ARG TARGETPLATFORM
RUN set -eux && \
    # Zig compiles musl from source for every target, which needs more than the
    # 1024 file descriptors (soft and hard limit) that `podman build` gives, so
    # it cannot be raised from here. Docker and buildx give a larger hard limit
    # and let the soft limit be raised, as done below. With podman, build with
    # --ulimit nofile=65536:65536 or the link step fails with
    # `error: sub-compilation of musl libc.a failed`.
    ulimit -Sn 8192 || echo "warning: file descriptor limit stays at $(ulimit -Sn)" && \
    # Values are the jemalloc page size log2, which has to be greater than or
    # equal to the page size of the systems running the binary. Upstream uses 12
    # for amd64 and 16 for arm64, the latter supporting the 16 KiB page systems
    # of Asahi Linux and Raspberry Pi 5.
    case "${TARGETPLATFORM}" in \
    linux/amd64) \
        RUST_TARGET=x86_64-unknown-linux-musl \
        JEMALLOC_LG_PAGE=12 ;; \
    linux/386) \
        RUST_TARGET=i686-unknown-linux-musl \
        JEMALLOC_LG_PAGE=12 ;; \
    linux/arm64) \
        RUST_TARGET=aarch64-unknown-linux-musl \
        JEMALLOC_LG_PAGE=16 ;; \
    linux/arm/v7) \
        RUST_TARGET=armv7-unknown-linux-musleabihf \
        ZIG_CC_TARGET=arm-linux-musleabihf \
        ZIG_CC_MCPU=generic+v7a \
        JEMALLOC_LG_PAGE=12 \
        USE_CC_WRAPPER=1 ;; \
    linux/arm/v6) \
        RUST_TARGET=arm-unknown-linux-musleabi \
        ZIG_CC_TARGET=arm-linux-musleabi \
        ZIG_CC_MCPU=generic+v6 \
        JEMALLOC_LG_PAGE=12 \
        USE_CC_WRAPPER=1 ;; \
    linux/ppc64le) \
        RUST_TARGET=powerpc64le-unknown-linux-musl \
        JEMALLOC_LG_PAGE=16 ;; \
    linux/riscv64) \
        RUST_TARGET=riscv64gc-unknown-linux-musl \
        JEMALLOC_LG_PAGE=16 ;; \
    linux/s390x) \
        UNAVAILABLE=1 ;; \
    *) \
        echo "There is no rumdl build for ${TARGETPLATFORM}" && exit 1 ;; \
    esac && \
    # No static binary is possible on linux/s390x (see the note on rumdl in the
    # README), so the substitute program of this repository is built there
    # instead: it prints a message and exits with 1, so that cross building with
    # all the architectures keeps working.
    if [ -n "${UNAVAILABLE:-}" ]; then \
        apk add --no-cache go && \
        wget -qO- https://raw.githubusercontent.com/qdm12/binpot/main/unavailable/build.sh | \
        sh -s -- "rumdl ${VERSION}" "${TARGETPLATFORM}" && \
        exit 0; \
    fi && \
    export JEMALLOC_SYS_WITH_LG_PAGE="${JEMALLOC_LG_PAGE}" && \
    # On the 32 bits Arm targets the flags rustc hands to the C compiler are not
    # understood by zig, so jemalloc (and any other C dependency) is built
    # through the wrapper instead of the cargo-zigbuild one. CC_<target> wins
    # over CC for cargo, both are set to the wrapper. ZIG_CC_TARGET and
    # ZIG_CC_MCPU are the zig equivalents of the Rust target and of the
    # -march/-mfpu/-mfloat-abi flags rustc uses, read by the wrapper.
    if [ -n "${USE_CC_WRAPPER:-}" ]; then \
        rust_target_key="CC_$(echo "${RUST_TARGET}" | tr '-' '_')" && \
        export "${rust_target_key}"=/usr/local/bin/zig-cc-wrapper && \
        export CC=/usr/local/bin/zig-cc-wrapper \
            ZIG_CC_WRAPPER_TARGET="${ZIG_CC_TARGET}" \
            ZIG_CC_WRAPPER_MCPU="${ZIG_CC_MCPU}"; \
    fi && \
    rustup target add "${RUST_TARGET}" && \
    cargo zigbuild --locked --release --bin rumdl --target "${RUST_TARGET}" && \
    mv "target/${RUST_TARGET}/release/rumdl" /tmp/bin && \
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
    org.opencontainers.image.url="https://github.com/qdm12/binpot/blob/main/dockerfiles/rumdl" \
    org.opencontainers.image.documentation="https://github.com/qdm12/binpot" \
    org.opencontainers.image.source="https://github.com/qdm12/binpot.git" \
    org.opencontainers.image.title="rumdl binary" \
    org.opencontainers.image.description="rumdl binary"
COPY --from=build /tmp/bin /bin
