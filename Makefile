DEFAULT_OUT := "chat-cli"

all: get pb build
gen: pb
	protoc --go_out=pb --go_opt=paths=source_relative --go-grpc_out=pb --go-grpc_opt=paths=source_relative --proto_path=protos protos/*.proto
pb:
	mkdir pb
clean:
	rm -rf pb/*
get:
	go mod download && go mod verify
build: get pb
ifeq ($(strip $(OUT)),)
	# use default executable output if not previously defined
	echo "OUT not defined. Using default build output" 
	go build -ldflags="-X main.SERVER_ADDR=${SERVER_ADDR}" -o $(DEFAULT_OUT) main.go
else
	go build -ldflags="-X main.SERVER_ADDR=${SERVER_ADDR}" -o $(OUT) main.go
endif
