FROM golang:1.14.6-alpine3.12 as builder
COPY go.mod go.sum /go/src/github.com/peterkrauz/go-rest-api-playground/
WORKDIR /go/src/github.com/peterkrauz/go-rest-api-playground
RUN go mod download
COPY . /go/src/github.com/peterkrauz/go-rest-api-playground
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o build/go-rest-api-playground github.com/peterkrauz/go-rest-api-playground

FROM alpine
RUN apk add --no-cache ca-certificates && update-ca-certificates
COPY --from=builder /go/src/github.com/peterkrauz/go-rest-api-playground/build/go-rest-api-playground /usr/bin/go-rest-api-playground
EXPOSE 8080 8080
ENTRYPOINT ["/usr/bin/go-rest-api-playground"]