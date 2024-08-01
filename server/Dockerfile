FROM node:alpine AS tailwind

WORKDIR /app

COPY ./static ./static
COPY ./pages ./pages
COPY ./package.json ./package-lock.json ./tailwind.config.js ./

RUN npm ci

RUN npm run build

FROM golang:alpine AS builder

WORKDIR /app

COPY ./cmd ./cmd
COPY ./internal ./internal
COPY ./static ./static
COPY ./pages ./pages
COPY ./go.mod ./go.sum ./

COPY --from=tailwind /app/static/css/output.css ./static/css/

RUN go mod download

RUN GOOS=linux go build -o ./server ./cmd/server/main.go

EXPOSE 3000

CMD ["./server"]
