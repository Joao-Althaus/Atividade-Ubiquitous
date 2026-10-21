FROM golang:1.23-alpine AS build
WORKDIR /src
COPY . .
RUN go mod tidy
RUN CGO_ENABLED=0 go build -o /app/server .

FROM alpine:3.21
RUN adduser -D app
USER app
COPY --from=build /app/server /server
EXPOSE 8080
CMD ["/server"]
