FROM golang:1.22-alpine AS build

WORKDIR /src

COPY go.mod ./
COPY cmd ./cmd
COPY internal ./internal

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/go-cli-app ./cmd/go-language-app

FROM alpine:3.20

RUN addgroup -S app && adduser -S app -G app

WORKDIR /app
COPY --from=build /app/go-cli-app ./go-cli-app

USER app
EXPOSE 8080

ENTRYPOINT ["./go-cli-app"]
