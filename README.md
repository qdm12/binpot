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

[Thinking of copying the binary for your host?](#Copy-the-binary-on-your-host)

## Programs available

💁 Only the last bug fix releases for each program feature release is kept in the following table.

| Program | Version | Image tags | Architectures |
| --- | --- | --- | --- |
| [`bit`](https://github.com/chriswalz/bit) | [`v1.1.2`](https://github.com/chriswalz/bit/releases/tag/v1.1.2) | [`:bit`](https://hub.docker.com/r/qmcgaw/binpot/tags?name=bit) | all |
| [`bit`](https://github.com/chriswalz/bit) | [`v1.1.2`](https://github.com/chriswalz/bit/releases/tag/v1.1.2) | [`:bit-v1.1.2`](https://hub.docker.com/r/qmcgaw/binpot/tags?name=bit-v1.1.2) | all |
| [`buildx`](https://github.com/docker/buildx) | [`v0.8.1`](https://github.com/docker/buildx/releases/tag/v0.8.1) | [`:buildx`](https://hub.docker.com/r/qmcgaw/binpot/tags?name=buildx) | all |
| [`buildx`](https://github.com/docker/buildx) | [`v0.8.1`](https://github.com/docker/buildx/releases/tag/v0.8.1) | [`:buildx-v0.8.1`](https://hub.docker.com/r/qmcgaw/binpot/tags?name=buildx-v0.8.1) | all |
| [`buildx`](https://github.com/docker/buildx) | [`v0.7.1`](https://github.com/docker/buildx/releases/tag/v0.7.1) | [`:buildx-v0.7.1`](https://hub.docker.com/r/qmcgaw/binpot/tags?name=buildx-v0.7.1) | all |
| [`buildx`](https://github.com/docker/buildx) | [`v0.6.3`](https://github.com/docker/buildx/releases/tag/v0.6.3) | [`:buildx-v0.6.3`](https://hub.docker.com/r/qmcgaw/binpot/tags?name=buildx-v0.6.3) | all |
| [`buildx`](https://github.com/docker/buildx) | [`v0.5.1`](https://github.com/docker/buildx/releases/tag/v0.5.1) | [`:buildx-v0.5.1`](https://hub.docker.com/r/qmcgaw/binpot/tags?name=buildx-v0.5.1) | all |
| [`compose`](https://github.com/docker/compose) | [`v2.3.4`](https://github.com/docker/compose/releases/tag/v2.3.4) | [`:compose`](https://hub.docker.com/r/qmcgaw/binpot/tags?name=compose) | all |
| [`compose`](https://github.com/docker/compose) | [`v2.3.4`](https://github.com/docker/compose/releases/tag/v2.3.4) | [`:compose-v2.3.4`](https://hub.docker.com/r/qmcgaw/binpot/tags?name=compose-v2.3.4) | all |
| [`compose`](https://github.com/docker/compose) | [`v2.2.3`](https://github.com/docker/compose/releases/tag/v2.2.3) | [`:compose-v2.2.3`](https://hub.docker.com/r/qmcgaw/binpot/tags?name=compose-v2.2.3) | all |
| [`compose`](https://github.com/docker/compose) | [`v2.1.0`](https://github.com/docker/compose/releases/tag/v2.1.0) | [`:compose-v2.1.0`](https://hub.docker.com/r/qmcgaw/binpot/tags?name=compose-v2.1.0) | all |
| [`compose`](https://github.com/docker/compose) | [`v2.0.1`](https://github.com/docker/compose/releases/tag/v2.0.1) | [`:compose-v2.0.1`](https://hub.docker.com/r/qmcgaw/binpot/tags?name=compose-v2.0.1) | all |
| [`dlv`](https://github.com/go-delve/delve) | [`v1.8.2`](https://github.com/go-delve/delve/releases/tag/v1.8.2) | [`:dlv`](https://hub.docker.com/r/qmcgaw/binpot/tags?name=dlv) | `linux/amd64` and `linux/arm64` |
| [`dlv`](https://github.com/go-delve/delve) | [`v1.8.2`](https://github.com/go-delve/delve/releases/tag/v1.8.2) | [`:dlv-v1.8.2`](https://hub.docker.com/r/qmcgaw/binpot/tags?name=dlv-v1.8.2) | `linux/amd64` and `linux/arm64` |
| [`dlv`](https://github.com/go-delve/delve) | [`v1.7.3`](https://github.com/go-delve/delve/releases/tag/v1.7.3) | [`:dlv-v1.7.3`](https://hub.docker.com/r/qmcgaw/binpot/tags?name=dlv-v1.7.3) | `linux/amd64` and `linux/arm64` |
| [`dlv`](https://github.com/go-delve/delve) | [`v1.6.1`](https://github.com/go-delve/delve/releases/tag/v1.6.1) | [`:dlv-v1.6.1`](https://hub.docker.com/r/qmcgaw/binpot/tags?name=dlv-v1.6.1) | `linux/amd64` and `linux/arm64` |
| [`docker`](https://github.com/docker/cli) | [`v20.10.12`](https://github.com/docker/cli/releases/tag/v20.10.12) | [`:docker`](https://hub.docker.com/r/qmcgaw/binpot/tags?name=docker) | all |
| [`docker`](https://github.com/docker/cli) | [`v20.10.12`](https://github.com/docker/cli/releases/tag/v20.10.12) | [`:docker-v20.10.12`](https://hub.docker.com/r/qmcgaw/binpot/tags?name=docker-v20.10.12) | all |
| [`gh`](https://github.com/cli/cli) | [`v2.7.0`](https://github.com/cli/cli/releases/tag/v2.7.0) | [`:gh`](https://hub.docker.com/r/qmcgaw/binpot/tags?name=gh) | all |
| [`gh`](https://github.com/cli/cli) | [`v2.7.0`](https://github.com/cli/cli/releases/tag/v2.7.0) | [`:gh-v2.7.0`](https://hub.docker.com/r/qmcgaw/binpot/tags?name=gh-v2.7.0) | all |
| [`gh`](https://github.com/cli/cli) | [`v2.6.0`](https://github.com/cli/cli/releases/tag/v2.6.0) | [`:gh-v2.6.0`](https://hub.docker.com/r/qmcgaw/binpot/tags?name=gh-v2.6.0) | all |
| [`gh`](https://github.com/cli/cli) | [`v2.5.2`](https://github.com/cli/cli/releases/tag/v2.5.2) | [`:gh-v2.5.2`](https://hub.docker.com/r/qmcgaw/binpot/tags?name=gh-v2.5.2) | all |
| [`gh`](https://github.com/cli/cli) | [`v2.4.0`](https://github.com/cli/cli/releases/tag/v2.4.0) | [`:gh-v2.4.0`](https://hub.docker.com/r/qmcgaw/binpot/tags?name=gh-v2.4.0) | all |
| [`gh`](https://github.com/cli/cli) | [`v2.3.0`](https://github.com/cli/cli/releases/tag/v2.3.0) | [`:gh-v2.3.0`](https://hub.docker.com/r/qmcgaw/binpot/tags?name=gh-v2.3.0) | all |
| [`gh`](https://github.com/cli/cli) | [`v2.2.0`](https://github.com/cli/cli/releases/tag/v2.2.0) | [`:gh-v2.2.0`](https://hub.docker.com/r/qmcgaw/binpot/tags?name=gh-v2.2.0) | all |
| [`gh`](https://github.com/cli/cli) | [`v2.1.0`](https://github.com/cli/cli/releases/tag/v2.1.0) | [`:gh-v2.1.0`](https://hub.docker.com/r/qmcgaw/binpot/tags?name=gh-v2.1.0) | all |
| [`gh`](https://github.com/cli/cli) | [`v2.0.0`](https://github.com/cli/cli/releases/tag/v2.0.0) | [`:gh-v2.0.0`](https://hub.docker.com/r/qmcgaw/binpot/tags?name=gh-v2.0.0) | all |
| [`gh`](https://github.com/cli/cli) | [`v1.14.0`](https://github.com/cli/cli/releases/tag/v1.14.0) | [`:gh-v1.14.0`](https://hub.docker.com/r/qmcgaw/binpot/tags?name=gh-v1.14.0) | all |
| [`gh`](https://github.com/cli/cli) | [`v1.13.1`](https://github.com/cli/cli/releases/tag/v1.13.1) | [`:gh-v1.13.1`](https://hub.docker.com/r/qmcgaw/binpot/tags?name=gh-v1.13.1) | all |
| [`gh`](https://github.com/cli/cli) | [`v1.12.1`](https://github.com/cli/cli/releases/tag/v1.12.1) | [`:gh-v1.12.1`](https://hub.docker.com/r/qmcgaw/binpot/tags?name=gh-v1.12.1) | all |
| [`gh`](https://github.com/cli/cli) | [`v1.11.0`](https://github.com/cli/cli/releases/tag/v1.11.0) | [`:gh-v1.11.0`](https://hub.docker.com/r/qmcgaw/binpot/tags?name=gh-v1.11.0) | all |
| [`golangci-lint`](https://github.com/golangci/golangci-lint) | [`v1.45.2`](https://github.com/golangci/golangci-lint/releases/tag/v1.45.2) | [`:golangci-lint`](https://hub.docker.com/r/qmcgaw/binpot/tags?name=golangci-lint) | all |
| [`golangci-lint`](https://github.com/golangci/golangci-lint) | [`v1.45.2`](https://github.com/golangci/golangci-lint/releases/tag/v1.45.2) | [`:golangci-lint-v1.45.2`](https://hub.docker.com/r/qmcgaw/binpot/tags?name=golangci-lint-v1.45.2) | all |
| [`golangci-lint`](https://github.com/golangci/golangci-lint) | [`v1.44.2`](https://github.com/golangci/golangci-lint/releases/tag/v1.44.2) | [`:golangci-lint-v1.44.2`](https://hub.docker.com/r/qmcgaw/binpot/tags?name=golangci-lint-v1.44.2) | all |
| [`golangci-lint`](https://github.com/golangci/golangci-lint) | [`v1.43.0`](https://github.com/golangci/golangci-lint/releases/tag/v1.43.0) | [`:golangci-lint-v1.43.0`](https://hub.docker.com/r/qmcgaw/binpot/tags?name=golangci-lint-v1.43.0) | all |
| [`golangci-lint`](https://github.com/golangci/golangci-lint) | [`v1.42.1`](https://github.com/golangci/golangci-lint/releases/tag/v1.42.1) | [`:golangci-lint-v1.42.1`](https://hub.docker.com/r/qmcgaw/binpot/tags?name=golangci-lint-v1.42.1) | all |
| [`golangci-lint`](https://github.com/golangci/golangci-lint) | [`v1.41.1`](https://github.com/golangci/golangci-lint/releases/tag/v1.41.1) | [`:golangci-lint-v1.41.1`](https://hub.docker.com/r/qmcgaw/binpot/tags?name=golangci-lint-v1.41.1) | all |
| [`golangci-lint`](https://github.com/golangci/golangci-lint) | [`v1.40.1`](https://github.com/golangci/golangci-lint/releases/tag/v1.40.1) | [`:golangci-lint-v1.40.1`](https://hub.docker.com/r/qmcgaw/binpot/tags?name=golangci-lint-v1.40.1) | all |
| [`gomock`](https://github.com/golang/mock) | [`v1.6.0`](https://github.com/golang/mock/releases/tag/v1.6.0) | [`:gomock`](https://hub.docker.com/r/qmcgaw/binpot/tags?name=gomock) | all |
| [`gomock`](https://github.com/golang/mock) | [`v1.6.0`](https://github.com/golang/mock/releases/tag/v1.6.0) | [`:gomock-v1.6.0`](https://hub.docker.com/r/qmcgaw/binpot/tags?name=gomock-v1.6.0) | all |
| [`gomodifytags`](https://github.com/fatih/gomodifytags) | [`v1.16.0`](https://github.com/fatih/gomodifytags/releases/tag/v1.16.0) | [`:gomodifytags`](https://hub.docker.com/r/qmcgaw/binpot/tags?name=gomodifytags) | all |
| [`gomodifytags`](https://github.com/fatih/gomodifytags) | [`v1.16.0`](https://github.com/fatih/gomodifytags/releases/tag/v1.16.0) | [`:gomodifytags-v1.16.0`](https://hub.docker.com/r/qmcgaw/binpot/tags?name=gomodifytags-v1.15.0) | all |
| [`gomodifytags`](https://github.com/fatih/gomodifytags) | [`v1.15.0`](https://github.com/fatih/gomodifytags/releases/tag/v1.15.0) | [`:gomodifytags-v1.15.0`](https://hub.docker.com/r/qmcgaw/binpot/tags?name=gomodifytags-v1.15.0) | all |
| [`gomodifytags`](https://github.com/fatih/gomodifytags) | [`v1.14.0`](https://github.com/fatih/gomodifytags/releases/tag/v1.14.0) | [`:gomodifytags-v1.14.0`](https://hub.docker.com/r/qmcgaw/binpot/tags?name=gomodifytags-v1.14.0) | all |
| [`gomodifytags`](https://github.com/fatih/gomodifytags) | [`v1.13.0`](https://github.com/fatih/gomodifytags/releases/tag/v1.13.0) | [`:gomodifytags-v1.13.0`](https://hub.docker.com/r/qmcgaw/binpot/tags?name=gomodifytags-v1.13.0) | all |
| [`go-outline`](https://github.com/ramya-rao-a/go-outline) | [`9736a4b`](https://github.com/ramya-rao-a/go-outline/commit/9736a4bde949f321d201e5eaa5ae2bcde011bf00) | [`:go-outline`](https://hub.docker.com/r/qmcgaw/binpot/tags?name=go-outline) | all |
| [`go-outline`](https://github.com/ramya-rao-a/go-outline) | [`9736a4b`](https://github.com/ramya-rao-a/go-outline/commit/9736a4bde949f321d201e5eaa5ae2bcde011bf00) | [`:go-outline-9736a4b`](https://hub.docker.com/r/qmcgaw/binpot/tags?name=go-outline-9736a4b) | all |
| [`gopkgs`](https://github.com/uudashr/gopkgs) | [`v2.1.2`](https://github.com/uudashr/gopkgs/releases/tag/v2.1.2) | [`:gopkgs`](https://hub.docker.com/r/qmcgaw/binpot/tags?name=gopkgs) | all |
| [`gopkgs`](https://github.com/uudashr/gopkgs) | [`v2.1.2`](https://github.com/uudashr/gopkgs/releases/tag/v2.1.2) | [`:gopkgs-v2.1.2`](https://hub.docker.com/r/qmcgaw/binpot/tags?name=gopkgs-v2.1.2) | all |
| [`goplay`](https://github.com/haya14busa/goplay) | [`v1.0.0`](https://github.com/haya14busa/goplay/releases/tag/v1.0.0) | [`:goplay`](https://hub.docker.com/r/qmcgaw/binpot/tags?name=goplay) | all |
| [`goplay`](https://github.com/haya14busa/goplay) | [`v1.0.0`](https://github.com/haya14busa/goplay/releases/tag/v1.0.0) | [`:goplay-v1.0.0`](https://hub.docker.com/r/qmcgaw/binpot/tags?name=goplay-v1.0.0) | all |
| [`gopls`](https://github.com/golang/tools/tree/master/gopls) | [`v0.8.2`](https://github.com/golang/tools/releases/tag/gopls/v0.8.2) | [`:gopls`](https://hub.docker.com/r/qmcgaw/binpot/tags?name=gopls) | all |
| [`gopls`](https://github.com/golang/tools/tree/master/gopls) | [`v0.8.2`](https://github.com/golang/tools/releases/tag/gopls/v0.8.2) | [`:gopls-v0.8.2`](https://hub.docker.com/r/qmcgaw/binpot/tags?name=gopls-v0.8.2) | all |
| [`gopls`](https://github.com/golang/tools/tree/master/gopls) | [`v0.7.5`](https://github.com/golang/tools/releases/tag/gopls/v0.7.5) | [`:gopls-v0.7.5`](https://hub.docker.com/r/qmcgaw/binpot/tags?name=gopls-v0.7.5) | all |
| [`gotests`](https://github.com/cweill/gotests) | [`v1.6.0`](https://github.com/cweill/gotests/releases/tag/v1.6.0) | [`:gotests`](https://hub.docker.com/r/qmcgaw/binpot/tags?name=gotests) | all |
| [`gotests`](https://github.com/cweill/gotests) | [`v1.6.0`](https://github.com/cweill/gotests/releases/tag/v1.6.0) | [`:gotests-v1.6.0`](https://hub.docker.com/r/qmcgaw/binpot/tags?name=gotests-v1.6.0) | all |
| [`helm`](https://github.com/helm/helm) | [`v3.8.1`](https://github.com/helm/helm/releases/tag/v3.8.1) | [`:helm`](https://hub.docker.com/r/qmcgaw/binpot/tags?name=helm) | all |
| [`helm`](https://github.com/helm/helm) | [`v3.8.1`](https://github.com/helm/helm/releases/tag/v3.8.1) | [`:helm-v3.8.1`](https://hub.docker.com/r/qmcgaw/binpot/tags?name=helm-v3.8.1) | all |
| [`helm`](https://github.com/helm/helm) | [`v3.7.2`](https://github.com/helm/helm/releases/tag/v3.7.2) | [`:helm-v3.7.2`](https://hub.docker.com/r/qmcgaw/binpot/tags?name=helm-v3.7.2) | all |
| [`helm`](https://github.com/helm/helm) | [`v3.6.3`](https://github.com/helm/helm/releases/tag/v3.6.3) | [`:helm-v3.6.3`](https://hub.docker.com/r/qmcgaw/binpot/tags?name=helm-v3.6.3) | all |
| [`impl`](https://github.com/josharian/impl) | [`v1.1.0`](https://github.com/josharian/impl/releases/tag/v1.1.0) | [`:impl`](https://hub.docker.com/r/qmcgaw/binpot/tags?name=impl) | all |
| [`impl`](https://github.com/josharian/impl) | [`v1.1.0`](https://github.com/josharian/impl/releases/tag/v1.1.0) | [`:impl-v1.1.0`](https://hub.docker.com/r/qmcgaw/binpot/tags?name=impl-v1.1.0) | all |
| [`kubectl`](https://github.com/kubernetes/kubernetes) | [`v1.23.4`](https://github.com/kubernetes/kubernetes/releases/tag/v1.23.4) | [`:kubectl`](https://hub.docker.com/r/qmcgaw/binpot/tags?name=kubectl) | all |
| [`kubectl`](https://github.com/kubernetes/kubernetes) | [`v1.23.4`](https://github.com/kubernetes/kubernetes/releases/tag/v1.23.4) | [`:kubectl-v1.23.4`](https://hub.docker.com/r/qmcgaw/binpot/tags?name=kubectl-v1.23.4) | all |
| [`kubectl`](https://github.com/kubernetes/kubernetes) | [`v1.22.4`](https://github.com/kubernetes/kubernetes/releases/tag/v1.22.4) | [`:kubectl-v1.22.4`](https://hub.docker.com/r/qmcgaw/binpot/tags?name=kubectl-v1.22.4) | all |
| [`kubectl`](https://github.com/kubernetes/kubernetes) | [`v1.21.3`](https://github.com/kubernetes/kubernetes/releases/tag/v1.21.3) | [`:kubectl-v1.21.3`](https://hub.docker.com/r/qmcgaw/binpot/tags?name=kubectl-v1.21.3) | all |
| [`kubectx`](https://github.com/ahmetb/kubectx) | [`v0.9.4`](https://github.com/ahmetb/kubectx/releases/tag/v0.9.4) | [`:kubectx`](https://hub.docker.com/r/qmcgaw/binpot/tags?name=kubectx) | all |
| [`kubectx`](https://github.com/ahmetb/kubectx) | [`v0.9.4`](https://github.com/ahmetb/kubectx/releases/tag/v0.9.4) | [`:kubectx-v0.9.4`](https://hub.docker.com/r/qmcgaw/binpot/tags?name=kubectx-v0.9.4) | all |
| [`kubens`](https://github.com/ahmetb/kubectx) | [`v0.9.4`](https://github.com/ahmetb/kubectx/releases/tag/v0.9.4) | [`:kubens`](https://hub.docker.com/r/qmcgaw/binpot/tags?name=kubens) | all |
| [`kubens`](https://github.com/ahmetb/kubectx) | [`v0.9.4`](https://github.com/ahmetb/kubectx/releases/tag/v0.9.4) | [`:kubens-v0.9.4`](https://hub.docker.com/r/qmcgaw/binpot/tags?name=kubens-v0.9.4) | all |
| [`logo-ls`](https://github.com/Yash-Handa/logo-ls) | [`v1.3.7`](https://github.com/Yash-Handa/logo-ls/releases/tag/v1.3.7) | [`:logo-ls`](https://hub.docker.com/r/qmcgaw/binpot/tags?name=logo-ls) | all |
| [`logo-ls`](https://github.com/Yash-Handa/logo-ls) | [`v1.3.7`](https://github.com/Yash-Handa/logo-ls/releases/tag/v1.3.7) | [`:logo-ls-v1.3.7`](https://hub.docker.com/r/qmcgaw/binpot/tags?name=logo-ls-v1.3.7) | all |
| [`mockery`](https://github.com/vektra/mockery) | [`v2.10.1`](https://github.com/vektra/mockery/releases/tag/v2.10.1) | [`:mockery`](https://hub.docker.com/r/qmcgaw/binpot/tags?name=mockery) | all |
| [`mockery`](https://github.com/vektra/mockery) | [`v2.10.1`](https://github.com/vektra/mockery/releases/tag/v2.10.1) | [`:mockery-v2.10.1`](https://hub.docker.com/r/qmcgaw/binpot/tags?name=mockery-v2.10.1) | all |
| [`mockery`](https://github.com/vektra/mockery) | [`v2.9.5`](https://github.com/vektra/mockery/releases/tag/v2.9.5) | [`:mockery-v2.9.5`](https://hub.docker.com/r/qmcgaw/binpot/tags?name=mockery-v2.9.5) | all |
| [`mockgen`](https://github.com/golang/mock) | [`v1.6.0`](https://github.com/golang/mock/releases/tag/v1.6.0) | [`:mockgen`](https://hub.docker.com/r/qmcgaw/binpot/tags?name=mockgen) | all |
| [`mockgen`](https://github.com/golang/mock) | [`v1.6.0`](https://github.com/golang/mock/releases/tag/v1.6.0) | [`:mockgen-v1.6.0`](https://hub.docker.com/r/qmcgaw/binpot/tags?name=mockgen-v1.6.0) | all |
| [`stern`](https://github.com/stern/stern) | [`v1.21.0`](https://github.com/stern/stern/releases/tag/v1.21.0) | [`:stern`](https://hub.docker.com/r/qmcgaw/binpot/tags?name=stern) | all |
| [`stern`](https://github.com/stern/stern) | [`v1.21.0`](https://github.com/stern/stern/releases/tag/v1.21.0) | [`:stern-v1.21.0`](https://hub.docker.com/r/qmcgaw/binpot/tags?name=stern-v1.21.0) | all |
| [`stern`](https://github.com/stern/stern) | [`v1.20.1`](https://github.com/stern/stern/releases/tag/v1.20.1) | [`:stern-v1.20.1`](https://hub.docker.com/r/qmcgaw/binpot/tags?name=stern-v1.20.1) | all |
| [`stern`](https://github.com/stern/stern) | [`v1.19.0`](https://github.com/stern/stern/releases/tag/v1.19.0) | [`:stern-v1.19.0`](https://hub.docker.com/r/qmcgaw/binpot/tags?name=stern-v1.19.0) | all |
| [`stern`](https://github.com/stern/stern) | [`v1.18.0`](https://github.com/stern/stern/releases/tag/v1.18.0) | [`:stern-v1.18.0`](https://hub.docker.com/r/qmcgaw/binpot/tags?name=stern-v1.18.0) | all |
| [`supervisord`](https://github.com/ochinchina/supervisord) | [`v0.7.3`](https://github.com/ochinchina/supervisord/releases/tag/v0.7.3) | [`:supervisord`](https://hub.docker.com/r/qmcgaw/binpot/tags?name=supervisord) | all |
| [`supervisord`](https://github.com/ochinchina/supervisord) | [`v0.7.3`](https://github.com/ochinchina/supervisord/releases/tag/v0.7.3) | [`:supervisord-v0.7.3`](https://hub.docker.com/r/qmcgaw/binpot/tags?name=supervisord-v0.7.3) | all |

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
