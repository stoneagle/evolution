.PHONY: run-web, stop-web rm-web

PWD := $(shell pwd)
USER := $(shell id -u)
GROUP := $(shell id -g)

run-local: 
	cd hack && docker-compose -f docker-compose-local.yml -p "airflow-$(USER)" up -d

stop-local: 
	cd hack && docker-compose -f docker-compose-local.yml -p "airflow-$(USER)" stop 

rm-local: 
	cd hack && docker-compose -f docker-compose-local.yml -p "airflow-$(USER)" rm 

build-img:
	cd hack/dockerfile && docker build -f ./Dockerfile-assets -t puckel/docker-airflow:1.8.2-assets .
