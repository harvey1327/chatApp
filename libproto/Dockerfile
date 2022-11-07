FROM golang:1.19.1-alpine3.16

ENV VERSION "21.7"
RUN apk update && apk --no-cache add curl protoc protobuf git
RUN curl -LO https://github.com/protocolbuffers/protobuf/releases/download/v$VERSION/protoc-$VERSION-linux-x86_64.zip && unzip protoc-$VERSION-linux-x86_64.zip -d /usr/bin && rm protoc-21.7-linux-x86_64.zip
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

COPY build_proto.sh .

CMD ["./build_proto.sh"]