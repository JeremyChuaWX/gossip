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
COPY ./go.mod ./go.sum ./

RUN go mod download

RUN go build -o ./server ./cmd/server/main.go

FROM gcr.io/distroless/base-debian12 AS runner

WORKDIR /app

COPY --from=tailwind /app/static/css/output.css ./static/css/
COPY --from=builder /app/server ./
COPY ./static ./static
COPY ./pages ./pages

EXPOSE 3000

CMD ["./server"]
