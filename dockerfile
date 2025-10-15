FROM alpine:latest AS builder

RUN apk update && apk upgrade
RUN apk add curl tar git bash npm
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
RUN git clone https://github.com/KhoalaS/godel.git
RUN git clone https://github.com/KhoalaS/Vue98.git

WORKDIR /docker/app/Vue98
RUN npm install &&\
    npm run build &&\
    bun link

WORKDIR /docker/app/godel
RUN git fetch &&\
    git pull origin

RUN cd ui && bun install && bun run build
RUN go build -o build/server cmd/server/server.go

FROM alpine:latest

WORKDIR /app
COPY --from=builder /docker/app/godel/build/server ./
COPY .env ./

RUN addgroup --system godel && adduser --system --ingroup godel godel
RUN chown -R godel:godel /app
USER godel

ENTRYPOINT [ "/app/server" ]