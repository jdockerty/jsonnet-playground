FROM golang:1.22 as build

WORKDIR /go/src/app
COPY . .

RUN go mod download
RUN go vet -v ./...

RUN CGO_ENABLED=0 go build -o /go/bin/playground cmd/server/*

FROM gcr.io/distroless/static-debian12

COPY --from=build /go/bin/playground /

EXPOSE 8080
CMD ["/playground"]
