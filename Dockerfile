FROM golang:1.24-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /todo-api .

FROM alpine:3.20
RUN apk add --no-cache bash
ADD https://github.com/vishnubob/wait-for-it/raw/master/wait-for-it.sh /wait-for-it.sh
RUN chmod +x /wait-for-it.sh
WORKDIR /app
COPY --from=builder /todo-api .
COPY .env .
EXPOSE 3000
CMD ["/wait-for-it.sh", "db:5432", "--", "./todo-api"]