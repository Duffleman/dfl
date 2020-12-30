FROM golang:1.15.6-alpine as builder
ENV CGO_ENABLED=0
RUN apk add --no-cache git
RUN mkdir -p /usr/local/app
WORKDIR /usr/local/app
COPY ./go.mod .
COPY ./go.sum .
RUN go mod download
COPY . /usr/local/app
RUN go test -v ./...
RUN go install -v ./cmd/...

FROM alpine
RUN apk add --no-cache ca-certificates tzdata
RUN mkdir -p /usr/local/app
WORKDIR /usr/local/app
ENTRYPOINT ["./svc"]
COPY --from=builder /usr/local/app/resources/auth.html resources/
COPY --from=builder /usr/local/app/resources/auth_code.html resources/
COPY --from=builder /usr/local/app/resources/markdown.html resources/
COPY --from=builder /usr/local/app/resources/not_found.html resources/
COPY --from=builder /usr/local/app/resources/nsfw.html resources/
COPY --from=builder /usr/local/app/resources/robots.txt resources/
COPY --from=builder /usr/local/app/resources/syntax_highlight.html resources/
COPY --from=builder /go/bin/svc .
