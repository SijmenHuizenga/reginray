# redeploy: make disable remove build set enable -k

build:
	CGO_ENABLED=0 GOOS=linux go build -a -o aspicio
	docker build -t rootfsimage .
	docker create --name buildimage rootfsimage true
	mkdir -p .out/rootfs
	cp config.json .out
	docker export buildimage | tar -x -C .out/rootfs
	docker rm -vf buildimage
	docker rmi rootfsimage
	docker plugin create quay.io/aspicio-docker-plugin:snapshot .out
	rm -rf .out
	rm aspicio

enable:
	docker plugin enable quay.io/aspicio-docker-plugin:snapshot

disable:
	docker plugin disable quay.io/aspicio-docker-plugin:snapshot

remove:
	docker plugin remove quay.io/aspicio-docker-plugin:snapshot

set:
	docker plugin set quay.io/aspicio-docker-plugin:snapshot ASPICIO_BACKEND=http://192.168.0.25:8080