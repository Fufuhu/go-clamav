FROM golang:1.22 AS builder
COPY ./ ./src
WORKDIR /go/src
RUN go build -o /bin/go-clamav .

FROM debian:bookworm-slim
COPY --from=builder /bin/go-clamav /bin/go-clamav
ENTRYPOINT ["/bin/go-clamav"]
CMD ["poll"]