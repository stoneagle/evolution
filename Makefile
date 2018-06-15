.PHONY: run-web, stop-web rm-web

PWD := $(shell pwd)
USER := $(shell id -u)
GROUP := $(shell id -g)
DATE := $(shell date "+%F")
REGISTRY_PREFIX := "stoneagle/quant-engine"
PROJ := "quant"

# backend
run-web: 
	export REGISTRY_PREFIX=$(REGISTRY_PREFIX) && \
	cd hack/swarm && docker-compose -f docker-compose.yml -p "$(PROJ)-$(USER)-web" up

stop-web: 
	export REGISTRY_PREFIX=$(REGISTRY_PREFIX) && \
	cd hack/swarm && docker-compose -f docker-compose.yml -p "$(PROJ)-$(USER)-web" stop 

rm-web: 
	export REGISTRY_PREFIX=$(REGISTRY_PREFIX) && \
	cd hack/swarm && docker-compose -f docker-compose.yml -p "$(PROJ)-$(USER)-web" rm 

build-golang:
	cd hack/dockerfile && docker build -f ./Dockerfile-golang -t quant/golang:1.10 .

# init
init-db:
	docker exec -w /go/src/quant/backend/initial -it quant-wuzhongyang-golang go run init.go 

# frontend
run-ng:
	cd frontend && ng serve --environment=dev

# grafana 
run-grafana: 
	cd hack && docker-compose -f docker-compose-grafana.yml -p "grafana-$(USER)" up -d

stop-grafana: 
	cd hack && docker-compose -f docker-compose-grafana.yml -p "grafana-$(USER)" stop 

rm-grafana: 
	cd hack && docker-compose -f docker-compose-grafana.yml -p "grafana-$(USER)" rm 

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
		thrift --gen go -out $(THRIFT_PREFIX)/backend/rpc $(THRIFT_PREFIX)/hack/thrift/engine.thrift

thrift-python:
	docker run -it --rm \
		-u $(USER):$(GROUP) \
		-v $(PWD)/engine:$(THRIFT_PREFIX)/engine \
		-v $(PWD)/hack:$(THRIFT_PREFIX)/hack \
		thrift:0.11.0 \
		thrift --gen py -out $(THRIFT_PREFIX)/engine/rpc $(THRIFT_PREFIX)/hack/thrift/engine.thrift

# engine
build-engine-basic:
	cd ./hack/dockerfile && docker build -f ./Dockerfile.engine -t $(REGISTRY_PREFIX):quant . --network=host
