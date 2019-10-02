.PHONY: run-web, stop-web rm-web

PWD := $(shell pwd)
DIRNAME := $(shell dirname "$(PWD)")
USERNAME := $(shell id -nu)
USER := $(shell id -u)
GROUP := $(shell id -g)
PROJ := evolution
PREFIX_THRIFT := /tmp/thrift
PREFIX_REGISTRY := stoneagle
PREFIX_DEVELOP := $(PREFIX_REGISTRY)/develop
RELEASE_DIR = release
SYS_TIME = time
GOVERSION = 1.10

# SYSTEM can be quant/time/system/basic
run:
	cd hack/swarm && docker-compose -f docker-compose-$(SYSTEM).yml -p "$(PROJ)-$(USER)-$(SYSTEM)" up $(PARAMS)
stop:
	cd hack/swarm && docker-compose -f docker-compose-$(SYSTEM).yml -p "$(PROJ)-$(USER)-$(SYSTEM)" stop 
rm:
	cd hack/swarm && docker-compose -f docker-compose-$(SYSTEM).yml -p "$(PROJ)-$(USER)-$(SYSTEM)" rm 

test-run:
	cd hack/release/swarm && docker-compose up
test-run-db:
	cd hack/release/swarm && docker-compose -f docker-compose-mysql.yaml -p "evolution-mysql" up
test-stop:
	cd hack/release/swarm && docker-compose stop 
test-rm:
	cd hack/release/swarm && docker-compose rm

# SYSTEM can be quant/time/system
init-db:
	docker exec -w /go/src/evolution/backend/$(SYSTEM)/initial -it $(PROJ)-$(USERNAME)-$(SYSTEM)-backend go run init.go 

# transfer
transfer-time-db:
	docker exec -w /go/src/evolution/backend/time/models/transfer -it $(PROJ)-$(USERNAME)-$(SYS_TIME)-backend go run transfer.go 

init-influxdb:
	sudo docker run --rm \
		-e INFLUXDB_DB=quant -e INFLUXDB_ADMIN_ENABLED=true \
		-e INFLUXDB_ADMIN_USER=admin -e INFLUXDB_ADMIN_PASSWORD=a1b2c3d4E \
		-e INFLUXDB_USER=quant -e INFLUXDB_USER_PASSWORD=a1b2c3d4E \
		-v /home/$(USERNAME)/database/influxdb:/var/lib/influxdb \
		influxdb:1.5.3 /init-influxdb.sh

# build
check-release:
ifeq "$(wildcard $(RELEASE_DIR))" ""
	echo "release directory not existed, make it"
	mkdir $(RELEASE_DIR)
else
	echo "release directory existed"
endif

# frontend
run-frontend:
	cd frontend && ng serve --environment=dev --poll=2000 --disable-host-check 
# grafana 
init-plugin:
	cd plugin/ashare && npm install --registry=http://rgistry.npm.taobao.org && ./node_modules/grunt/bin/grunt

build-grafana:
	cd hack/dockerfile && docker build -f ./Dockerfile-grafana -t grafana/grafana:4.6.2-1000 .

# thrift
thrift-golang:
	docker run -it --rm \
		-u $(USER):$(GROUP) \
		-v $(PWD)/quant:$(PREFIX_THRIFT)/quant \
		-v $(PWD)/hack:$(PREFIX_THRIFT)/hack \
		thrift:0.11.0 \
		thrift --gen go -out $(PREFIX_THRIFT)/quant/rpc $(PREFIX_THRIFT)/hack/thrift/engine.thrift && \
	sed -i 's:"engine":"evolution/quant/rpc/engine":' ./quant/rpc/engine/*/*.go

thrift-python:
	docker run -it --rm \
		-u $(USER):$(GROUP) \
		-v $(PWD)/engine:$(PREFIX_THRIFT)/engine \
		-v $(PWD)/hack:$(PREFIX_THRIFT)/hack \
		thrift:0.11.0 \
		thrift --gen py -out $(PREFIX_THRIFT)/engine/rpc $(PREFIX_THRIFT)/hack/thrift/engine.thrift

# engine
build-engine-basic:
	cd ./hack/dockerfile && docker build -f ./Dockerfile.engine -t $(PREFIX_DEVELOP):quant-engine . --network=host

# drone-release
local-drone:
	REGISTRY=$(PREFIX_REGISTRY) \
	GOPATH=/go/src \
	USERNAME=$(USERNAME) \
	PROJ=$(PROJ) \
	WORKPATH=$(DIRNAME) \
	GOVERSION=$(GOVERSION) \
	RELEASE_DIR=$(RELEASE_DIR) \
	DOCKER_USERNAME=username \
	DOCKER_PASSWORD=password \
	drone exec \
		--build-event "tag" \
		--commit-branch "release"

# local-release
# prepare golang vendor need proxy
prepare: prepare-backend prepare-frontend

build: check-release 
	make build-backend SYSTEM=time && \
	make build-backend SYSTEM=system && \
	make build-frontend local-release

prepare-backend:
	cd backend && glide update

build-backend: 
	docker run -it --rm \
		-u $(USER):$(GROUP) \
		-e CGO_ENABLED:0 \
		-v $(PWD)/$(RELEASE_DIR):/tmp/$(RELEASE_DIR) \
		-v $(PWD)/backend:/go/src/$(PROJ)/backend \
		-w /go/src/$(PROJ)/backend/$(SYSTEM) \
		golang:$(GOVERSION) \
		go build -o /tmp/$(RELEASE_DIR)/$(SYSTEM) \
		-a -tags netgo -installsuffix netgo -ldflags '-w' && \
	docker run -it --rm \
		-u $(USER):$(GROUP) \
		-e CGO_ENABLED:0 \
		-v $(PWD)/$(RELEASE_DIR):/tmp/$(RELEASE_DIR) \
		-v $(PWD)/backend:/go/src/$(PROJ)/backend \
		-w /go/src/$(PROJ)/backend/$(SYSTEM)/initial \
		golang:$(GOVERSION) \
		go build -o /tmp/$(RELEASE_DIR)/$(SYSTEM)-init-db \
		-a -tags netgo -installsuffix netgo -ldflags '-w' && \
	cp ./hack/release/Dockerfile.backend ./hack/release/Dockerfile && \
	sed -i "s:PROJ:$(PROJ):g" ./hack/release/Dockerfile && \
	sed -i "s:SYSTEM:$(SYSTEM):g" ./hack/release/Dockerfile && \
	docker build -f ./hack/release/Dockerfile -t $(PREFIX_RELEASE)-$(SYSTEM):$(TAG_RELEASE) . && \
	rm ./hack/release/Dockerfile

prepare-frontend:
	cd frontend && npm install

build-frontend:
	docker run -it --rm \
		-u $(USER):$(GROUP) \
		-v $(PWD)/$(RELEASE_DIR):/tmp/$(RELEASE_DIR) \
		-v $(PWD)/frontend:/tmp/frontend \
		-w /tmp/frontend \
		alexsuch/angular-cli:6.0.5 \
		ng build --environment=prod && \
	docker build -f ./hack/release/Dockerfile.frontend -t $(PREFIX_RELEASE)-frontend:$(TAG_RELEASE) .
