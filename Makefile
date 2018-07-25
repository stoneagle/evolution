.PHONY: run-web, stop-web rm-web

PWD := $(shell pwd)
USERNAME := $(shell id -nu)
USER := $(shell id -u)
GROUP := $(shell id -g)
DATE := $(shell date "+%F")
REGISTRY_DEVELOP_PREFIX := stoneagle/develop
PROJ := evolution

# quant 
run-quant: 
	cd hack/swarm && docker-compose -f docker-compose-quant.yml -p "$(PROJ)-$(USER)-quant" up
stop-quant: 
	cd hack/swarm && docker-compose -f docker-compose-quant.yml -p "$(PROJ)-$(USER)-quant" stop 
rm-quant: 
	cd hack/swarm && docker-compose -f docker-compose-quant.yml -p "$(PROJ)-$(USER)-quant" rm 
	
# time 
run-time: 
	cd hack/swarm && docker-compose -f docker-compose-time.yml -p "$(PROJ)-$(USER)-time" up
stop-time: 
	cd hack/swarm && docker-compose -f docker-compose-time.yml -p "$(PROJ)-$(USER)-time" stop
rm-time: 
	cd hack/swarm && docker-compose -f docker-compose-time.yml -p "$(PROJ)-$(USER)-time" rm

# system 
run-system: 
	cd hack/swarm && docker-compose -f docker-compose-system.yml -p "$(PROJ)-$(USER)-system" up
stop-system: 
	cd hack/swarm && docker-compose -f docker-compose-system.yml -p "$(PROJ)-$(USER)-system" stop
rm-system: 
	cd hack/swarm && docker-compose -f docker-compose-system.yml -p "$(PROJ)-$(USER)-system" rm

# envoy 
run-envoy: 
	cd hack/swarm && docker-compose -f docker-compose-envoy.yml -p "$(PROJ)-$(USER)-envoy" up -d
stop-envoy: 
	cd hack/swarm && docker-compose -f docker-compose-envoy.yml -p "$(PROJ)-$(USER)-envoy" stop 
rm-envoy: 
	cd hack/swarm && docker-compose -f docker-compose-envoy.yml -p "$(PROJ)-$(USER)-envoy" rm 

# init
init-quant-db:
	docker exec -w /go/src/evolution/backend/quant/initial -it evolution-wuzhongyang-quant-backend go run init.go 
init-time-db:
	docker exec -w /go/src/evolution/backend/time/initial -it evolution-wuzhongyang-time-backend go run init.go 
init-system-db:
	docker exec -w /go/src/evolution/backend/system/initial -it evolution-wuzhongyang-system-backend go run init.go 

# transfer
transfer-time-db:
	docker exec -w /go/src/evolution/backend/time/models/transfer -it evolution-wuzhongyang-time-backend go run transfer.go 

init-influxdb:
	sudo docker run --rm \
		-e INFLUXDB_DB=quant -e INFLUXDB_ADMIN_ENABLED=true \
		-e INFLUXDB_ADMIN_USER=admin -e INFLUXDB_ADMIN_PASSWORD=a1b2c3d4E \
		-e INFLUXDB_USER=quant -e INFLUXDB_USER_PASSWORD=a1b2c3d4E \
		-v /home/$(USERNAME)/database/influxdb:/var/lib/influxdb \
		influxdb:1.5.3 /init-influxdb.sh

# frontend
run-ng:
	cd frontend && ng serve --environment=dev --poll=2000

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
		-v $(PWD)/quant:$(THRIFT_PREFIX)/quant \
		-v $(PWD)/hack:$(THRIFT_PREFIX)/hack \
		thrift:0.11.0 \
		thrift --gen go -out $(THRIFT_PREFIX)/quant/rpc $(THRIFT_PREFIX)/hack/thrift/engine.thrift && \
	sed -i 's:"engine":"evolution/quant/rpc/engine":' ./quant/rpc/engine/*/*.go

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
