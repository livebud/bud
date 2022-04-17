test: precommit

precommit:
	@ go generate ./...
	@ go test ./...

install:
	@ go install

npm.link.livebud:
	@ (cd livebud && npm link --quiet --no-progress)

example.basic: npm.link.livebud
	@ (cd example/basic && npm link livebud)
	@ go run main.go -C example/basic run

example.basic.watch:
	@ watch -- $(MAKE) example.basic

example.scratch: npm.link.livebud
	@ rm -fr example/scratch
	@ go run main.go create --link=true example/scratch
	@ go run main.go -C example/scratch new controller / index show
	@ go run main.go -C example/scratch new controller users/admin:admin index show
	@ (cd example/scratch && npm link livebud)
	@ go run main.go -C example/scratch run

example.scratch.watch:
	@ watch -- $(MAKE) example.scratch

example.hn: npm.link.livebud
	@ (cd example/hn && npm link livebud)
	@ go run main.go -C example/hn run

example.hn.embed: npm.link.livebud
	@ (cd example/hn && npm link livebud)
	@ go run main.go -C example/hn build --embed
	@ mv example/hn/bud/app $(TMPDIR)/bud_app
	@ $(TMPDIR)/bud_app

example.hn.watch:
	@ watch -- $(MAKE) example.hn
