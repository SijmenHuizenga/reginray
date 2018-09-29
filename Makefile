all: clean build-frontend build-backend

build-frontend:
	cd frontend && npm run build
	mkdir -p build/static
	cp -rf frontend/build/* build/static/

build-backend:
	cd backend && make
	mkdir -p build
	mv backend/aspicio-backend build/app

clean:
	rm -rf build

run:
	cp config.ini build/config.ini
	cd build && PORT=8080 MONGO_SERVER=mongodb://localhost:27017 ./app

run-db:
	docker run --log-driver quay.io/aspicio-docker-plugin:snapshot -p 27017:27017 -v aspicio-db:/data/db mongo