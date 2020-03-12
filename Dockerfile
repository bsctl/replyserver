FROM golang:1.13 as builder
WORKDIR /code
ADD go.* /code/
RUN go mod download
ADD . .
RUN go build -o /replyserver main.go

FROM gcr.io/distroless/base
ENV VERSION "v0.0.2"
EXPOSE 1968 1969
EXPOSE 1936
WORKDIR /
COPY --from=builder /replyserver /usr/bin/replyserver
ENTRYPOINT ["/usr/bin/replyserver"]