# godel

A download server and maybe a Zapier for downloads.

![Demo](https://raw.githubusercontent.com/KhoalaS/godel/refs/heads/main/media/sc.png)

## Build

### Prerequisites

-   npm >= 10
-   bun >= 1.2
-   go >= 1.25.0

For now you need to install Vue98 and link it use `bun`

```bash
git clone https://github.com/KhoalaS/Vue98.git
cd Vue98
npm install
npm run build
bun link
```

(TBD if it will be used in the future)

Then you can build the frontend and backend using `make`

```bash
git clone https://github.com/KhoalaS/godel.git
cd godel

make install
make build_ui
make build

# to start the server
./build/server

```

If you dont have `make` installed

```bash
cd ui
bun install && bun run build
cd ..
go build -o build/server cmd/server/server.go

```

## Docker

There is also a dockerfile in this project to build the server and UI.

```bash
# build the image
docker build -t godel/server .

# start container
docker-compose up -d
```
