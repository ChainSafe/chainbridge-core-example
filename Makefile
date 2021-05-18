
build:
	@echo "  >  \033[32mBuilding binary...\033[0m "
	GOARCH=amd64 go build -o build/chainbridge-example
