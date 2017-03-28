package main

/*
#cgo LDFLAGS:-lTDBAPI
#include "include/TDBAPI.h"
#include "include/TDBAPIStruct.h"
#include <stdlib.h>
#include <string.h>
*/
import "C"

import (
	"fmt"
)


func main(){

	JsonParse := NewJsonStruct()

	cfg := conf{}

	JsonParse.Load("conf.json", &cfg)

	fmt.Println(cfg)

	//connect to TDB server and get handle
	hTdb := TDBConnection(cfg)

	//connect to Influxdb server and get handle
	c := InfluxConnection(cfg)

	writeData(hTdb, cfg, c)
}
