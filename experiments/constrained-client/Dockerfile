FROM golang:1.22-bookworm AS build

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -trimpath -o /out/rbpvss-bench ./cmd/rbpvss-bench
RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -trimpath -o /out/constrained-receiver ./cmd/constrained-receiver
RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -trimpath -o /out/full-comparison-bench ./cmd/full-comparison-bench

FROM debian:bookworm-slim
RUN apt-get update \
    && apt-get install --no-install-recommends -y ca-certificates curl iproute2 \
    && rm -rf /var/lib/apt/lists/*
COPY --from=build /out/rbpvss-bench /usr/local/bin/rbpvss-bench
COPY --from=build /out/constrained-receiver /usr/local/bin/constrained-receiver
COPY --from=build /out/full-comparison-bench /usr/local/bin/full-comparison-bench
ENTRYPOINT ["rbpvss-bench"]
