FROM golang:1.24 AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/gbw


FROM gcr.io/distroless/static-debian13

COPY --from=build /app/bin/gbw /
EXPOSE 8080

CMD ["/gbw"]
