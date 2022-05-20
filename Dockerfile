FROM --platform=$BUILDPLATFORM golang:1.18.2 as builder
ARG TARGETPLATFORM
ARG BUILDPLATFORM
ARG TARGETOS
ARG TARGETARCH
WORKDIR /app
COPY go.mod go.sum /app
RUN go mod download -x
COPY . .
RUN go test ./internal/...
RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build -ldflags='-w -s -extldflags "-static"' -o service ./cmd/server

FROM --platform=$TARGETPLATFORM openjdk
COPY --from=builder /app/service /overseer
COPY log4j.xml /log4j.xml
CMD ["/overseer", "run"]
