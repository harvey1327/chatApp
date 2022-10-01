#!/usr/bin/env sh

# Generates the protobuff files
ROOT=/input/user
PROTO_DEST=/generated
MODULE=github.com/chatapp/proto

# Uses protoc to generate the files
printf "protoc -I=/go%s --go_out=/go%s --go_opt=module=%s%s --go-grpc_out=/go%s --go-grpc_opt=module=%s%s /go%s/*.proto\n" "${ROOT}" "${PROTO_DEST}" "${MODULE}" "${PROTO_DEST}" "${PROTO_DEST}" "${MODULE}" "${PROTO_DEST}" "${ROOT}"
protoc -I=/go${ROOT} \
--go_out=/go${PROTO_DEST} \
--go_opt=module=${MODULE}${PROTO_DEST} \
--go-grpc_out=/go${PROTO_DEST} \
--go-grpc_opt=module=${MODULE}${PROTO_DEST} \
/go${ROOT}/*.proto