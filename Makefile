GO=go


build: 
	$(GO) build -o bin/* *.go

run:
	sudo $(GO) run *.go

