grcp:
    go run ./cmd/grpcServer/ -envVars

server:
    go run ./cmd/httpAndGraphQl/ -envVars

docker-compose:
    docker-compose up --build -d

test:
    go test ./...