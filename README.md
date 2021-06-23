# Binpot

The statically compiled, cross architecture, Docker based, binaries pot.

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
- Each Docker image is based on scratch and only contain the binary program at the path `/bin`.
- Each Docker image and binary program is built for each of the possible CPU architecture supported by Docker, unless indicated otherwise.
- Each binary has permissions `500` (read and execute for the user owner)
- You can change the ownership in the `COPY` command using `--chown=1000` for example, without duplicating the binary in your Docker image layers.

**Need help?** ▶️ [Create a discussion](https://github.com/qdm12/binpot/discussions)

## Programs available

| Program | Version | Image tag | Architectures |
| --- | --- | --- | --- |
| `golangci-lint` | `v1.41.1` | `:golangci-lint` | all but `linux/riscv64` |
| `golangci-lint` | `v1.41.1` | `:golangci-lint-v1.41.1` | all but `linux/riscv64` |
| `helm` | `v3.6.1` | `:helm` | all but `linux/riscv64` |
| `helm` | `v3.6.1` | `:helm-v3.6.1` | all but `linux/riscv64` |
| `kubectl` | `v1.21.1` | `:kubectl` | all |
| `kubectl` | `v1.21.1` | `:kubectl-v1.21.1` | all |
| `kubectx` | `v0.9.3` | `:kubectx` | all |
| `kubectx` | `v0.9.3` | `:kubectx-v0.9.3` | all |
| `kubens` | `v0.9.3` | `:kubens` | all |
| `kubens` | `v0.9.3` | `:kubens-v0.9.3` | all |
| `stern` | `v1.18.0` | `:stern` | all |
| `stern` | `v1.18.0` | `:stern-v1.18.0` | all |

ℹ️ `all` architectures means: linux/amd64,linux/386,linux/arm64,linux/arm/v7,linux/arm/v6,linux/ppc64le,linux/s390x,linux/riscv64

**Want more!?** ▶️ [Create an issue!](https://github.com/qdm12/binpot/issues)

## How it works

1. For each program, a Dockerfile describes how to build it. The final binary is placed on a final [scratch](https://hub.docker.com/_/scratch) based Docker image. *Example:* `helm` has [`./dockerfiles/helm/Dockerfile`](dockerfiles/helm/Dockerfile)
2. For each program, a Github Action workflow is triggered when its Dockerfile or the workflow itself is changed. This workflow takes care of:
    1. Cross build the program for many CPU architectures
    2. Pushing the images containing the program to Docker Hub
