FROM alpine:latest

RUN apk update && apk add sqlite

WORKDIR /db

COPY schema.sql .

# Initialize the database with the schema
RUN sqlite3 scans.db < schema.sql

CMD ["tail", "-f", "/dev/null"]