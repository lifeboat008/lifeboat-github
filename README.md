<p align="center"><img src="assets/logo.svg" alt="Lifeboat logo" width="112"></p>

# lifeboat-github

[![Go CI](https://github.com/lifeboat008/lifeboat-github/actions/workflows/ci.yml/badge.svg)](https://github.com/lifeboat008/lifeboat-github/actions/workflows/ci.yml)

Go GitHub App adapter that turns verified repository activity into reviewable evidence candidates.

## Owns

- GitHub App installation identity and webhook HMAC verification.
- Event delivery deduplication.
- Mapping supported events to the `lifeboat-api` evidence contract.

## Does not own

Project funding rules, reviewer decisions, payout authorization, or Stellar payments.

## Run

Set `LIFEBOAT_GITHUB_WEBHOOK_SECRET` to the GitHub App webhook secret, `LIFEBOAT_GITHUB_API_TOKEN` to a Lifeboat API actor token with the `github` role, and `LIFEBOAT_API_URL` to the API base URL. Set `LIFEBOAT_GITHUB_PROJECTS_JSON` to a mapping such as `{"owner/repo":{"project_id":"p1","installation_id":123}}`. Set `LIFEBOAT_GITHUB_LISTEN` for the listen address; its default is `127.0.0.1:8081`. Run `go run ./cmd/lifeboat-github` and configure GitHub to POST to `/webhook` through HTTPS.

The adapter verifies the GitHub HMAC-SHA256 signature, installation ID, repository, and evidence URL. It accepts merged pull requests, published releases, and closed issues. The API enforces unique delivery IDs, so retries do not create duplicate evidence. Events only create reviewable evidence; they never approve a claim or trigger a payment. Go 1.25 or newer is required. The tagged `lifeboat-protocol` module is public; no module token is needed. Run `go test ./...` and `go vet ./...` before opening a pull request.

The [product requirements](https://github.com/lifeboat008/lifeboat-api/blob/main/product/docs/PRD.md), [architecture](https://github.com/lifeboat008/lifeboat-api/blob/main/product/docs/ARCHITECTURE.md), and [Wave plan](https://github.com/lifeboat008/lifeboat-api/blob/main/product/docs/WAVE.md) are versioned in `lifeboat-api`.

## Documentation

Full documentation, including the API reference, role guides, security notes, and operations runbooks, is at [cjay-1.gitbook.io/lifeboat-docs](https://cjay-1.gitbook.io/lifeboat-docs/).

## Contributing and security

See [CONTRIBUTING.md](CONTRIBUTING.md) for pull requests, [SECURITY.md](SECURITY.md) for private vulnerability reports, and [LICENSE](LICENSE) for MIT terms.

Maintainers: [lifeboat008](https://github.com/lifeboat008). Discuss public work in [issues](https://github.com/lifeboat008/lifeboat-github/issues); report vulnerabilities privately as described in SECURITY.md. This pilot has not had a formal security audit.
