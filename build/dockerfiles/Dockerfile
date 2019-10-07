FROM golang:1.13.1-alpine3.10 AS builder

WORKDIR /go/src/github.com/yu-ame/mixiapps-restful
COPY . . 
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o /main ./cmd/mixiapps-restful

# final stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=builder /main /mixiapps-restful/main
COPY --from=builder /go/src/github.com/yu-ame/mixiapps-restful/web /mixiapps-restful/web
COPY --from=builder /go/src/github.com/yu-ame/mixiapps-restful/configs /mixiapps-restful/configs
RUN chmod +x /mixiapps-restful/main
WORKDIR /mixiapps-restful
CMD ["/mixiapps-restful/main"]
