FROM --platform=$TARGETPLATFORM golang:1.17 as builder
ARG TARGETPLATFORM
ARG BUILDPLATFORM
ARG TARGETARCH
RUN uname -a
RUN echo $BUILDPLATFORM $TARGETPLATFORM
WORKDIR /app
COPY . .
RUN GOARCH=$TARGETARCH go build -o overseer ./cmd/server

FROM --platform=$TARGETPLATFORM openjdk
COPY --from=builder /app/overseer /overseer
COPY log4j.xml /log4j.xml
CMD ["/overseer", "run"]
