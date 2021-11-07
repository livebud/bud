precommit:
	@ go generate ./...
	@ go test ./...

clear:
	@ clear

hackernews.run:
	@ rm -rf ../hackernews/run
	@ mkdir ../hackernews/run
	@ go run main.go -C ../hackernews/run run

hackernews.build:
	@ rm -rf ../hackernews/build
	@ mkdir -p ../hackernews/build
	@ go run main.go -C ../hackernews/build build
	@ cd ../hackernews/build && ./bud/main

hackernews.deploy:
	@ rm -rf ../hackernews/deploy
	@ mkdir -p ../hackernews/deploy
	@ go run main.go -C ../hackernews/deploy deploy

example:
	@ watch -- $(MAKE) clear hackernews.run hackernews.build hackernews.deploy
