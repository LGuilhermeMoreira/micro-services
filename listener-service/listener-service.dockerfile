FROM golang:latest AS builder
RUN mkdir /app
COPY . /app
WORKDIR /app
RUN go mod download
RUN CGO_ENABLED=0 go build -o listenerApp ./cmd/listen
RUN chmod +x /app/listenerApp

#build a small image
FROM alpine:latest
RUN mkdir /app
COPY --from=builder /app/listenerApp /app
CMD [ "/app/listenerApp" ]
