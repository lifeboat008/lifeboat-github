<p align="center"><img src="assets/logo.svg" alt="Lifeboat logo" width="112"></p>

# lifeboat-github

Go GitHub App adapter that turns verified repository activity into reviewable evidence candidates.

## Owns

- GitHub App installation identity and webhook HMAC verification.
- Event delivery deduplication.
- Mapping supported events to the `lifeboat-api` evidence contract.

## Does not own

Project funding rules, reviewer decisions, payout authorization, or Stellar payments.

## Run

Set `LIFEBOAT_GITHUB_WEBHOOK_SECRET` to the GitHub App webhook secret, `LIFEBOAT_GITHUB_API_TOKEN` to a Lifeboat API actor token with the `github` role, and `LIFEBOAT_API_URL` to the API base URL. Set `LIFEBOAT_GITHUB_PROJECTS_JSON` to a mapping such as `{"owner/repo":{"project_id":"p1","installation_id":123}}`. Set `LIFEBOAT_GITHUB_LISTEN` for the listen address; its default is `127.0.0.1:8081`. Run `go run ./cmd/lifeboat-github` and configure GitHub to POST to `/webhook` through HTTPS.

The adapter verifies the GitHub HMAC-SHA256 signature, installation ID, repository, and evidence URL. It accepts merged pull requests, published releases, and closed issues. The API enforces unique delivery IDs, so retries do not create duplicate evidence. Events only create reviewable evidence; they never approve a claim or trigger a payment. Go 1.26 and access to the tagged private `lifeboat-protocol` module are required to build from source.

Product PRD and architecture live in the parent `lifeboat/docs` folder in the local workspace.

## Contributing and security

See [CONTRIBUTING.md](CONTRIBUTING.md) for pull requests, [SECURITY.md](SECURITY.md) for private vulnerability reports, and [LICENSE](LICENSE) for MIT terms.
