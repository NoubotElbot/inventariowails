VITE_PORT ?= 9245

.PHONY: dev build package test seed fmt fmt-check install-frontend

# ── Development ──────────────────────────────────────────────────────────────

## dev: Run the application in development mode (live reload)
dev:
	wails3 dev -config ./build/config.yml -port $(VITE_PORT)

# ── Production ───────────────────────────────────────────────────────────────

## build: Build the desktop application for the host platform
build:
	wails3 build

## package: Create a distributable package for the host platform
package:
	wails3 package

# ── Database ─────────────────────────────────────────────────────────────────

## seed: Seed the database with sample categories and products
seed:
	go run ./cmd/seed

# ── Testing ───────────────────────────────────────────────────────────────────

## test: Run all Go tests
test:
	go test -v .\internal\modules\...\application .\internal\modules\...\domain

# ── Frontend ─────────────────────────────────────────────────────────────────

## install-frontend: Install frontend dependencies
install-frontend:
	cd frontend && pnpm install

## fmt: Format frontend source files
fmt:
	cd frontend && pnpm run format

## fmt-check: Check frontend formatting without writing
fmt-check:
	cd frontend && pnpm run format:check

# ── Help ──────────────────────────────────────────────────────────────────────

## help: List available commands
help:
	@grep -E '^## ' Makefile | sed 's/## /  /'

.DEFAULT_GOAL := help
