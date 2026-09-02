# Open Prompt Registry

A Go service for governing prompts used with an Open Guardrail Layer. It provides immutable prompt versions and an explicit release workflow:

`draft -> in_review -> approved -> rolled_out`

A version cannot be approved until it has at least one passing test. The rollout operation marks the selected version as the prompt's production version.

## Run

```sh
go run ./cmd/prompt-registry
```

The server listens on `http://localhost:8080`.

## API

| Operation | Endpoint |
| --- | --- |
| Health check | `GET /healthz` |
| Create a prompt | `POST /prompts` |
| Read prompt and versions | `GET /prompts/{id}` |
| Create a draft version | `POST /prompts/{id}/versions` |
| Record a test | `POST /prompts/{id}/versions/{number}/tests` |
| Request review | `POST /prompts/{id}/versions/{number}/review` |
| Approve a version | `POST /prompts/{id}/versions/{number}/approve` |
| Roll out a version | `POST /prompts/{id}/versions/{number}/rollout` |

### Example workflow

```sh
curl -X POST localhost:8080/prompts -H 'Content-Type: application/json' \
  -d '{"name":"assistant","description":"Main support prompt"}'

curl -X POST localhost:8080/prompts/prompt-1/versions -H 'Content-Type: application/json' \
  -d '{"template":"Help with: {{question}}","author":"author@example.com"}'

curl -X POST localhost:8080/prompts/prompt-1/versions/1/tests -H 'Content-Type: application/json' \
  -d '{"name":"prompt-injection suite","passed":true}'

curl -X POST localhost:8080/prompts/prompt-1/versions/1/review -H 'Content-Type: application/json' \
  -d '{"actor":"reviewer@example.com"}'

curl -X POST localhost:8080/prompts/prompt-1/versions/1/approve -H 'Content-Type: application/json' \
  -d '{"actor":"approver@example.com"}'

curl -X POST localhost:8080/prompts/prompt-1/versions/1/rollout
```

## Test

```sh
go test ./...
```

This first implementation stores state in memory. Replace `registry.Service` storage with a durable repository when deploying it across multiple instances.
