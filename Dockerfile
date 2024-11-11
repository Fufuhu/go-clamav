FROM golang:1.22 AS builder
COPY ./ /go/src
WORKDIR /go/src
RUN go mod download && go build -o /bin/go-clamav .

FROM debian:bookworm-slim
RUN apt-get update && apt-get install -y ca-certificates
COPY --from=builder /bin/go-clamav /bin/go-clamav
ENTRYPOINT ["/bin/go-clamav"]
CMD ["poll"]