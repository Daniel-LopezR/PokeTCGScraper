build:
	@go build -o bin/PokeTcgScrapper

run: build
	@./bin/PokeTcgScrapper
