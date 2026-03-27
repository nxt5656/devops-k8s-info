FROM golang:1.26-alpine AS builder

WORKDIR /src
ARG TARGETOS
ARG TARGETARCH

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=${TARGETOS:-linux} GOARCH=${TARGETARCH:-amd64} go build -o /out/k8s-info ./cmd/server

FROM gcr.io/distroless/base-debian12:nonroot

WORKDIR /app

COPY --from=builder /out/k8s-info /app/k8s-info

ENV SERVER_ADDR=:8080
ENV API_BASE_SEGMENT=k8s-info

EXPOSE 8080

ENTRYPOINT ["/app/k8s-info"]
