FROM golang:1.27.3@sha256:e519d3b0a5b285e5aa53ab9eef02bc6bbb9cb04608389b743a136a9ce75f0d1e AS builder
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
