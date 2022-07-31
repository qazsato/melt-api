.PHONY: init serve build clean deploy gomodgen

init:
	npm install

serve:
	npm run offline

clean:
	rm -rf ./bin ./vendor Gopkg.lock

deploy: clean build
	npm run deploy:prod

gomodgen:
	chmod u+x gomod.sh
	./gomod.sh
