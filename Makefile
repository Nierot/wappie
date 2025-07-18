run: generate
	go run .

dev: generate
	air

build: generate
	go build
	
generate:
	go generate