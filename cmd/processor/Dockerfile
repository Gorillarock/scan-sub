FROM golang:1.22-alpine AS build

RUN apk add --no-cache gcc musl-dev

WORKDIR /app
COPY . .

ENV CGO_ENABLED=1
RUN go build -o processor ./cmd/processor

FROM alpine:latest

RUN apk add --no-cache ca-certificates sqlite

WORKDIR /app
COPY --from=build /app/processor /app/processor

ENTRYPOINT ["/app/processor"]