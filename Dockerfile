# Build executable binary
FROM golang:alpine3.13 AS builder

ENV USER=authg
ENV UID=1000

RUN adduser \
  --disabled-password \
  --gecos "" \
  --home "$(pwd)" \
  --no-create-home \
  --uid "1000" \
  "authg"

RUN apk --update add ca-certificates

WORKDIR /build

COPY go.mod .
COPY go.sum .

RUN go mod download
RUN go mod verify

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

COPY . .

RUN go build -o auth-guardian .

WORKDIR /dist

RUN cp /build/auth-guardian .

# Build image
FROM scratch

LABEL Name=auth-guardian \
      Release=https://github.com/StevenCyb/Auth-Guardian \
      Url=https://github.com/StevenCyb/Auth-Guardian \
      Help=https://github.com/StevenCyb/Auth-Guardian/issues

COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

COPY --from=builder /dist/auth-guardian /

USER authg:authg
ENTRYPOINT ["/auth-guardian"]