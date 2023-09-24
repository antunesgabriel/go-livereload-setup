FROM golang:1.21.1-alpine3.18 as buildstage

WORKDIR /go/src

RUN apk update

COPY . .

RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o api ./cmd/api

FROM busybox

ENV PORT=8080

COPY --from=buildstage /go/src/api /app/

EXPOSE $PORT

CMD [ "/app/api" ]
