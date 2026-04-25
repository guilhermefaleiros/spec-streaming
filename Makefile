.PHONY: test-backend test-frontend test-e2e run-api run-worker verify db-up db-migrate dev

# --- Database ---

db-up:
	docker compose up -d postgres flyway

db-migrate:
	docker compose up flyway

db-down:
	docker compose down

# --- Development ---

dev: db-up
	@echo ""
	@echo "Infrastructure is up! Now run these in separate terminals:"
	@echo "  make run-api"
	@echo "  make run-worker"
	@echo "  make run-frontend"
	@echo ""

run-api:
	@echo "Starting API on http://localhost:8080"
	cd backend && go run ./cmd/api

run-worker:
	@echo "Starting Worker..."
	cd backend && go run ./cmd/worker

run-frontend:
	@echo "Starting Frontend on http://localhost:5173"
	cd frontend && npm run dev

# --- Tests ---

test-backend:
	cd backend && go test ./...

test-frontend:
	cd frontend && npm test -- --run

test-e2e:
	cd frontend && npx playwright test

verify:
	$(MAKE) test-backend
	$(MAKE) test-frontend
