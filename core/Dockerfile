# builder image
FROM golang:1.16-alpine as base
WORKDIR /app
COPY . ./
RUN go get -d -v
RUN go get github.com/githubnemo/CompileDaemon

FROM base as builder
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64  go build -a -installsuffix cgo -o server .

# generate clean, final image for end users
FROM scratch
COPY --from=builder /app/server /app/server

ENTRYPOINT [ "/app/server" ]
