FROM golang:1.24-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/komiac-support-backend

FROM alpine:3.20

WORKDIR /app

COPY --from=build /app/server /app/server

EXPOSE 8080

CMD ["/app/server"]
