precommit: test.dev

install:
	@ go install

##
# Examples
##

example.basic:
	@ (cd livebud && npm link --quiet --no-progress)
	@ (cd example/basic && npm link livebud)
	@ go run main.go -C example/basic run

example.basic.watch:
	@ watch -- $(MAKE) example.basic

example.scratch:
	@ rm -fr example/scratch
	@ go run main.go create --link=true example/scratch
	@ go run main.go -C example/scratch new controller / index show
	@ go run main.go -C example/scratch new controller users/admin:admin index show
	@ (cd livebud && npm link --quiet --no-progress)
	@ (cd example/scratch && npm link livebud)
	@ go run main.go -C example/scratch run

example.scratch.watch:
	@ watch -- $(MAKE) example.scratch

example.hn:
	@ (cd livebud && npm link --quiet --no-progress)
	@ (cd example/hn && npm link livebud)
	@ go run main.go -C example/hn run

example.hn.embed:
	@ (cd livebud && npm link --quiet --no-progress)
	@ (cd example/hn && npm link livebud)
	@ go run main.go -C example/hn build --embed
	@ mv example/hn/bud/app $(TMPDIR)/bud_app
	@ $(TMPDIR)/bud_app

example.hn.watch:
	@ watch -- $(MAKE) example.hn

##
# Go
##

GO_SOURCE := ./internal/... ./package/... ./runtime/...

# Run go generate
go.generate:
	@ go generate $(GO_SOURCE)

# TODO: add -race back in
go.test:
	@ go test $(GO_SOURCE)

go.vet:
	@ go vet $(GO_SOURCE)

go.fmt:
	@ test -z "$(shell go fmt $(GO_SOURCE))"

##
# Test
##

test:
	@ $(MAKE) --no-print-directory -j6 test.dev

test.dev: go.test go.vet go.fmt
test.ci: go.test go.vet go.fmt

##
# Build
##

build.darwin:


build.linux:

##
# Publish
##

publish:
	@ npm --version > /dev/null || (echo "The 'npm' command must be in your path to publish" && false)
	@ echo "Checking for uncommitted/untracked changes..." && test -z "`git status --porcelain | grep -vE 'M (CHANGELOG\.md|version\.txt)'`" || \
		(echo "Refusing to publish with these uncommitted/untracked changes:" && \
		git status --porcelain | grep -vE 'M (CHANGELOG\.md|version\.txt)' && false)
	@ echo "Checking for main branch..." && test main = "`git rev-parse --abbrev-ref HEAD`" || \
		(echo "Refusing to publish from non-main branch `git rev-parse --abbrev-ref HEAD`" && false)
	@ echo "Checking for unpushed commits..." && git fetch
	@ test "" = "`git cherry`" || (echo "Refusing to publish with unpushed commits" && false)

