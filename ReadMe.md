docker compose up -d
docker compose down -v

docker build . -t goprotobuff:v1
docker run --rm -v ${PWD}/generated:/go/generated -v ${PWD}/protofiles:/go/input goprotobuff:v1

go run cmd/main.go