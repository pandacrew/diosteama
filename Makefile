project=diosteama
remote=fary.pandacrew.net
source=${project}.go

all: build

${project}:
	go build .

sync: build ## Deploy to remote server
	scp ${project} ${remote}:

build: ${project} ## Build project

run: ## Run without building
	go run .

clean: ## Delete generated files (so far the generated executable file)
	rm ${project}

help:   ## Shows this message.
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: all sync build run clean help
