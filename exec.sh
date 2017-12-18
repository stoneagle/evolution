#!/bin/bash

if [ ! -d ./logs  ]; then
    mkdir -p ./logs
fi
runmode=dev LD_LIBRARY_PATH=/usr/lib python3 ./dags/scripts/bash.py 
