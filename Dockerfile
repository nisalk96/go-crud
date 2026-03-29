# syntax=docker/dockerfile:1

FROM golang:1.22-bookworm AS build

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/server ./cmd/server

FROM gcr.io/distroless/static-debian12:nonroot

WORKDIR /

COPY --from=build /out/server /server

USER nonroot:nonroot

EXPOSE 8080

ENV HTTP_ADDR=:8080
ENV UPLOAD_DIR=/data/covers

ENTRYPOINT ["/server"]
