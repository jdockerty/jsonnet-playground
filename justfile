

# Build with ko, does not push to the registry
build REGISTRY="ghcr.io/jdockerty/jsonnet-playground":
    KO_DOCKER_REPO={{ REGISTRY }} KO_DATA_PATH="assets" ko build ./cmd/server --push=false --platform=linux/arm64,linux/amd64

# Build and push to the registry
push REGISTRY="ghcr.io/jdockerty/jsonnet-playground":
    KO_DOCKER_REPO={{ REGISTRY }} KO_DATA_PATH="assets" ko build ./cmd/server --platform=linux/arm64,linux/amd64

deps:
    go install github.com/google/ko@latest
