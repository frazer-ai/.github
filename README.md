# .github

The `frazer-ai` organization's shared GitHub configuration: the public org
profile, default community files, the Renovate preset and the reusable
workflows every repository calls. Plan: task 1a in `docs`
(`plan/implementation-plan.md`) and `plan/repositories.md` §9.

| Path | What it is |
|---|---|
| `profile/README.md` | The public org profile |
| `SECURITY.md`, `.github/pull_request_template.md`, `.github/ISSUE_TEMPLATE/` | Defaults for every repository without its own |
| `renovate-config.json` | The shared Renovate preset |
| `.github/workflows/go-check.yml` | Go gate: gofmt, `go mod tidy -diff`, vet, golangci-lint, race tests (shuffled, uncached), govulncheck; with `tools`, the tests' tools at the caller's `mise.toml` versions |
| `.github/workflows/web-check.yml` | Web gate: `pnpm install --frozen-lockfile`, then `typecheck lint format:check test build` when defined |
| `.github/workflows/tofu-check.yml` | OpenTofu gate: fmt, init without a backend, validate, TFLint, Trivy config scan (digest-pinned, signature-verified image; its checks bundle from GHCR only, no version check or telemetry), `tofu test` |
| `.github/workflows/image-build.yml` | The one image builder: ko for a Go package, or BuildKit (pinned by digest) for a `dockerfile` inside its build context, multi-arch; on push, the image pushed by digest (crane), provenance and SBOM attestations with actions/attest, and the same predicates attested again on Sigstore's public-good instance (Rekor) by cosign, which frazer-install verifies offline; signing with the OpenBao transit key, and only then its `tags` (latest); with `vulnerability-gate`, Grype checks the image's CycloneDX SBOM before anything is attested or tagged, on PR builds too |
| `.github/workflows/scan.yml` | Security scanners: Betterleaks (a pull request's own commits; a pushed branch's history; every ref with `all-refs: true` for a scheduled run), OSV-Scanner (lockfiles), Trivy (vulnerabilities and misconfigurations) from its digest-pinned, signature-verified image, with its databases and checks bundle from GHCR only (not mirror.gcr.io) and no request to check.trivy.dev |
| `.github/workflows/names.yml` | Naming gates: no other cloud's region IDs in tracked files, Conventional Commit PR titles |
| `.github/workflows/license-gate.yml` | Licence gate: Syft SBOM, SPDX allow-list, exceptions file, generated NOTICE |
| `.github/workflows/self-check.yml` | This repository's CI: lints the workflows and runs each one against `testdata/`, including fixtures that must fail (`license-bad`, `vuln`, `image-vuln`; the last is also published, and must get no tag and no attestation), publishes the fixture images, ko's and the Dockerfile's, to `ghcr.io/frazer-ai/image-build-selftest` with their attestations, and builds a ko fixture on that private package as its base, which a token without `packages: read` must not read (not on a fork's pull request) |

## Using the workflows

Pin the reusable workflow by commit SHA, with the tag in a comment, and give
the calling job only the permissions it needs:

```yaml
jobs:
  go:
    uses: frazer-ai/.github/.github/workflows/go-check.yml@<sha> # v1
    with:
      working-directory: .
```

Each workflow documents its inputs at the top. The same workflow runs in the
pull-request check and in the release, so a release always ran what the PR ran.

### Tools for the tests

Tests that run a tool, such as `helm template` in a chart test, name it in
go-check's `tools` input. Each line gives a key of the caller's `mise.toml`
`[tools]` table and the Go package that builds the tool:

```yaml
jobs:
  go:
    uses: frazer-ai/.github/.github/workflows/go-check.yml@<sha>
    with:
      tools: |
        helm helm.sh/helm/v4/cmd/helm
```

The test job runs `go install <package>@v<version>` with the version
`mise.toml` pins, which must be exact. The binary goes into the job's own
`$RUNNER_TEMP/bin`, which comes first on PATH. A step then checks that each
tool resolves there and lists the tools in the job summary. On a self-hosted
runner, `~/go/bin` keeps what earlier jobs installed. A tool installed this
way never comes from there, and is never another version than the pinned one.
`mise-toml` names another manifest, relative to the repository root. The
`go-tools` self-check proves the input on `testdata/go-tools`.

### An image from a Dockerfile

`image-build` builds a Go package with ko (`importpath`). An image that is
not one Go package gives `dockerfile` instead, relative to
`working-directory`, which is its build context and must contain it. The
image gets the same provenance, SBOM and signatures, with this workflow as
the builder:

```yaml
jobs:
  resolver:
    permissions:
      contents: read
      packages: write
      id-token: write
      attestations: write
    uses: frazer-ai/.github/.github/workflows/image-build.yml@<sha> # v1
    with:
      runs-on: ubuntu-24.04
      working-directory: images/unbound
      dockerfile: Dockerfile
      repository: ghcr.io/frazer-ai/frazer-dns-resolver
      platforms: linux/amd64
      push: true
```

BuildKit runs from `moby/buildkit` pinned by digest (`buildkit-version`,
`buildkit-digest`), and sets every layer file's time and the image's
creation time to the commit's (`SOURCE_DATE_EPOCH`), so a Dockerfile whose
inputs are pinned gives one image per commit. RUN steps execute on the
runner's architecture and nothing installs emulation: build for the
runner's platform, or cross-compile with `FROM --platform=$BUILDPLATFORM`.

### A private base image

Every build, a PR build (`push: false`) too, logs in to ghcr.io with the
job's token, so a base image that is a private ghcr.io package can be
pulled. The calling job grants `packages: read`:

```yaml
jobs:
  image:
    permissions:
      contents: read
      packages: read
    uses: frazer-ai/.github/.github/workflows/image-build.yml@<sha> # v1
    with:
      importpath: ./cmd/frazer-api
```

The package must also grant the calling repository read access. A package
that a workflow of the same repository published already has it. For any
other repository, add it in the package's settings, under Manage Actions
access.

### A vulnerability gate

`vulnerability-gate: true` scans the image's CycloneDX SBOM with Grype
(`grype-version`, built with `go install`) after the SBOM step and before
the attestations, so an image with a finding at `fail-on` (`high`) or above
gets no attestations and no tag, and verification refuses it. A PR build runs the same
gate on an OCI layout of its build, so the PR shows the failure first.
Grype reads the SBOM rather than the image, so a component that no package
database lists, such as a program compiled from source, counts once the
image carries its own CycloneDX fragment and `syft-config` enables Syft's
`sbom-cataloger`:

```yaml
    uses: frazer-ai/.github/.github/workflows/image-build.yml@<sha> # v1
    with:
      working-directory: images/unbound
      dockerfile: Dockerfile
      syft-config: images/unbound/syft.yaml     # select-catalogers: [+sbom-cataloger]
      vulnerability-gate: true
      grype-config: images/unbound/.grype.yaml  # ignore rules, each with a reason
      vex: images/unbound/vex.openvex.json      # not_affected, with a justification
```

The three paths are relative to the repository root. Without
`grype-config`, Grype uses no configuration at all, not one it finds by
itself. The job's `vulnerabilities` and `ignored-vulnerabilities` outputs
list the IDs the gate reported and set aside, and its summary has the table.

Grype's database comes from `grype-db-url`, by default Anchore's listing
(`https://grype.anchore.io/databases`), into a directory of the job's own;
the gate fails when the listing cannot be read or the database is more than
five days old.

### Tags

A publishing build pushes the image by digest only, untagged, and moves its
`tags` (default `latest`, whitespace-separated) to that digest as its last
step: after the vulnerability gate, the attestations and the signature. A tag
therefore never names an image that failed the gate or lacks its
attestations. Neither builder pushes by itself: the build writes an OCI
layout, which crane (`crane-version`, built with `go install`) pushes by
digest and then tags. `tags: ""` leaves the image untagged.

### Renovate

```json
{ "extends": ["local>frazer-ai/.github:renovate-config"] }
```

In `frazer`, also set `"baseBranchPatterns": ["staging"]` so dependency updates
go through staging (F-850). Tool versions in workflows, Makefiles and
`mise.toml` are kept current when annotated:

```yaml
        # renovate: datasource=github-releases depName=golangci/golangci-lint
        default: v2.14.0
```

### Licence exceptions

A component that the licence gate rejects can be allowed in the calling
repository's `.github/license-exceptions.txt`, one per line with a reason:

```
example.com/some/module@v1.2.3 dual-licensed; we use it under MIT (see NOTICE)
```

## Rules for this repository

- Every action is pinned by commit SHA; `zizmor` and `actionlint` run on every PR.
- Changes to `.github/workflows/` need two approvals (CODEOWNERS and the
  ruleset from `org`).
- Releases are tagged `v1`, `v1.1.0` and so on; callers pin SHAs, never tags.
- This repository is public: never put a secret, an internal host name or
  customer data here.
