FROM ubuntu:24.04 AS builder
RUN apt update
RUN apt install -y curl tar git bash unzip
RUN curl -fsSL https://deb.nodesource.com/setup_22.x | bash
RUN apt install -y nodejs

RUN npm install -g npm@latest &&\
    npm install -g bun

# Set Go version
ENV GO_VERSION=1.25.1

# Download and install Go
RUN curl -LO https://go.dev/dl/go${GO_VERSION}.linux-amd64.tar.gz \
    && tar -C /usr/local -xzf go${GO_VERSION}.linux-amd64.tar.gz \
    && rm go${GO_VERSION}.linux-amd64.tar.gz

# Set Go environment
ENV PATH=/usr/local/go/bin:$PATH
ENV GOPATH=/go
ENV PATH=$GOPATH/bin:$PATH

# Set Bun version
ENV BUN_VERSION=1.2.16

# Install Bun
RUN curl -fsSL https://bun.sh/install | bash

# Add Bun to PATH
ENV PATH=/root/.bun/bin:$PATH

WORKDIR /docker/app
RUN git clone https://github.com/KhoalaS/Vue98.git

WORKDIR /docker/app/Vue98
RUN npm install &&\
    npm run build &&\
    bun link

COPY ui/package.json ui/bun.lock /docker/app/godel/ui/

WORKDIR /docker/app/godel/ui 
RUN bun install
COPY ./ /docker/app/godel
RUN bun run build

WORKDIR /docker/app/godel
RUN apt -y install build-essential
RUN go mod download
RUN go build -o build/server cmd/server/server.go

FROM ubuntu:24.04 AS runtime 
RUN apt update
RUN apt install -y curl tar passwd

WORKDIR /tmp/wasmer
RUN curl -sSfL https://github.com/wasmerio/wasmer/releases/download/v6.1.0/wasmer-linux-amd64.tar.gz -o wasmer.tar.gz
RUN tar xvf wasmer.tar.gz && cp lib/libwasmer.so /usr/lib

WORKDIR /app
RUN rm -rf /tmp/wasmer
COPY --from=builder /docker/app/godel/build/server ./
COPY .env ./

RUN useradd -r -U godel
RUN chown -R godel:godel /app
USER godel

ENTRYPOINT [ "/app/server" ]