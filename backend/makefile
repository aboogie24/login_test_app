GIT_HASH ?= $(shell git log --format="%h" -n 1)


build: 
	docker build -t finance_app_backend:latest .
	docker build -t alfred2424/finance_app_backend:${GIT_HASH} .

login: 
	docker login -u alfred2424 -p fatboy24

push: login build
	docker push alfred2424/finance_app_backend:${GIT_HASH} 

release: push
	docker pull alfred2424/finance_app_backend:${GIT_HASH} 
	docker tag alfred2424/finance_app_backend:${GIT_HASH} alfred2424/finance_app_backend:latest
	docker push alfred2424/finance_app_backend:latest

run: build 
	docker-compose up 

local: 
	cd api && pwd && go run main.go 
	