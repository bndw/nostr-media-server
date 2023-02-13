FROM golang:1.20-alpine as builder
RUN apk --no-cache add git

WORKDIR /go/src/github.com/bndw/nostr-media-server
COPY go.* ./
RUN go mod download

COPY . .
RUN go build -o /bin/nostr-media-server .

# --- Execution Stage

FROM alpine:latest

COPY --from=builder /bin/nostr-media-server /bin/

EXPOSE 80
ENTRYPOINT ["/bin/nostr-media-server"]
