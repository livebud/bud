test: precommit

precommit:
	@ go generate ./...
	@ go test ./...

install:
	@ go install

example.basic:
	@ go run main.go -C example/basic run

example.basic.watch:
	@ watch -- $(MAKE) example.basic

example.hello:
	@ rm -fr example/hello
	@ go run main.go create --link=true example/hello
	@ go run main.go -C example/hello new controller / index show
	@ go run main.go -C example/hello new controller users/admin:admin index show
	@ go run main.go -C example/hello run

example.hello.watch:
	@ watch -- $(MAKE) example.hello

example.hn:
	@ go run main.go -C example/hn run

example.hn.embed:
	@ rm -fr example/hn/bud
	@ go run main.go -C example/hn build --embed
	@ mv example/hn/bud/app $(TMPDIR)/bud_app
	@ $(TMPDIR)/bud_app

example.hn.watch:
	@ watch -- $(MAKE) example.hn
