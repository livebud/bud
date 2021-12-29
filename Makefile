test: precommit

precommit:
	@ go generate ./...
	@ go test ./...

TAILWIND := $(realpath $(PWD)/../bud-tailwind)

e2e.hackernews:
	# @ clear
	@ rm -rf build/hn
	@ mkdir -p build
	@ go install
	@ cp -R example/hn build/hn
	@ go mod edit -replace=gitlab.com/mnm/bud=$(PWD) build/hn/go.mod
	@ go mod edit -replace=gitlab.com/mnm/bud-tailwind=$(TAILWIND) build/hn/go.mod
	# @ bud -C build/hn deploy --access-key=1 --secret-key=2
	# @ go run main.go -C build/hn new view users
	@ go run main.go -C build/hn run

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

.PHONY: example
