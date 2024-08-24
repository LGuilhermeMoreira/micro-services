FROM golang:latest AS builder
RUN mkdir /app
COPY  . /app
WORKDIR /app
RUN go mod download
RUN CGO_ENABLED=0 go build -o authenticationApp ./cmd/api
RUN chmod +x /app/authenticationApp

#build a small image
FROM alpine:latest
RUN mkdir /app
COPY --from=builder /app/authenticationApp /app
CMD [ "/app/authenticationApp" ]