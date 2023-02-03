FROM golang:1.20-alpine as builder
RUN apk --no-cache add git

WORKDIR /go/src/github.com/bndw/nostr-media-server
COPY go.* ./
RUN go mod download

COPY . .
RUN go build -o /bin/nostr-media-server .

# --- Execution Stage

FROM alpine:latest

WORKDIR /root/
COPY --from=builder /bin/nostr-media-server .

EXPOSE 80
CMD ["./nostr-media-server"]
