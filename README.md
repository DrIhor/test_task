# test_task

Envirinment variables:

1. STORAGE=['in-memory', 'postgres'] - type storage to use
1. SERVER_PORT - Can be passed to use localhost
1. SERVER_HOST
1. POSTGRE_HOST
1. POSTGRE_PORT
1. POSTGRE_USER
1. POSTGRE_PASS
1. POSTGRE_DB


# update protobuff

protoc --go_out=. --go_opt=paths=source_relative  --go-grpc_out=. --go-grpc_opt=paths=source_relative *.proto