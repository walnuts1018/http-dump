FROM golang:1.21.0 AS builder
ENV ROOT=/build
RUN mkdir ${ROOT}
WORKDIR ${ROOT}

COPY ./ ./
RUN go get

RUN CGO_ENABLED=0 GOOS=linux go build -o main $ROOT/main.go && chmod +x ./main

FROM alpine:3.18
WORKDIR /app

COPY --from=builder /build/main ./
EXPOSE 8080

CMD ["./main"]
