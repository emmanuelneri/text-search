FROM golang:1.15-alpine as builder
RUN apk update && apk add --no-cache git
WORKDIR build
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /bin/app ./cmd

FROM scratch
COPY --from=builder /bin/app /bin/app
ENV SERVER_PORT "8080"
ENV ELASTIC_URL "http://localhost:9200"
ENV ELASTIC_INDEX "users"
ENV ELASTIC_DEBUG "false"
ENTRYPOINT ["/bin/app"]
EXPOSE 8080