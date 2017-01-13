#!/usr/bin/env bash

g++ -c -fpic TDBAPI.cpp -I./

export LD_LIBRARY_PATH=$LD_LIBRARY_PATH:./

g++ -g -o TDBAPI.out TDBAPI.o -L./lib/ -lTDBAPI
