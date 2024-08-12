FROM  golang:latest as builder
RUN mkdir /app
COPY  . /app
WORKDIR /app
RUN go mod download
RUN CGO_ENABLED=0 go build -o brokerApp ./cmd/api
RUN chmod +x /app/brokerApp

#build a small image
FROM alpine:latest
RUN mkdir /app
COPY --from=builder /app/brokerApp /app
CMD [ "/app/brokerApp" ]