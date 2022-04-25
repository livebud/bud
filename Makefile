BUD_VERSION := $(shell cat version.txt)

precommit: test.dev

install:
	@ go install

##
# Examples
##

example.basic:
	@ (cd example/basic && npm link ../../livebud)
	@ go run main.go -C example/basic run

example.basic.watch:
	@ watch -- $(MAKE) example.basic

example.scratch:
	@ rm -fr example/scratch
	@ go run main.go create --link=true example/scratch
	@ go run main.go -C example/scratch new controller / index show
	@ go run main.go -C example/scratch new controller users/admin:admin index show
	@ (cd example/scratch && npm link ../../livebud)
	@ go run main.go -C example/scratch run

example.scratch.watch:
	@ watch -- $(MAKE) example.scratch

example.hn:
	@ (cd example/hn && npm link ../../livebud)
	@ go run main.go -C example/hn run

example.hn.embed:
	@ (cd example/hn && npm link ../../livebud)
	@ go run main.go -C example/hn build --embed
	@ mv example/hn/bud/app $(TMPDIR)/bud_app
	@ $(TMPDIR)/bud_app

example.hn.watch:
	@ watch -- $(MAKE) example.hn

##
# Go
##

GO_SOURCE := ./internal/... ./package/... ./runtime/...
GO_FLAGS := -trimpath -ldflags="-s -w"

go.tools:
	@ go install \
		github.com/evanw/esbuild/cmd/esbuild \
		github.com/pointlander/peg \
		src.techknowlogick.com/xgo

go.mod.tidy:
	@ go mod tidy

# Run go generate
go.generate:
	@ go generate $(GO_SOURCE)

# TODO: add -race back in
go.test:
	@ go test --failfast --timeout=20m $(GO_SOURCE)

go.vet:
	@ go vet $(GO_SOURCE)

go.fmt:
	@ test -z "$(shell go fmt $(GO_SOURCE))"

# Use xgo to cross-compile for OSX, Linux and Windows
go.build.darwin:
	@ xgo --targets=darwin/amd64 --dest=release --out=bud-v$(BUD_VERSION) $(GO_FLAGS) ./ 1> /dev/null

go.build.linux:
	@ xgo --targets=linux/amd64 --dest=release --out=bud-v$(BUD_VERSION) $(GO_FLAGS) ./ 1> /dev/null

# v8go on Windows isn't supported at the moment.
# You'll encounter: "/usr/bin/x86_64-w64-mingw32-ld: cannot find -lv8"
# See:
# - https://github.com/rogchap/v8go#windows
# - https://github.com/rogchap/v8go/pull/234
go.build.windows:
	@ xgo --targets=windows/amd64 --dest=release --out=bud-v$(BUD_VERSION) --trimpath --ldflags="-s -w" ./ 1> /dev/null

##
# BudJS
##

budjs.ci:
	@ (cd livebud && npm ci)

budjs.check:
	@ (cd livebud && ./node_modules/.bin/tsc)

budjs.test:
	@ (cd livebud && ./node_modules/.bin/mocha -r ts-eager/register **/*_test.ts)

##
# Test
##

test: test.dev
test.dev: go.tools go.generate go.fmt go.vet budjs.check go.test budjs.test
test.all: go.tools go.generate go.fmt go.vet budjs.check go.test budjs.test

##
# CI
##

ci.npm:
	@ npm ci

ci.macos: test.all
ci.ubuntu: test.all

##
# Build
##

# TODO windows support
build:
	@ $(MAKE) --no-print-directory -j4 \
		go.build.darwin \
		go.build.linux

##
# Publish
#
# Adapted from esbuild's excellent publish task:
# https://github.com/evanw/esbuild/blob/master/Makefile
#
# TODO: remove Makefile & docs/release.md from the excluded files
##
publish:
	@ npm --version > /dev/null || (echo "The 'npm' command must be in your path to publish" && false)
	@ echo "Checking for uncommitted/untracked changes..." && test -z "`git status --porcelain | grep -vE 'M (CHANGELOG\.md|version\.txt|Makefile|docs/release.md)'`" || \
		(echo "Refusing to publish with these uncommitted/untracked changes:" && \
		git status --porcelain | grep -vE 'M (CHANGELOG\.md|version\.txt|Makefile|docs/release.md)' && false)
	@ echo "Checking for main branch..." && test main = "`git rev-parse --abbrev-ref HEAD`" || \
		(echo "Refusing to publish from non-main branch `git rev-parse --abbrev-ref HEAD`" && false)
	@ echo "Checking for unpushed commits..." && git fetch
	@ test "" = "`git cherry`" || (echo "Refusing to publish with unpushed commits" && false)

	# Build the releases
	@$(MAKE) --no-print-directory build


