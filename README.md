# Binpot

The statically compiled, cross architecture, Docker based, binaries pot.

<img height="150" src="https://raw.githubusercontent.com/qdm12/binpot/main/binpot.svg">

[Blog post](https://qqq.ninja/blog/post/binpot/)

[![MIT](https://img.shields.io/github/license/qdm12/binpot)](https://github.com/qdm12/binpot/main/LICENSE)

## Usage

The usage is focused to build other Docker images.

For example:

```Dockerfile
FROM alpine:3.14
COPY --from=qmcgaw/binpot:helm /bin /usr/local/bin/helm
```

[🔍 Search for all image tags 💡](https://hub.docker.com/r/qmcgaw/binpot/tags)

- The Docker image tags for `qmcgaw/binpot` follow the following formatting:
  - `:name` for the latest stable version of the program `name`
  - `:name-v0.0.0` for the semver version of the program `name`
- Each Docker image is based on [scratch](https://hub.docker.com/_/scratch) and only contain the binary program at the path `/bin`.
- Each Docker image and binary program is built for each of the possible CPU architecture supported by Docker, unless indicated otherwise.
- Each binary has permissions `500` (read and execute for the user owner)
- You can change the ownership in the `COPY` command using `--chown=1000` for example, without duplicating the binary in your Docker image layers.
- Each Docker image has an entrypoint `[ "/bin" ]` so you can run it as well

**Need help?** ▶️ [Create a discussion](https://github.com/qdm12/binpot/discussions)

[Thinking of copying the binary for your host?](#copy-the-binary-on-your-host)

## Programs available

| Program | Last version | Image tags | Architectures |
| --- | --- | --- | --- |
| [`bit`](https://github.com/chriswalz/bit) | [`v1.1.2`](https://github.com/chriswalz/bit/releases/tag/v1.1.2) | [**Docker Hub**](https://hub.docker.com/r/qmcgaw/binpot/tags?name=bit) | all |
| [`buildx`](https://github.com/docker/buildx) | [`v0.17.1`](https://github.com/docker/buildx/releases/tag/v0.17.1) | [**Docker Hub**](https://hub.docker.com/r/qmcgaw/binpot/tags?name=buildx) | all |
| [`compose`](https://github.com/docker/compose) | [`v2.29.7`](https://github.com/docker/compose/releases/tag/v2.29.7) | [**Docker Hub**](https://hub.docker.com/r/qmcgaw/binpot/tags?name=compose) | all |
| [`dlv`](https://github.com/go-delve/delve) | [`v1.23.1`](https://github.com/go-delve/delve/releases/tag/v1.23.1) | [**Docker Hub**](https://hub.docker.com/r/qmcgaw/binpot/tags?name=dlv) | `linux/amd64` and `linux/arm64` |
| [`docker`](https://github.com/docker/cli) | [`v25.0.2`](https://github.com/docker/cli/releases/tag/v25.0.2) | [**Docker Hub**](https://hub.docker.com/r/qmcgaw/binpot/tags?name=docker) | all |
| [`gh`](https://github.com/cli/cli) | [`v2.58.0`](https://github.com/cli/cli/releases/tag/v2.58.0) | [**Docker Hub**](https://hub.docker.com/r/qmcgaw/binpot/tags?name=gh) | all |
| [`golangci-lint`](https://github.com/golangci/golangci-lint) | [`v1.61.0`](https://github.com/golangci/golangci-lint/releases/tag/v1.61.0) | [**Docker Hub**](https://hub.docker.com/r/qmcgaw/binpot/tags?name=golangci-lint) | all |
| [`gomock`](https://github.com/golang/mock) | [`v1.6.0`](https://github.com/golang/mock/releases/tag/v1.6.0) | [**Docker Hub**](https://hub.docker.com/r/qmcgaw/binpot/tags?name=gomock) | all |
| [`gomodifytags`](https://github.com/fatih/gomodifytags) | [`v1.16.0`](https://github.com/fatih/gomodifytags/releases/tag/v1.16.0) | [**Docker Hub**](https://hub.docker.com/r/qmcgaw/binpot/tags?name=gomodifytags) | all |
| [`go-outline`](https://github.com/ramya-rao-a/go-outline) | [`9736a4b`](https://github.com/ramya-rao-a/go-outline/commit/9736a4bde949f321d201e5eaa5ae2bcde011bf00) | [**Docker Hub**](https://hub.docker.com/r/qmcgaw/binpot/tags?name=go-outline) | all |
| [`gopkgs`](https://github.com/uudashr/gopkgs) | [`v2.1.2`](https://github.com/uudashr/gopkgs/releases/tag/v2.1.2) | [**Docker Hub**](https://hub.docker.com/r/qmcgaw/binpot/tags?name=gopkgs) | all |
| [`goplay`](https://github.com/haya14busa/goplay) | [`v1.0.0`](https://github.com/haya14busa/goplay/releases/tag/v1.0.0) | [**Docker Hub**](https://hub.docker.com/r/qmcgaw/binpot/tags?name=goplay) | all |
| [`gopls`](https://github.com/golang/tools/tree/master/gopls) | [`v0.14.2`](https://github.com/golang/tools/releases/tag/gopls/v0.14.2) | [**Docker Hub**](https://hub.docker.com/r/qmcgaw/binpot/tags?name=gopls) | all |
| [`gotests`](https://github.com/cweill/gotests) | [`v1.6.0`](https://github.com/cweill/gotests/releases/tag/v1.6.0) | [**Docker Hub**](https://hub.docker.com/r/qmcgaw/binpot/tags?name=gotests) | all |
| [`helm`](https://github.com/helm/helm) | [`v3.16.2`](https://github.com/helm/helm/releases/tag/v3.16.2) | [**Docker Hub**](https://hub.docker.com/r/qmcgaw/binpot/tags?name=helm) | all |
| [`impl`](https://github.com/josharian/impl) | [`v1.2.0`](https://github.com/josharian/impl/releases/tag/v1.2.0) | [**Docker Hub**](https://hub.docker.com/r/qmcgaw/binpot/tags?name=impl) | all |
| [`kubectl`](https://github.com/kubernetes/kubernetes) | [`v1.31.1`](https://github.com/kubernetes/kubernetes/releases/tag/v1.31.1) | [**Docker Hub**](https://hub.docker.com/r/qmcgaw/binpot/tags?name=kubectl) | all |
| [`kubectx`](https://github.com/ahmetb/kubectx) | [`v0.9.5`](https://github.com/ahmetb/kubectx/releases/tag/v0.9.5) | [**Docker Hub**](https://hub.docker.com/r/qmcgaw/binpot/tags?name=kubectx) | all |
| [`kubens`](https://github.com/ahmetb/kubectx) | [`v0.9.5`](https://github.com/ahmetb/kubectx/releases/tag/v0.9.5) | [**Docker Hub**](https://hub.docker.com/r/qmcgaw/binpot/tags?name=kubens) | all |
| [`logo-ls`](https://github.com/Yash-Handa/logo-ls) | [`v1.3.7`](https://github.com/Yash-Handa/logo-ls/releases/tag/v1.3.7) | [**Docker Hub**](https://hub.docker.com/r/qmcgaw/binpot/tags?name=logo-ls) | all |
| [`logo-ls`](https://github.com/Yash-Handa/logo-ls) | [`v1.3.7`](https://github.com/Yash-Handa/logo-ls/releases/tag/v1.3.7) | [**Docker Hub**](https://hub.docker.com/r/qmcgaw/binpot/tags?name=logo-ls-v1.3.7) | all |
| [`mockery`](https://github.com/vektra/mockery) | [`v2.46.2`](https://github.com/vektra/mockery/releases/tag/v2.46.2) | [**Docker Hub**](https://hub.docker.com/r/qmcgaw/binpot/tags?name=mockery) | all |
| [`mockgen`](https://github.com/golang/mock) | [`v1.6.0`](https://github.com/golang/mock/releases/tag/v1.6.0) | [**Docker Hub**](https://hub.docker.com/r/qmcgaw/binpot/tags?name=mockgen) | all |
| [`stern`](https://github.com/stern/stern) | [`v1.31.0`](https://github.com/stern/stern/releases/tag/v1.31.0) | [**Docker Hub**](https://hub.docker.com/r/qmcgaw/binpot/tags?name=stern) | all |
| [`supervisord`](https://github.com/ochinchina/supervisord) | [`v0.7.3`](https://github.com/ochinchina/supervisord/releases/tag/v0.7.3) | [**Docker Hub**](https://hub.docker.com/r/qmcgaw/binpot/tags?name=supervisord) | all |

ℹ️ `all` architectures means: linux/amd64, linux/386, linux/arm64, linux/arm/v7, linux/arm/v6, linux/ppc64le, linux/s390x, linux/riscv64

**Want more!?** ▶️ [Create an issue!](https://github.com/qdm12/binpot/issues)

## How it works

1. For each program, a Dockerfile describes how to build it. The final binary is placed on a final [scratch](https://hub.docker.com/_/scratch) based Docker image. *Example:* `helm` has [`./dockerfiles/helm/Dockerfile`](dockerfiles/helm/Dockerfile)
2. For each program, a Github Action workflow is triggered when its Dockerfile or the workflow itself is changed. This workflow takes care of:
    1. Cross build the program for all CPU architectures
        - If one architecture is not supported such as for `dlv`, build the [unavailable](unavailable) program
    2. Pushing the images containing the program to Docker Hub

## Note on `dlv`

💁 Concerning `dlv`: all images are built for all architectures even if the program does not support all of them.
A substitute Go program printing `dlv v1.7.0 is unavailable on <platform name>` and exiting with exit code `1` is used for unsupported platforms.
This is like so so you can still cross build with all the architectures, especially if the program is an optional dependency.
This is often the case for VSCode development containers for instance. In this case, if you try to build for `arm/v7` and need `dlv` as an optional dependency, your `COPY --from=qmcgaw/binpot:dlv` will not fail.

## Copy the binary on your host

If you want to use the binary directly on your host, you can do it with Docker.
This has the advantage that it will automatically get the right binary for your host platform.

In the following we want to install `helm` on our host.

For example:

```sh
export PROGRAM="helm" && \
  docker pull "qmcgaw/binpot:$PROGRAM" && \
  containerid="$(docker create qmcgaw/binpot:$PROGRAM)" && \
  docker cp "$containerid:/bin" "/usr/local/bin/$PROGRAM" && \
  docker rm "$containerid"
```

Test Helm works with `helm`

## TODOs

- Change version of go-outline once a release tag is made: [Github issue](https://github.com/ramya-rao-a/go-outline/issues/15)
