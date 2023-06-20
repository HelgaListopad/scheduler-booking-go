FROM golang:1.18 AS builder
WORKDIR /src
COPY . .
RUN go mod download
RUN CGO_ENABLED=1 GOOS=linux go build -o /backend -a -ldflags '-linkmode external -extldflags "-static"' .

FROM scratch
WORKDIR /app
COPY --from=builder /backend /app/backend
COPY ./demodata /app/demodata

EXPOSE 3000

CMD ["/app/backend"]
