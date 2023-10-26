# FROM "extiende" otra imagen
FROM golang:latest AS build

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . ./

RUN CGO_ENABLED=0 GOOS=linux go build  -o main ./cmd/tuiter/main.go

FROM public.ecr.aws/docker/library/alpine:latest

COPY --from=build /app/main .

EXPOSE 3000

CMD ["./main"]