FROM golang:1.23-bookworm

# ---- basic OS deps ----
RUN apt-get update && apt-get install -y \
    git \
    curl \
    ca-certificates \
    build-essential \
    pkg-config \
    lldb \
    && rm -rf /var/lib/apt/lists/*

# ---- Go environment ----
ENV GO111MODULE=on \
    GOPATH=/go \
    GOCACHE=/go/cache \
    CGO_ENABLED=1

# ---- tools (explicit, versionable) ----
RUN go install golang.org/x/tools/gopls@latest \
    && go install github.com/go-delve/delve/cmd/dlv@latest \
    && go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# ---- workdir ----
WORKDIR /workspace

# ---- non-root user (important for VS Code + volumes) ----
ARG USERNAME=vscode
ARG USER_UID=1000
ARG USER_GID=1000

RUN groupadd --gid $USER_GID $USERNAME \
    && useradd --uid $USER_UID --gid $USER_GID -m $USERNAME \
    && chown -R $USERNAME:$USERNAME /go /workspace

USER $USERNAME
