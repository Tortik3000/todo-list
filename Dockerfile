FROM golang:1.25 as builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o todo-list ./cmd/todo-list

FROM gcr.io/distroless/base-debian12 AS service
COPY --from=builder /app/todo-list /todo-list
ENTRYPOINT ["/todo-list"]

