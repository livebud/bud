precommit:
	@ go generate ./...
	@ go test ./...

e2e.hackernews:
	@ clear
	@ rm -rf ../hackernews
	@ mkdir ../hackernews
	@ cp -R example/hn/controller ../hackernews
	@ go run main.go -C ../hackernews run

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
