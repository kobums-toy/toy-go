tag=latest

all: server

server: dummy
	buildtool-model ./ 
	buildtool-router ./ > ./router/router.go
	go build -o bin/toysgo main.go

fswatch:
	fswatch -0 controllers | xargs -0 -n1 build/notify.sh

run:
	gin --port 9000 -a 9003 --bin bin/toysgo run main.go

allrun:
	fswatch -0 controllers | xargs -0 -n1 build/notify.sh &
	# gin --port 9000 -a 9003 --bin bin/toysgo run main.go
	go run .

test: dummy
	go test -v ./...

linux:
	env GOOS=linux GOARCH=amd64 go build -o bin/toysgo.linux main.go

dockerbuild:
	env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -ldflags '-s' -o bin/toysgo.linux main.go

docker: dockerbuild
	docker build -t kobums/toysgo:$(tag) .

dockerrun:
	docker run -d --name="toysgo" -p 9003:9003 kobums/toysgo

push: docker
	docker push kobums/toysgo:$(tag)

clean:
	rm -f bin/toysgo

dummy: