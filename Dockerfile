FROM golang:alpine as builder
RUN mkdir /build 
ADD . /build/
WORKDIR /build 
RUN go get ./...
RUN go build -o recipe-cms
RUN go install github.com/JamesTiberiusKirk/migrator/cmd/migrator@master

FROM alpine
COPY --from=builder /go/bin/migrator /app/
COPY --from=builder /build/site/public/ /app/site/public/
COPY --from=builder /build/recipe-cms /app/
COPY --from=builder /build/sql/ /app/sql/
WORKDIR /app

EXPOSE 80

ENTRYPOINT ["./recipe-cms"]
