FROM golang:1.24-bookworm AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o main main.go

FROM debian:bookworm-slim AS prod

WORKDIR /app

COPY --from=build /app/main .

EXPOSE ${PORT}

ENTRYPOINT ["/app/main"]

