# define variables
GO_VERSION := 1.20
LINTER := golangci-lint
BRANCH_NAME ?=
COMMIT_DIR ?= .
COMMIT_MSG ?= 'Automated commit: $(shell date)'

check-branch:
	@if [ -z "$(BRANCH_NAME)" ]; then \
		echo "Error: BRANCH_NAME is not set."; \
		exit 1; \
	fi

# targets
all: fmt vet test lint push

fmt:
	@echo "Formatting code..."
	go fmt ./...

vet: fmt
	@echo "Running go vet..."
	go vet ./...

test: vet
	@echo "Running go test..."
	go test ./...

lint: test
	@echo "Running static analysis with golangci-lint..."
	$(LINTER) run

push: check-branch
	@echo "Pushing source code to GitHub..."
	# git pull origin $(BRANCH_NAME)
	git add $(COMMIT_DIR)
	git commit -m $(COMMIT_MSG)
	git push origin $(BRANCH_NAME)

# development purpose
run: lint
	@echo "Running ./cmd/api/..."
	go run ./cmd/api/

.PHONY: all fmt vet test lint push check-branch run