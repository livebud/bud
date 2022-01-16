test: precommit

precommit:
	@ go generate ./...
	@ go test ./...

TAILWIND := $(realpath $(PWD)/../bud-tailwind)

e2e.hackernews:
	@ clear
	@ rm -rf _build/hn
	@ mkdir -p _build
	@ go install
	@ cp -R example/hn _build/hn
	@ go mod edit -replace=gitlab.com/mnm/bud=$(PWD) _build/hn/go.mod
	@ go mod edit -replace=gitlab.com/mnm/bud-tailwind=$(TAILWIND) _build/hn/go.mod
	@ bud -C _build/hn deploy --access-key=1 --secret-key=2
	@ bud -C _build/hn new view users
	@ bud -C _build/hn run

e2e.hackernews.run:
	@ (cd ../hackernews && go run bud/main.go)

# hackernews.build:
# 	# @ rm -rf ../hackernews
# 	# @ mkdir -p ../hackernews
# 	@ go run main.go -C ../hackernews build
# 	@ cd ../hackernews && ./bud/main

# hackernews.deploy:
# 	# @ rm -rf ../hackernews
# 	# @ mkdir -p ../hackernews
# 	@ go run main.go -C ../hackernews deploy

example:
	@ watch -- $(MAKE) e2e.hackernews

action:
	@ watch -- $(MAKE) test.action

test.action:
	@ clear
	@ go test ./internal/generator/action

.PHONY: example action
