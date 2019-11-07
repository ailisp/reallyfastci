all:
	go build github.com/ailisp/reallyfastci/cmd/rfci

clean:
	rm -f rfci

test:
	go test ./...

.PHONY: all test clean