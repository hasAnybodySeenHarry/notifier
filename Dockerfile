FROM golang:1.22-alpine AS builder

LABEL maintainer="harryd.io@proton.me"

ENV CGO_ENABLED=0

RUN addgroup -S myuser && adduser -S -D -G myuser myuser

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o /bin/app ./cmd/api

FROM scratch

COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group

WORKDIR /app

COPY --from=builder /bin/app /app/app

USER myuser:myuser

EXPOSE 8080

CMD ["./app"]