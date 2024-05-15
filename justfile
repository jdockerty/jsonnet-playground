

# Build with ko, does not push to the registry
build REGISTRY="ghcr.io/jdockerty/jsonnet-playground":
    KO_DOCKER_REPO={{ REGISTRY }} ko build github.com/jdockerty/jsonnet-playground/cmd/server --push=false --platform=linux/arm64,linux/amd64

build_local:
    ko build github.com/jdockerty/jsonnet-playground/cmd/server --local --push=false --platform=linux/arm64,linux/amd64

# Build and push to the registry
push REGISTRY="ghcr.io/jdockerty/jsonnet-playground":
    KO_DOCKER_REPO={{ REGISTRY }} ko build github.com/jdockerty/jsonnet-playground/cmd/server --platform=linux/arm64,linux/amd64

run:
    KO_DATA_PATH="cmd/server/kodata" go run cmd/server/cmd.go

# Run the server with hot reloading for templ components
run_reload:
    KO_DATA_PATH="cmd/server/kodata" templ generate --watch --proxy="http://127.0.0.1:8080" --cmd='go run cmd/server/cmd.go'

# Install required dependencies
deps:
    go install github.com/google/ko@latest

# Run various lint/generation tools
lint:
    go fmt ./...
    templ fmt -v .
    templ generate -v
