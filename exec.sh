#!/bin/bash

if [ ! -d ./logs  ]; then
    mkdir -p ./logs
fi
script=bash runmode=dev LD_LIBRARY_PATH=/usr/lib python3 ./dags/route.py 
