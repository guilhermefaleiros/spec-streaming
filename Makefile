.PHONY: test-backend test-frontend test-e2e run-api run-worker verify db-up db-migrate

db-up:
	docker compose up -d postgres flyway

db-migrate:
	docker compose up flyway

test-backend:
	cd backend && go test ./...

test-frontend:
	cd frontend && npm test -- --run

test-e2e:
	cd frontend && npx playwright test

run-api:
	cd backend && go run ./cmd/api

run-worker:
	cd backend && go run ./cmd/worker

verify:
	$(MAKE) test-backend
	$(MAKE) test-frontend
