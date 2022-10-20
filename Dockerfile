FROM golang:1.19 AS builder
WORKDIR /go/src/app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY main.go .
RUN CGO_ENABLED=0 go build -o app

FROM gcr.io/distroless/static-debian11 AS runtime
COPY --from=builder /go/src/app/app /usr/local/bin/app
ENTRYPOINT ["/usr/local/bin/app"]