GO=go


build: 
	$(GO) build -o bin/main main.go

run:
	sudo $(GO) run main.go

