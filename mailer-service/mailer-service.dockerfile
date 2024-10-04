FROM golang:latest AS builder
RUN mkdir /app
COPY  . /app
WORKDIR /app
RUN go mod download
RUN CGO_ENABLED=0 go build -o mailerApp ./cmd/api
RUN chmod +x /app/mailerApp

#build a small image
FROM alpine:latest
RUN mkdir /app
COPY --from=builder /app/mailerApp /app
COPY --from=builder /app/templates /templates
CMD [ "/app/mailerApp" ]