FROM golang:1.18-rc-bullseye

WORKDIR /app
COPY . .
RUN go install

CMD ["go", "run", "main.go"]