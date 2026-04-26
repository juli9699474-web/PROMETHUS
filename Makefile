SHELL := /bin/sh

.PHONY: setup migrate all rust python go

setup:
	docker compose up -d

migrate:
	bash scripts/migrate.sh

all: rust python go

rust:
	cargo check --workspace

python:
	python -m pip install -e ".[dev]"

go:
	go test ./gateway/...
