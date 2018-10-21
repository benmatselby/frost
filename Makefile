NAME := frost
DOCKER_PREFIX = benmatselby

.DEFAULT_GOAL := explain
.PHONY: explain
explain:
	### Welcome
	#
	#  ,---.,---.    .---.    .---.  _______
	#  | .-'| .-.\  / .-. )  ( .-._)|__   __|
	#  | `-.| `-'/  | | |(_)(_) \     )| |
	#  | .-'|   (   | | | | _  \ \   (_) |
	#  | |  | |\ \  \ `-' /( `-'  )    | |
	#  )\|  |_| \)\  )---'  `----'     `-'
	# (__)      (__)(_)
	#
	#
	### Installation
	#
	# $$ make all
	#
	### Targets
	@cat Makefile* | grep -E '^[a-zA-Z_-]+:.*?## .*$$' | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

GITCOMMIT := $(shell git rev-parse --short HEAD)

.PHONY: clean
clean: ## Clean the local dependencies
	rm -fr vendor

.PHONY: install
install: ## Install the local dependencies
	dep ensure

.PHONY: vet
vet: ## Vet the code
	go vet ./...

.PHONY: lint
lint: ## Lint the code
	golint -set_exit_status $(shell go list ./...)

.PHONY: build
build: ## Build the application
	go build .

.PHONY: static
static: ## Build the application
	CGO_ENABLED=0 go build -ldflags "-extldflags -static -X github.com/benmatselby/$(NAME)/version.GITCOMMIT=$(GITCOMMIT)" -o $(NAME) .

.PHONY: test
test: ## Run the unit tests
	go test ./... -coverprofile=profile.out

.PHONY: test-cov
test-cov: test ## Run the unit tests with coverage
	go tool cover -html=profile.out

.PHONY: all
all: clean install lint vet build test ## Run everything

.PHONY: static-all
static-all: clean install static test ## Run everything
