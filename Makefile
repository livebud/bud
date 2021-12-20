precommit:
	@ go generate ./...
	@ go test ./...

e2e.hackernews:
	@ clear
	@ rm -rf ../hackernews
	@ cp -R example/hn ../hackernews
	@ go run main.go -C ../hackernews deploy --access-key=1 --secret-key=2
	@ go run main.go -C ../hackernews new view users
	@ go run main.go -C ../hackernews run

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
