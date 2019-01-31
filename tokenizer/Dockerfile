FROM golang:1.11-alpine3.8

RUN apk update && apk add git

ENV GOPATH /go

RUN go get -v github.com/Lupino/tokenizer/tokenizer

FROM alpine:3.8

COPY --from=0 /go/bin/tokenizer /usr/bin/tokenizer
COPY --from=0 /go/src/github.com/huichen/sego/data/dictionary.txt /dictionary.txt

ENTRYPOINT ["tokenizer"]

CMD ["-host", "0.0.0.0:3000", "-dicts", "/dictionary.txt"]
