FROM golang:1.17.8-alpine as dev

WORKDIR /app
RUN apk update && apk add git
COPY . .
RUN go mod download

CMD ["go", "run", "main.go"]


FROM golang:1.17.8-alpine as builder

WORKDIR /app
RUN apk update && apk add git
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/main


FROM scratch as prod

WORKDIR /app
COPY --from=builder /app/main .

CMD ["./main"]
