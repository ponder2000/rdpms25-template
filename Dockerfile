# Set a build-time variable
ARG APP_NAME=app
ARG RELEASE_VERSION=unknown

FROM debian:bullseye AS builder

ARG APP_NAME
ARG RELEASE_VERSION

ENV APP_NAME=${APP_NAME}
ENV RELEASE_VERSION=${RELEASE_VERSION}

RUN apt-get update && apt-get install -y \
    wget \
    build-essential \
    git \
    sqlite3 \
    libsqlite3-dev \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

ENV GOLANG_VERSION=1.23.8
RUN wget https://go.dev/dl/go${GOLANG_VERSION}.linux-amd64.tar.gz && \
    tar -C /usr/local -xzf go${GOLANG_VERSION}.linux-amd64.tar.gz && \
    rm go${GOLANG_VERSION}.linux-amd64.tar.gz

ENV PATH="/usr/local/go/bin:$PATH"
ENV CGO_ENABLED=1
ENV GOOS=linux
ENV GOARCH=amd64

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Use the variable in the build
RUN go build -ldflags "-X main.version=${RELEASE_VERSION} -X main.buildTime=$(date -u +%Y-%m-%dT%H:%M:%SZ)" -o bin/${APP_NAME} main.go

# Runtime
FROM debian:bullseye-slim
ARG APP_NAME
ENV APP_NAME=${APP_NAME}

RUN apt-get update && apt-get install -y tzdata && rm -rf /var/lib/apt/lists/*

ENV TZ=Asia/Calcutta
ENV GIN_MODE=release

WORKDIR /app

COPY --from=builder /src/bin/${APP_NAME} .

CMD ["sh", "-c", "./$APP_NAME"]
