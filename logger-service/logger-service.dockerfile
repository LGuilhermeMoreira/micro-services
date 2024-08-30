FROM golang:latest AS builder
RUN mkdir /app
COPY  . /app
WORKDIR /app
RUN go mod download
RUN CGO_ENABLED=0 go build -o loggerApp ./cmd/api
RUN chmod +x /app/loggerApp

#build a small image
FROM alpine:latest
RUN mkdir /app
COPY --from=builder /app/loggerApp /app
CMD [ "/app/loggerApp" ]