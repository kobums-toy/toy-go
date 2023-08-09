tag=latest

all: server

server: dummy
	buildtool-model ./ 
	buildtool-router ./ > ./router/router.go
	go build -o bin/gofiber main.go

fswatch:
	fswatch -0 controllers | xargs -0 -n1 build/notify.sh

run:
	gin --port 9000 -a 9003 --bin bin/gofiber run main.go

allrun:
	fswatch -0 controllers | xargs -0 -n1 build/notify.sh &
	gin --port 9000 -a 9003 --bin bin/gofiber run main.go

test: dummy
	go test -v ./...

linux:
	env GOOS=linux GOARCH=amd64 go build -o bin/gofiber.linux main.go

dockerbuild:
	env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -ldflags '-s' -o bin/gofiber.linux main.go

docker: dockerbuild
	docker build -t kobums/gofiber:$(tag) .

dockerrun:
	docker run -d --name="gofiber" -p 9003:9003 kobums/gofiber

push: docker
	docker push kobums/gofiber:$(tag)

clean:
	rm -f bin/gofiber

dummy:
