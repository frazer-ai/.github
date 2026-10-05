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
| `.github/workflows/go-check.yml` | Go gate: gofmt, `go mod tidy -diff`, vet, golangci-lint, race tests (shuffled, uncached), govulncheck |
| `.github/workflows/web-check.yml` | Web gate: `pnpm install --frozen-lockfile`, then `typecheck lint format:check test build` when defined |
| `.github/workflows/tofu-check.yml` | OpenTofu gate: fmt, init without a backend, validate, TFLint, Trivy config scan (digest-pinned, signature-verified image), `tofu test` |
| `.github/workflows/image-build.yml` | The one image builder: ko, multi-arch; on push, provenance and SBOM attestations; signing with the OpenBao transit key |
| `.github/workflows/scan.yml` | Security scanners: Betterleaks (a pull request's own commits; a pushed branch's history; every ref with `all-refs: true` for a scheduled run), OSV-Scanner (lockfiles), Trivy (vulnerabilities and misconfigurations) from its digest-pinned, signature-verified image |
| `.github/workflows/names.yml` | Naming gates: no other cloud's region IDs in tracked files, Conventional Commit PR titles |
| `.github/workflows/license-gate.yml` | Licence gate: Syft SBOM, SPDX allow-list, exceptions file, generated NOTICE |
| `.github/workflows/self-check.yml` | This repository's CI: lints the workflows and runs each one against `testdata/`, including fixtures that must fail (`license-bad`, `vuln`) |

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
