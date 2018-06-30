.PHONY: run-web, stop-web rm-web

PWD := $(shell pwd)
USERNAME := $(shell id -nu)
USER := $(shell id -u)
GROUP := $(shell id -g)
DATE := $(shell date "+%F")
REGISTRY_DEVELOP_PREFIX := stoneagle/develop
PROJ := quant

# backend
run-web: 
	cd hack/swarm && docker-compose -f docker-compose.yml -p "$(PROJ)-$(USER)-web" up

stop-web: 
	cd hack/swarm && docker-compose -f docker-compose.yml -p "$(PROJ)-$(USER)-web" stop 

rm-web: 
	cd hack/swarm && docker-compose -f docker-compose.yml -p "$(PROJ)-$(USER)-web" rm 

# init
init-db:
	docker exec -w /go/src/quant/backend/initial -it quant-wuzhongyang-golang go run init.go 

init-influxdb:
	sudo docker run --rm \
		-e INFLUXDB_DB=quant -e INFLUXDB_ADMIN_ENABLED=true \
		-e INFLUXDB_ADMIN_USER=admin -e INFLUXDB_ADMIN_PASSWORD=a1b2c3d4E \
		-e INFLUXDB_USER=quant -e INFLUXDB_USER_PASSWORD=a1b2c3d4E \
		-v /home/$(USERNAME)/database/influxdb:/var/lib/influxdb \
		influxdb:1.5.3 /init-influxdb.sh

# frontend
run-ng:
	cd frontend && ng serve --environment=dev

# grafana 
init-plugin:
	cd plugin/ashare && npm install --registry=http://rgistry.npm.taobao.org && ./node_modules/grunt/bin/grunt

build-grafana:
	cd hack/dockerfile && docker build -f ./Dockerfile-grafana -t grafana/grafana:4.6.2-1000 .

# thrift
THRIFT_PREFIX = /tmp/thrift
thrift-golang:
	docker run -it --rm \
		-u $(USER):$(GROUP) \
		-v $(PWD)/backend:$(THRIFT_PREFIX)/backend \
		-v $(PWD)/hack:$(THRIFT_PREFIX)/hack \
		thrift:0.11.0 \
		thrift --gen go -out $(THRIFT_PREFIX)/backend/rpc $(THRIFT_PREFIX)/hack/thrift/engine.thrift && \
	sed -i 's:"engine":"quant/backend/rpc/engine":' ./backend/rpc/engine/*/*.go

thrift-python:
	docker run -it --rm \
		-u $(USER):$(GROUP) \
		-v $(PWD)/engine:$(THRIFT_PREFIX)/engine \
		-v $(PWD)/hack:$(THRIFT_PREFIX)/hack \
		thrift:0.11.0 \
		thrift --gen py -out $(THRIFT_PREFIX)/engine/rpc $(THRIFT_PREFIX)/hack/thrift/engine.thrift

# engine
build-engine-basic:
	cd ./hack/dockerfile && docker build -f ./Dockerfile.engine -t $(REGISTRY_DEVELOP_PREFIX):quant-engine . --network=host

build-golang:
	cd hack/dockerfile && docker build -f ./Dockerfile-golang -t $(REGISTRY_DEVELOP_PREFIX):golang-beego-1.10 .

