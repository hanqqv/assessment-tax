FROM golang:1.22.2-alpine as build-base

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go test -tags=unit -v ./...

RUN go build -o ./out/assessment-tax

### -----------

FROM alpine:3.19
COPY --from=build-base /app/out/assessment-tax /app/assessment-tax

EXPOSE 8080

CMD ["/app/assessment-tax"]