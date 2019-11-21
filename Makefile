all:
	mkdir -p target/script/
	go build github.com/ailisp/reallyfastci/cmd/rfci
	mv rfci target/
	cp -r script/*.py script/Pipfile script/Pipfile.lock target/script/
	cp config/config.yaml.example target/config.yaml
	cp script/build.sh.example target/build.sh

clean:
	rm -rf target

test:
	go test ./...

.PHONY: all test clean