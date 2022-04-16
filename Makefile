test: precommit

install:
	@ go install

precommit:
	@ go generate ./...
	@ go test ./...

TAILWIND := $(realpath $(PWD)/../bud-tailwind)

e2e.basic:
	@ clear
	@ go run main.go -C example/basic run

basic:
	@ watch -- $(MAKE) e2e.basic

e2e.hello:
	# @ clear
	@ rm -fr example/hello
	@ go run main.go create --link=true example/hello
	@ go run main.go -C example/hello new controller / index show
	@ go run main.go -C example/hello new controller users/admin:admin index show
	@ go run main.go -C example/hello run

hello:
	@ watch -- $(MAKE) e2e.hello

e2e.hn:
	@ go run main.go -C example/hn run

e2e.hn.embed:
	@ rm -fr example/hn/bud
	@ go run main.go -C example/hn build --embed
	@ mv example/hn/bud/app $(TMPDIR)/bud_app
	@ $(TMPDIR)/bud_app

hn:
	@ watch -- $(MAKE) e2e.hn

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

action:
	@ watch -- $(MAKE) test.action

test.action:
	@ clear
	@ go test ./internal/generator/action

.PHONY: example action
