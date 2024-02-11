FROM golang:alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY ./cmd ./internal ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /server ./cmd/server/main.go

EXPOSE 3000

CMD ["/server"]
