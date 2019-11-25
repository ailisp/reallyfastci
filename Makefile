all: go python config

go:
	mkdir -p target
	go build github.com/ailisp/reallyfastci/cmd/rfci
	mv rfci target/

python:
	mkdir -p target/script
	cp -r script/*.py target/script/
	cp -r script/Pipfile.lock script/Pipfile target/
	cd target && pipenv install

config:
	mkdir -p target
	cp config/config.yaml.example target/config.yaml
	cp script/build.sh.example target/build.sh

clean:
	rm -rf target

test:
	go test ./...

.PHONY: all test clean go python config