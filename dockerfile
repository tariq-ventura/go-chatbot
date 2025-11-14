FROM golang:1.24.5 AS bin-stage

SHELL ["/bin/bash", "-c"]

RUN mkdir -p /github.com/tariq-ventura/go-chatbot

WORKDIR /github.com/tariq-ventura/go-chatbot

COPY . .

RUN go mod download && go mod tidy

RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/api ./cmd/server/main.go

FROM alpine:latest AS release-stage

WORKDIR /

COPY --from=bin-stage /bin/api /bin/api
VOLUME /var/log

EXPOSE 3000

ENTRYPOINT ["/bin/api"]