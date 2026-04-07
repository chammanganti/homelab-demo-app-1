FROM golang:1.25-alpine AS builder
WORKDIR /app

ARG TARGETARCH=arm64

COPY go.mod ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=${TARGETARCH} go build -ldflags="-w -s" -o homelab-demo-app-1 .

FROM gcr.io/distroless/static:nonroot
COPY --from=builder /app/homelab-demo-app-1 /homelab-demo-app-1
ENTRYPOINT ["/homelab-demo-app-1"]
