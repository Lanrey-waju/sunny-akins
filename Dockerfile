FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY . .

RUN go mod download
RUN go mod verify

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /server ./cmd/web/.

FROM gcr.io/distroless/static-debian11

COPY --from=builder /server .

EXPOSE 8080

CMD ["/server"]
