FROM golang:1.17 as builder
WORKDIR /app
COPY . .
RUN rm -fR minecraft-overseer && go get .
RUN go build -o overseer main.go
RUN ./overseer --help

FROM openjdk
COPY --from=builder /app/overseer /overseer
COPY log4j.xml /log4j.xml
CMD ["/overseer", "server"]