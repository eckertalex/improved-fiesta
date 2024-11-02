include .env

# ==================================================================================== #
# HELPERS
# ==================================================================================== #

## help: print this help message
.PHONY: help
help:
	@echo "Usage:"
	@sed -n "s/^##//p" ${MAKEFILE_LIST} | column -t -s ":" |  sed -e "s/^/ /"

.PHONY: confirm
confirm:
	@echo "Are you sure? (y/n) \c"
	@read answer; \
	if [ "$$answer" != "y" ]; then \
		echo "Aborting."; \
		exit 1; \
	fi

# ==================================================================================== #
# QUALITY CONTROL
# ==================================================================================== #

## audit: run quality control checks
.PHONY: audit
audit: test
	go mod tidy -diff
	go mod verify
	test -z "$(shell gofmt -l .)" 
	go vet ./...
	go run honnef.co/go/tools/cmd/staticcheck@latest -checks=all,-ST1000,-U1000 ./...
	go run golang.org/x/vuln/cmd/govulncheck@latest ./...

## test: run all tests
.PHONY: test
test:
	go test -v -race -buildvcs ./...

## itest: run integration tests
.PHONY: itest
itest:
	./scripts/api-smoke-test.sh

# ==================================================================================== #
# DEVELOPMENT
# ==================================================================================== #

## tidy: tidy and format all .go files
.PHONY: tidy
tidy:
	go mod tidy
	go fmt ./...

## build/api: build the cmd/api application
.PHONY: build/api
build/api:
	go build -o=./tmp/api ./cmd/api

## run/api: run the cmd/api application
.PHONY: run/api
run/api: build/api
	@./tmp/api \
		-port=${IMPROVED_FIESTA_PORT} \
		-env=${IMPROVED_FIESTA_ENV} \
		-db-dsn=${IMPROVED_FIESTA_DB_DSN} \
		-smtp-host=${IMPROVED_FIESTA_SMTP_HOST} \
		-smtp-port=${IMPROVED_FIESTA_SMTP_PORT} \
		-smtp-username=${IMPROVED_FIESTA_SMTP_USERNAME} \
		-smtp-password=${IMPROVED_FIESTA_SMTP_PASSWORD} \
		-smtp-sender=${IMPROVED_FIESTA_SMTP_SENDER}

## watch/api: watch the cmd/api application
.PHONY: watch/api
watch/api:
	@go run github.com/air-verse/air@latest \
		--build.cmd "make build/api" \
		--build.bin "./tmp/api \
			-port=${IMPROVED_FIESTA_PORT} \
			-env=${IMPROVED_FIESTA_ENV} \
			-db-dsn=${IMPROVED_FIESTA_DB_DSN} \
			-smtp-host=${IMPROVED_FIESTA_SMTP_HOST} \
			-smtp-port=${IMPROVED_FIESTA_SMTP_PORT} \
			-smtp-username=${IMPROVED_FIESTA_SMTP_USERNAME} \
			-smtp-password=${IMPROVED_FIESTA_SMTP_PASSWORD} \
			-smtp-sender=${IMPROVED_FIESTA_SMTP_SENDER}"
	
# ==================================================================================== #
# DB
# ==================================================================================== #

## db/seed: seed database
.PHONY: db/seed
db/seed: confirm
	go run ./cmd/seed -database=${IMPROVED_FIESTA_DB_DSN}

## db/connect: create to the local database
.PHONY: db/connect
db/connect:
	sqlite3 ${IMPROVED_FIESTA_DB_DSN}

## db/migrations/new name=$1: create a new migration
.PHONY: db/migrations/new
db/migrations/new:
	go run -tags 'sqlite3' github.com/golang-migrate/migrate/v4/cmd/migrate@latest create -seq -ext=.sql -dir=./migrations ${name}

## db/migrations/up: apply all up migrations
.PHONY: db/migrations/up
db/migrations/up: confirm
	go run -tags 'sqlite3' github.com/golang-migrate/migrate/v4/cmd/migrate@latest -path=./migrations -database="sqlite3://${IMPROVED_FIESTA_DB_DSN}" up

## db/migrations/down: apply all down migrations
.PHONY: db/migrations/down
db/migrations/down: confirm
	go run -tags 'sqlite3' github.com/golang-migrate/migrate/v4/cmd/migrate@latest -path=./migrations -database="sqlite3://${IMPROVED_FIESTA_DB_DSN}" down

## db/migrations/goto version=$1: migrate to a specific version number
.PHONY: db/migrations/goto
db/migrations/goto: confirm
	go run -tags 'sqlite3' github.com/golang-migrate/migrate/v4/cmd/migrate@latest -path=./migrations -database="sqlite3://${IMPROVED_FIESTA_DB_DSN}" goto ${version}

## db/migrations/force version=$1: force database migration version number
.PHONY: db/migrations/force
db/migrations/force: confirm
	go run -tags 'sqlite3' github.com/golang-migrate/migrate/v4/cmd/migrate@latest -path=./migrations -database="sqlite3://${IMPROVED_FIESTA_DB_DSN}" force ${version}

## db/migrations/version: print the current migration version
.PHONY: db/migrations/version
db/migrations/version:
	go run -tags 'sqlite3' github.com/golang-migrate/migrate/v4/cmd/migrate@latest -path=./migrations -database="sqlite3://${IMPROVED_FIESTA_DB_DSN}" version

# vim: set tabstop=4 shiftwidth=4 noexpandtab
