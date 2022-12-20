.PHONY: test
test:
	go test -race -cover ./pkg/lexer && \
	go test -race -cover ./pkg/parser && \
	go test -race -cover ./pkg/eval
