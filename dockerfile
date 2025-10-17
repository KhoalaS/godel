FROM ubuntu:24.04 AS builder

RUN apt update && apt install -y curl tar git bash unzip build-essential

RUN curl -fsSL https://deb.nodesource.com/setup_22.x | bash && \
    apt install -y nodejs && \
    npm install -g npm@latest bun

# Set Go version
ENV GO_VERSION=1.25.1
RUN curl -LO https://go.dev/dl/go${GO_VERSION}.linux-amd64.tar.gz && \
    tar -C /usr/local -xzf go${GO_VERSION}.linux-amd64.tar.gz && \
    rm go${GO_VERSION}.linux-amd64.tar.gz
ENV PATH=/usr/local/go/bin:$PATH
ENV GOPATH=/go
ENV PATH=$GOPATH/bin:$PATH

# Set Bun path
ENV PATH=/root/.bun/bin:$PATH

WORKDIR /docker/app
RUN git clone https://github.com/KhoalaS/Vue98.git
WORKDIR /docker/app/Vue98
RUN npm install && npm run build && bun link

COPY ui/package.json ui/bun.lock /docker/app/godel/ui/
WORKDIR /docker/app/godel/ui
RUN bun install
COPY ./ /docker/app/godel
RUN bun run build

WORKDIR /docker/app/godel
RUN go mod download
ENV GOEXPERIMENT=greenteagc
RUN go build -ldflags="-s -w" -o build/server cmd/server/server.go

FROM alpine:latest AS wasmer
RUN apk update && apk add tar curl
WORKDIR /tmp/wasmer
RUN curl -sSfL https://github.com/wasmerio/wasmer/releases/download/v6.1.0/wasmer-linux-amd64.tar.gz -o wasmer.tar.gz && \
    tar xvf wasmer.tar.gz

FROM ubuntu:24.04 AS runtime
RUN apt update && apt install -y curl tar ca-certificates --no-install-recommends && rm -rf /var/lib/apt/lists/*

COPY --from=wasmer /tmp/wasmer/lib/libwasmer.so /usr/lib/

WORKDIR /app
COPY --from=builder /docker/app/godel/build/server ./
COPY .env ./

RUN useradd -r -U godel && chown -R godel:godel /app
USER godel

ENTRYPOINT [ "/app/server" ]