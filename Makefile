#build
build:
	@go build -o bin/lumina-api cmd/main.go

#run
run: build
	@./bin/lumina-api