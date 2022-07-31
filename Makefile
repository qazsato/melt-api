.PHONY: init serve build clean deploy gomodgen

init:
	npm install
	@make build

serve:
	npm run offline

build: gomodgen
	export GO111MODULE=on
	env GOOS=linux go build -ldflags="-s -w" -o bin/apps functions/apps/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/sites functions/sites/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/images functions/images/main.go

clean:
	rm -rf ./bin ./vendor Gopkg.lock

deploy: clean build
	npm run deploy:prod

gomodgen:
	chmod u+x gomod.sh
	./gomod.sh
