###############################################################
# BUILD STAGE
###############################################################
FROM golang:1.18.0-alpine3.15 as builder
LABEL stage=builder
RUN apk update && apk add git
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o main ./cmd/api

###############################################################
# DISTRIBUTION STAGE
###############################################################
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /app/
COPY --from=builder /app .

ENV GIN_MODE=release
ENV PORT 6000

EXPOSE 6000
CMD ["./main"]
