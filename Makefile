.PHONY: test-backend test-frontend test-e2e

test-backend:
	cd backend && go test ./...

test-frontend:
	cd frontend && npm test -- --run

test-e2e:
	cd frontend && npx playwright test
