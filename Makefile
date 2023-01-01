.PHONY: test
test:
	go test -race -v -cover ./pkg/lexer && \
	go test -race -v -cover ./pkg/parser && \
	go test -race -v -cover ./pkg/eval && \
	go test -race -v -cover ./pkg/bytecode && \
	go test -race -v -cover ./pkg/compiler && \
	go test -race -v -cover ./pkg/vm 
