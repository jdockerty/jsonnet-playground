

# Build with ko, does not push to the registry
build REGISTRY="ghcr.io/jdockerty/jsonnet-playground":
    KO_DOCKER_REPO={{ REGISTRY }} ko build ./cmd/server --push=false

# Build and push to the registry
push REGISTRY="ghcr.io/jdockerty/jsonnet-playground":
    KO_DOCKER_REPO={{ REGISTRY }} ko build ./cmd/server

deps:
    go install github.com/google/ko@latest
