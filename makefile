# Makefile (Local Version)

.PHONY: lint lint-backend lint-frontend test test-backend test-frontend format format-backend format-frontend

# ==============================================================================
# LINTING TASKS (CHECKS ONLY)
# ==============================================================================

## ✨ Runs linters to check for issues
lint: lint-backend lint-frontend

## Checks backend code for issues
lint-backend:
	@echo "--- Checking Backend (Go) [LOCAL] ---"
	@(cd backend && golangci-lint run ./...)

## Checks frontend code for issues
lint-frontend:
	@echo "--- Checking Frontend (Next.js) [LOCAL] ---"
	@(cd frontend && pnpm lint)

# ==============================================================================
# FORMATTING TASKS (AUTO-FIXES ISSUES)
# ==============================================================================

## 💅 Formats and fixes all code automatically
format: format-backend format-frontend

## Formats and fixes backend code
format-backend:
	@echo "--- Formatting Backend (Go) [LOCAL] ---"
	@(cd backend && golangci-lint run --fix ./... && gofmt -w ./)
## Formats and fixes frontend code
format-frontend:
	@echo "--- Formatting Frontend (Next.js) [LOCAL] ---"
	@(cd frontend && pnpm prettier --write . && pnpm eslint . --fix)


# ==============================================================================
# TESTING TASKS (RUNS INSIDE DOCKER CONTAINERS)
# ==============================================================================

## 🧪 Runs tests for all services inside containers
test: test-backend test-frontend

## Go tests for the backend service
test-backend:
	@echo "--- Testing Backend (Go) [DOCKER] ---"
	@docker-compose exec backend go test -v ./...

## Jest/Vitest tests for the frontend service
test-frontend:
	@echo "--- Testing Frontend (Next.js) [DOCKER] ---"
	@docker-compose exec frontend pnpm test
