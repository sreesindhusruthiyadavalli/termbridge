FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o termbridge main.go

FROM alpine:latest
RUN apk --no-cache add bash
WORKDIR /app
COPY --from=builder /app/termbridge .
COPY static ./static
EXPOSE 8080
CMD ["./termbridge", "-port=8080"]
