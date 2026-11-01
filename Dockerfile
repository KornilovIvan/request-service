FROM golang:1.26.6-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
ENV CGO_ENABLED=0
RUN go build -o /out/app ./cmd/app

FROM alpine:3.20
RUN apk add --no-cache ca-certificates
COPY --from=build /out/app /app
EXPOSE 8080
ENTRYPOINT ["/app"]
