# AGENTS

Guidance for AI agent sessions working in this repository.

Binpot builds 27 single-binary developer tools as tiny Docker images: one `Dockerfile` and one
build workflow per program, pushed to Docker Hub and ghcr.io by
[`.github/workflows/reusable.yml`](.github/workflows/reusable.yml).

## The one rule to know

A program's version is written in **three** places, and they must agree:

1. `ARG VERSION=vX.Y.Z` in `dockerfiles/<program>/Dockerfile`
2. `version: vX.Y.Z` in `.github/workflows/<program>.yml`, passed to the reusable workflow
3. The program's row in the table in [`README.md`](README.md), which also embeds the version in
   two URLs

Never edit one of the three by hand. Use the checker below: it reads all three, refuses to touch
a program whose pins disagree, and rewrites them together.

## Upgrading programs

```sh
export GH_TOKEN=...                  # optional; raises the Github rate limit from 60 to 5000/h

go run ./cmd/versions                # check every program
go run ./cmd/versions stern helm     # check a few
go run ./cmd/versions -update        # bump every out-of-date program, leave it unstaged
go run ./cmd/versions -update -commit  # ...and commit one commit per program
```

Useful flags: `-dry-run` (what would change, changes nothing), `-quiet` (table only, no per-program
progress), `-json`, `-allow-drift` (override the disagreement guard, only when you are fixing the
drift yourself), `-dir` (run against another checkout), `-pages` (repositories listing more than
200 releases need a higher page count). Exit status is 1 if any program errored, so it is safe in
CI.

The tool needs no Docker and builds nothing. It only edits text.

## What the tool guarantees

- **Prefers published releases over tags**, falling back to tags only for repositories that
  publish no releases (`docker/cli`). A raw tag listing is not trustworthy on its own:
  `ochinchina/supervisord` has a stray `0.61` tag that outranks its real `v0.7.6` release.
- **Ignores pre-releases and drafts** — `v1.37.0-rc.1`, `-alpha.0`, `+insecure` are never selected.
- **Splices only the version bytes**, so CRLF endings, trailing blanks and indentation survive.
  `dockerfiles/gopls/Dockerfile` is CRLF and has a trailing space after the version.
- **Never pushes.** Committing is the last thing it does.

## Upgraded by hand, never by the tool

| program | why |
| --- | --- |
| `go-outline` | pinned to a commit (`9736a4b`); upstream publishes no releases |
| `gitstatus` | has no workflow file, so only two of the three pins exist |

## Quirks to know before editing

- [`.github/workflows/gomoditytags.yml`](.github/workflows/gomoditytags.yml) is a typo'd **file
  name** declaring `name: gomodifytags`. Leave the name alone; the tool finds it by declared name.
- `gomock` and `mockgen` share one repository (`uber-go/mock`) and their images together;
  `kubectx` and `kubens` share `ahmetb/kubectx`.
- `gopls` builds from `golang/tools` using the tag `gopls/vX.Y.Z`, not a bare version.
- `rumdl` links with upstream's `lto = "fat"` and `codegen-units = 1`, which peaks at around 4 GiB of memory per architecture, and buildx runs one build step per platform at once: the seven real platforms need around 24 GiB together, more than the 15.6 GiB of a GitHub runner, so the runner dies mid-link without uploading any log ("The hosted runner lost communication with the server", run 44845650721). `.github/workflows/rumdl.yml` caps buildkit with `[worker.oci] max-parallelism = 1` through the reusable workflow's `buildkitd_config` input; keep that cap unless the link gets lighter.
- `dockerfiles/kubectl/Dockerfile` carries pre-existing stale `gitMajor`/`gitMinor` ldflags. They
  are not the pinned version and are not yours to "fix" as part of an upgrade.
- The repository has no Go code of its own besides `cmd/versions`; `go.mod` is `go 1.27`, stdlib
  only, so there is no `go.sum` and no dependency work to do.
- `.pi/` is a local editor directory. Never add it to a commit.

## Adding a program

1. `dockerfiles/<program>/Dockerfile`, ending on a `FROM scratch` stage that copies the binary in
   with `COPY --from=build /tmp/bin /bin`.
2. `.github/workflows/<program>.yml` calling the reusable workflow with `name` and `version`.
3. A row in the `README.md` table, matching the existing row shape exactly.
4. One entry in the registry in [`internal/program/program.go`](internal/program/program.go) —
   `Name`, `Repo`, `TagPrefix` when upstream namespaces its tags, and `Track: TrackCommit` for a
   commit-pinned program. Without this line the checker silently ignores the new program.

## Commit conventions

- An upgrade is one commit per program, titled `` `programname` vX.Y.Z `` (backticks around the
  name, no other prefix). `-update -commit` already produces this.
- Tooling and infrastructure follow Conventional Commits as used by history:
  `chore(dev): ...`, `chore(ci): ...`, `fix: ...`.
- Do not push. Leave the commits for review.

## Verify before finishing

```sh
gofmt -l . && go vet ./... && go test ./...   # code
go run ./cmd/versions                          # every program reports "up to date"
git show --stat HEAD                           # an upgrade touches exactly 3 files
```
