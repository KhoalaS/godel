install:
	cd ui && bun install
	go mod tidy
	
build: worker.go cmd/**/*.go pkg/**/*.go ui/dist/**/*
	go build -o build/server.exe cmd/server/server.go

build_ui: ui/src/**/*.vue ui/src/**/*.ts
	cd ui && bun run build

clean:
	rm build/server.exe
	rm -r ui/dist

build_all: build_ui build