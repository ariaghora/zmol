.PHONY: test
test:
	go test -v -cover ./pkg/lexer && \
	go test -v -cover ./pkg/parser
