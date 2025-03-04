FROM golang:1.24.1 AS builder
ENV ROOT=/build
RUN mkdir ${ROOT}
WORKDIR ${ROOT}

RUN --mount=type=cache,target=/go/pkg/mod/,sharing=locked \
    --mount=type=cache,target=/root/.cache/go-build,sharing=locked \
    --mount=type=bind,source=go.sum,target=go.sum \
    --mount=type=bind,source=go.mod,target=go.mod \
    go mod download -x

COPY . .
RUN --mount=type=cache,target=/go/pkg/mod/,sharing=locked \
    --mount=type=cache,target=/root/.cache/go-build,sharing=locked \
    CGO_ENABLED=0 GOOS=linux go build -o main $ROOT && chmod +x ./main

FROM debian:bookworm-slim
WORKDIR /app

COPY --from=builder /build/main ./
EXPOSE 8080
USER 1000

CMD ["./main"]
