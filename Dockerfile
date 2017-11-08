FROM golang:latest
WORKDIR /go/src/gotrading
COPY . .

RUN go-wrapper download
RUN go-wrapper install

CMD ["go-wrapper", "run"] # ["app"]
