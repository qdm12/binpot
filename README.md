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

**Want more!?** ▶️ [Create an issue!](https://github.com/qdm12/binpot/issues)
