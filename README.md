# Binpot

The statically compiled, cross architecture, Docker based, binaries pot.

<img height="250" src="https://raw.githubusercontent.com/qdm12/binpot/main/binpot.svg?sanitize=true">

## Usage

The usage is focused to build other Docker images.

For example:

```Dockerfile
FROM alpine:3.13
COPY --from=qmcgaw/binpot:helm /bin /usr/local/bin/helm
```

- The Docker image tags for `qmcgaw/binpot` follow the following formatting:
  - `:name` for the latest stable version of the program `name`
  - `:name-v0.0.0` for the semver version of the program `name`
- Each Docker image is based on [scratch](https://hub.docker.com/_/scratch) and only contain the binary program at the path `/bin`.
- Each Docker image and binary program is built for each of the possible CPU architecture supported by Docker, unless indicated otherwise.
- Each binary has permissions `500` (read and execute for the user owner)
- You can change the ownership in the `COPY` command using `--chown=1000` for example, without duplicating the binary in your Docker image layers.
- Each Docker image has an entrypoint `[ "/bin" ]` so you can run it as well

**Need help?** ▶️ [Create a discussion](https://github.com/qdm12/binpot/discussions)

[Thinking of copying the binary for your host?](#Copy-the-binary-on-your-host)

## Programs available

| Program | Version | Image tag | Architectures |
| --- | --- | --- | --- |
| `bit` | `v1.1.1` | `:bit` | all |
| `bit` | `v1.1.1` | `:bit-v1.1.1` | all |
| `buildx` | `31309b9` | `:buildx` | all |
| `buildx` | `31309b9` | `:buildx-31309b9` | all |
| `buildx` | `v0.5.1` | `:buildx-v0.5.1` | all |
| `compose` | `v2.0.0-beta.4` | `:compose` | all |
| `compose` | `v2.0.0-beta.4` | `:compose-v2.0.0-beta.4` | all |
| `dlv` | `v1.6.1` | `:dlv` | `linux/amd64` and `linux/arm64` |
| `dlv` | `v1.6.1` | `:dlv-v1.6.1` | `linux/amd64` and `linux/arm64` |
| `docker` | `v20.10.7` | `:docker` | all |
| `docker` | `v20.10.7` | `:docker-v20.10.7` | all |
| `gh` | `v1.11.0` | `:gh` | all |
| `gh` | `v1.11.0` | `:gh-v1.11.0` | all |
| `golangci-lint` | `v1.41.1` | `:golangci-lint` | all |
| `golangci-lint` | `v1.41.1` | `:golangci-lint-v1.41.1` | all |
| `gomock` | `v1.6.0` | `:gomock` | all |
| `gomock` | `v1.6.0` | `:gomock-v1.6.0` | all |
| `gomodifytags` | `v1.13.0` | `:gomodifytags` | all |
| `gomodifytags` | `v1.13.0` | `:gomodifytags-v1.13.0` | all |
| `go-outline` | `9736a4b` | `:go-outline` | all |
| `go-outline` | `9736a4b` | `:go-outline-9736a4b` | all |
| `gopkgs` | `v2.1.2` | `:gopkgs` | all |
| `gopkgs` | `v2.1.2` | `:gopkgs-v2.1.2` | all |
| `goplay` | `v1.0.0` | `:goplay` | all |
| `goplay` | `v1.0.0` | `:goplay-v1.0.0` | all |
| `gopls` | `v0.7.0` | `:gopls` | all |
| `gopls` | `v0.7.0` | `:gopls-v0.7.0` | all |
| `gotests` | `v1.6.0` | `:gotests` | all |
| `gotests` | `v1.6.0` | `:gotests-v1.6.0` | all |
| `helm` | `v3.6.1` | `:helm` | all |
| `helm` | `v3.6.1` | `:helm-v3.6.1` | all |
| `impl` | `v1.1.0` | `:impl` | all |
| `impl` | `v1.1.0` | `:impl-v1.1.0` | all |
| `kubectl` | `v1.21.1` | `:kubectl` | all |
| `kubectl` | `v1.21.1` | `:kubectl-v1.21.1` | all |
| `kubectx` | `v0.9.3` | `:kubectx` | all |
| `kubectx` | `v0.9.3` | `:kubectx-v0.9.3` | all |
| `kubens` | `v0.9.3` | `:kubens` | all |
| `kubens` | `v0.9.3` | `:kubens-v0.9.3` | all |
| `logo-ls` | `v1.6.1` | `:logo-ls` | all |
| `logo-ls` | `v1.6.1` | `:logo-ls-v1.3.7` | all |
| `mockery` | `v2.9.0` | `:mockery` | all |
| `mockery` | `v2.9.0` | `:mockery-v2.9.0` | all |
| `mockgen` | `v1.6.0` | `:mockgen` | all |
| `mockgen` | `v1.6.0` | `:mockgen-v1.6.0` | all |
| `stern` | `v1.18.0` | `:stern` | all |
| `stern` | `v1.18.0` | `:stern-v1.18.0` | all |

ℹ️ `all` architectures means: linux/amd64,linux/386,linux/arm64,linux/arm/v7,linux/arm/v6,linux/ppc64le,linux/s390x,linux/riscv64

**Want more!?** ▶️ [Create an issue!](https://github.com/qdm12/binpot/issues)

## How it works

1. For each program, a Dockerfile describes how to build it. The final binary is placed on a final [scratch](https://hub.docker.com/_/scratch) based Docker image. *Example:* `helm` has [`./dockerfiles/helm/Dockerfile`](dockerfiles/helm/Dockerfile)
2. For each program, a Github Action workflow is triggered when its Dockerfile or the workflow itself is changed. This workflow takes care of:
    1. Cross build the program for all CPU architectures
        - If one architecture is not supported such as for `dlv`, build the [unavailable](unavailable) program
    2. Pushing the images containing the program to Docker Hub

## Note on `dlv`

💁 Concerning `dlv`: all images are built for all architectures even if the program does not support all of them.
A substitute Go program printing `dlv v1.6.1 is unavailable on <platform name>` and exiting with exit code `1` is used for unsupported platforms.
This is like so so you can still cross build with all the architectures, especially if the program is an optional dependency.
This is often the case for VSCode development containers for instance. In this case, if you try to build for `arm/v7` and need `dlv` as an optional dependency, your `COPY --from=qmcgaw/binpot:dlv` will not fail.

## Copy the binary on your host

If you want to use the binary directly on your host, you can do it with Docker.
This has the advantage that it will automatically get the right binary for your host platform.

In the following we want to install `helm` on our host.

For example:

```sh
PROGRAM="helm" docker pull "qmcgaw/binpot:$PROGRAM" && \
  containerid="$(docker create qmcgaw/binpot:$PROGRAM)" && \
  docker cp "$containerid:/bin" "/usr/local/bin/$PROGRAM" && \
  docker rm "$containerid"
```

Test Helm works with `helm`

## TODOs

- Change version of Buildx once a release tag after v0.5.1 is made for Buildx to support RISCV64 officially: [Github issue](https://github.com/docker/buildx/issues/643)
- Change Dockerfile if block when Golangci-lint supports RISV64 officially: [Github issue](https://github.com/golangci/golangci-lint/issues/2079)
- Change Dockerfile if block when Helm supports RISV64 officially: [Github issue](https://github.com/helm/helm/issues/9858)
- Change version of go-outline once a release tag is made: [Github issue](https://github.com/ramya-rao-a/go-outline/issues/15)
