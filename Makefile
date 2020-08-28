
project_name=wordy
uid=$(shell id -u)
runas=--user ${uid}
mounts=-v ${PWD}:/home/dev/${project_name}:z #-v ${HOME}/.netrc:/home/dev/.netrc
dev_image=${project_name}_dev

.PHONY: help
help: ## Display this help section
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z0-9_-]+:.*?## / {printf "\033[36m%-38s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

.DEFAULT_GOAL := help

checks: ## Lint the project for common errors or bad practices
	golangci-lint run

build: clean ## Build the project
	go build -o wordy main.go

clean: ## Cleanup the project
	rm -f wordy

test: ## Run tests
	go test -v

dev_image: ## Build the development image
	docker build --build-arg uid=${uid} --target dev -t ${dev_image} .

shell: dev_image ## Start a development shell to work on the project
	docker run -it --rm --name ${project_name} ${mounts} ${runas} -w /home/dev/${project_name} ${dev_image}

