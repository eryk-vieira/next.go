build-cli:
	cd cli && go build -o ../nextgo ./main.go

build: build-cli
	./nextgo build

run:
	./nextgo run
