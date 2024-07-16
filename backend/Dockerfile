FROM golang:1.21

WORKDIR /app

COPY go.mod go.sum ./
RUN go install github.com/cosmtrek/air@v1.40.4
RUN go mod download

COPY . .

RUN go build -o main ./cmd

CMD ["make", "start"]