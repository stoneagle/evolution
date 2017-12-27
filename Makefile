.PHONY: run-web, stop-web rm-web

PWD := $(shell pwd)
USER := $(shell id -u)
GROUP := $(shell id -g)
DATE := $(shell date "+%F")

run-airflow: 
	cd hack && docker-compose -f docker-compose-local.yml -p "airflow-$(USER)" up -d
stop-airflow: 
	cd hack && docker-compose -f docker-compose-local.yml -p "airflow-$(USER)" stop 
rm-airflow: 
	cd hack && docker-compose -f docker-compose-local.yml -p "airflow-$(USER)" rm 

run-grafana: 
	cd hack && docker-compose -f docker-compose-grafana.yml -p "grafana-$(USER)" up -d
stop-grafana: 
	cd hack && docker-compose -f docker-compose-grafana.yml -p "grafana-$(USER)" stop 
rm-grafana: 
	cd hack && docker-compose -f docker-compose-grafana.yml -p "grafana-$(USER)" rm 

local-bash:
	nohup bash ./exec.sh > ./logs/$(DATE).log 2>&1 &
bitmex-test:
	LD_LIBRARY_PATH=/usr/lib python3 ./dags/route.py -s bitmex -f test 
bitmex-watch:
	LD_LIBRARY_PATH=/usr/lib python3 ./dags/route.py -s bitmex -f watch 
ashare-test:
	LD_LIBRARY_PATH=/usr/lib python3 ./dags/route.py -s ashare -f test 
ashare-watch:
	LD_LIBRARY_PATH=/usr/lib python3 ./dags/route.py -s ashare -f watch

build-talib:
	mkdir tmp && cd tmp && \
		wget http://prdownloads.sourceforge.net/ta-lib/ta-lib-0.4.0-src.tar.gz && \
		tar -xvf ./ta-lib-0.4.0-src.tar.gz && cd ./ta-lib && ./configure --prefix=/usr && \
		make && sudo make install && cd ../.. && \
		rm -rf ./tmp 

local-python:
	pip3 install --index https://pypi.mirrors.ustc.edu.cn/simple/ --user -r ./hack/dockerfile/requirements.utils.txt && \
	pip3 install --index https://pypi.mirrors.ustc.edu.cn/simple/ --user -r ./hack/dockerfile/requirements.tushare.txt

build-img:
	cd hack/dockerfile && docker build -f ./Dockerfile -t puckel/docker-airflow:1.8.2-assets-2 .

init-plugin:
	cd plugin/ashare && npm install --registry=http://rgistry.npm.taobao.org && ./node_modules/grunt/bin/grunt

build-grafana:
	cd hack/dockerfile && docker build -f ./Dockerfile-grafana -t grafana/grafana:4.6.2-1000 .
