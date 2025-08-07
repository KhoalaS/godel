cd ui
bun run build
cd ..
go build -o build/server.exe cmd/server/server.go
