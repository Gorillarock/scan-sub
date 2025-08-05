FROM golang:1.22-alpine AS build

RUN apk add --no-cache gcc musl-dev

WORKDIR /app
COPY . .

ENV CGO_ENABLED=1
RUN go build -o test_db_reader ./cmd/test_db_reader

FROM alpine:latest

RUN apk add --no-cache ca-certificates sqlite

WORKDIR /app
COPY --from=build /app/test_db_reader /app/test_db_reader

ENTRYPOINT ["/app/test_db_reader"]