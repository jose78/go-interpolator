fmt:
	@echo "==> Formating code"
	go fmt

test: all
	@echo "==> Running tests"
	GO111MODULE=on go test -v

test-cover: 
	@echo "==> Running Tests with coverage"
	GO111MODULE=on go test -cover .


