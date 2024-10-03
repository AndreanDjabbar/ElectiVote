FROM golang:1.23-alpine

WORKDIR /app

RUN go install github.com/air-verse/air@latest

COPY go.mod go.sum ./

RUN mkdir -p logs

RUN go mod download

COPY . .

RUN go build -o main ./cmd/ElectiVote

EXPOSE 8080

CMD ["air", "-c", ".air.toml"]