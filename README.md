<p align="center"><img src="assets/logo.svg" alt="Lifeboat logo" width="112"></p>

<h1 align="center">lifeboat-github</h1>

<p align="center">The GitHub App adapter that turns verified repository activity into reviewable evidence for Lifeboat.</p>

<p align="center">
  <a href="https://github.com/lifeboat008/lifeboat-github/actions/workflows/ci.yml"><img src="https://github.com/lifeboat008/lifeboat-github/actions/workflows/ci.yml/badge.svg" alt="Go CI"></a>
  <img src="https://img.shields.io/badge/license-MIT-blue" alt="MIT license">
  <img src="https://img.shields.io/badge/go-1.25%2B-00ADD8?logo=go&logoColor=white" alt="Go 1.25+">
  <img src="https://img.shields.io/badge/GitHub-App%20webhooks-181717?logo=github&logoColor=white" alt="GitHub App webhooks">
  <img src="https://img.shields.io/badge/signature-HMAC--SHA256-success" alt="HMAC-SHA256 signatures">
  <img src="https://img.shields.io/badge/Stellar-testnet%20pilot-7D00FF" alt="Stellar testnet pilot">
  <img src="https://img.shields.io/badge/status-pilot-orange" alt="Pilot status">
  <a href="https://cjay-1.gitbook.io/lifeboat-docs/"><img src="https://img.shields.io/badge/docs-GitBook-3884FF" alt="Documentation"></a>
</p>

## Contents

- [What is Lifeboat](#what-is-lifeboat)
- [What this repository does](#what-this-repository-does)
- [Which events become evidence](#which-events-become-evidence)
- [Quick start](#quick-start)
- [Configuration](#configuration)
- [Run it](#run-it)
- [Response codes](#response-codes)
- [How the four repositories fit together](#how-the-four-repositories-fit-together)
- [Project status](#project-status)
- [Open work](#open-work)
- [Documentation](#documentation)
- [Contributing and security](#contributing-and-security)
- [Maintainers](#maintainers)
- [Contributors](#contributors)
- [License](#license)

## What is Lifeboat

Lifeboat helps companies keep the open-source projects they depend on healthy. A sponsor commits a budget to a defined maintenance plan. A maintainer submits evidence of finished work as a claim. A named human reviewer approves or rejects it. Only an approved claim is paid, through a Stellar payment whose transaction hash anyone can verify.

GitHub activity is evidence, never approval. Lifeboat never pays because a project looks inactive.

## What this repository does

| Owns | Does not own |
| --- | --- |
| GitHub App installation identity and webhook HMAC verification | Project funding rules |
| Event delivery de-duplication | Reviewer decisions and payout authorization |
| Mapping supported events to the `lifeboat-api` evidence contract | Stellar payments |

Events only create reviewable evidence. They never approve a claim or trigger a payment.

## Which events become evidence

| GitHub event | Condition | Evidence kind |
| --- | --- | --- |
| `pull_request` | action `closed` and merged | `merged_pr` |
| `release` | action `published` | `release` |
| `issues` | action `closed` | `closed_issue` |

Every other event is acknowledged with `204` and ignored.

The adapter verifies the HMAC-SHA256 signature, the installation ID, the repository, and that the evidence URL belongs to that repository. The API enforces unique delivery IDs, so GitHub retries never create duplicate evidence.

## Quick start

With Go 1.25 or newer:

```bash
git clone https://github.com/lifeboat008/lifeboat-github.git
cd lifeboat-github
go test ./...
go vet ./...
```

The tagged `lifeboat-protocol` module is public, so no module token is needed. The tests use signed synthetic webhooks and need no network.

## Configuration

| Variable | Required | Default | Meaning |
| --- | --- | --- | --- |
| `LIFEBOAT_GITHUB_WEBHOOK_SECRET` | yes | none | GitHub App webhook secret, at least 24 characters |
| `LIFEBOAT_GITHUB_API_TOKEN` | yes | none | Token of a Lifeboat API actor with the `github` role, at least 24 characters |
| `LIFEBOAT_API_URL` | yes | none | Base URL of `lifeboat-api` |
| `LIFEBOAT_GITHUB_PROJECTS_JSON` | yes | none | Map of `owner/repo` to a project binding |
| `LIFEBOAT_GITHUB_LISTEN` | no | `127.0.0.1:8081` | Listen address |

Example project mapping:

```json
{"owner/repo": {"project_id": "p1", "installation_id": 123}}
```

## Run it

```bash
go run ./cmd/lifeboat-github
```

Configure the GitHub App to POST to `/webhook` over HTTPS, with a TLS reverse proxy in front of the loopback listener. Subscribe the App to Pull request, Release, and Issues events and grant it read-only permissions.

## Response codes

| Status | Meaning |
| --- | --- |
| `201` | Evidence recorded |
| `200` | Duplicate delivery acknowledged |
| `204` | Event type or action ignored |
| `400` | Bad delivery ID, body, or evidence URL |
| `401` | Missing or invalid signature |
| `403` | Repository or installation not enrolled |
| `415` | Content type is not JSON |
| `502` | API unreachable or rejected the evidence |

## How the four repositories fit together

```text
lifeboat-protocol ──► lifeboat-ledger ──► lifeboat-api ◄── lifeboat-github
```

| Repository | Responsibility |
| --- | --- |
| [lifeboat-protocol](https://github.com/lifeboat008/lifeboat-protocol) | Domain types, validation, budget arithmetic |
| [lifeboat-ledger](https://github.com/lifeboat008/lifeboat-ledger) | Stellar testnet payment construction, submission, reconciliation |
| [lifeboat-api](https://github.com/lifeboat008/lifeboat-api) | HTTP API, persistence, authorization, audit history |
| [lifeboat-github](https://github.com/lifeboat008/lifeboat-github) | GitHub webhook verification and evidence submission (this repository) |

## Project status

Lifeboat is a **Stellar testnet pilot** with no formal security audit. The repository mapping is static and read once at startup, so adding a repository means restarting the adapter, and App suspension is not handled yet. See [Project status](https://cjay-1.gitbook.io/lifeboat-docs/project-status).

## Open work

- [Handle GitHub App installation suspension and removal](https://github.com/lifeboat008/lifeboat-github/issues/1) (complexity: medium)
- [Replace static repository binding JSON with an auditable enrollment flow](https://github.com/lifeboat008/lifeboat-github/issues/2) (complexity: high)

The product requirements, architecture, and Wave plan are versioned in `lifeboat-api` ([PRD](https://github.com/lifeboat008/lifeboat-api/blob/main/product/docs/PRD.md), [architecture](https://github.com/lifeboat008/lifeboat-api/blob/main/product/docs/ARCHITECTURE.md), [Wave plan](https://github.com/lifeboat008/lifeboat-api/blob/main/product/docs/WAVE.md)).

## Documentation

Full documentation, including [GitHub adapter setup](https://cjay-1.gitbook.io/lifeboat-docs/developer-guide/github-adapter-setup), the API reference, role guides, security notes, and operations runbooks, is at **[cjay-1.gitbook.io/lifeboat-docs](https://cjay-1.gitbook.io/lifeboat-docs/)**.

## Contributing and security

- Read [CONTRIBUTING.md](CONTRIBUTING.md) before opening a pull request. Run `gofmt`, `go vet ./...`, and `go test ./...`, and add tests for changed behavior.
- Use synthetic webhook fixtures only. Never commit secrets or real payloads.
- Report vulnerabilities privately as described in [SECURITY.md](SECURITY.md). Do not post exploit details in a public issue.

## Maintainers

| Maintainer | Contact |
| --- | --- |
| [lifeboat008](https://github.com/lifeboat008) | [Open an issue](https://github.com/lifeboat008/lifeboat-github/issues) for public work; use SECURITY.md for private reports |

## Contributors

<a href="https://github.com/lifeboat008/lifeboat-github/graphs/contributors">
  <img src="https://contrib.rocks/image?repo=lifeboat008/lifeboat-github" alt="Contributors">
</a>

## License

[MIT](LICENSE). This pilot has not had a formal security audit.
