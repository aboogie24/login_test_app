GIT_HASH ?= $(shell git log --format="%h" -n 1)


build: 
	docker build -t login_test_backend:latest ./backend/.
	docker build -t login_test_frontend:latest ./frontend/.

run: build 
	docker-compose up 

local: 
	cd api && pwd && go run main.go 