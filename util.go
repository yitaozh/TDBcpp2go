package main

/*
#cgo LDFLAGS: -lTDBAPI
#include "TDAPI.h"
 */
import "C"
import (
	"fmt"
	"time"
)

func GetTickCount() int64 {
	return time.Now().Unix()
}

func int2str(n int) string {

}

func array2str(arr *[]int, len uint) string {
	var str string
	for i:=0; i<len; i++ {
		if i==len-1 {
			str += int2str(arr[i]) + " "
		}else {
			str += int2str(arr[i]) + ","
		}
	}
	return str
}

//请求代码表
func GetCodeTable(hTdb C.THANDLE, szMarket string)  {
	var (	pCodetable *[]C.TDBDefine_Code
		pCount int
		outPutTable bool)
	ret := C.TDB_GetCodeTable(hTdb, szMarket, &pCodetable, &pCount)

	if ret == C.TDB_NO_DATA {
		fmt.Println("无代码表！")
		return
	}

	fmt.Println("---------------------------Code Table--------------------")
	fmt.Printf("收到代码表项数：%d，\n\n",pCount)
	//输出
	if outPutTable {
		for i:=0; i<pCount; i++ {
			fmt.Printf("交易所代码 chWindCode:%s \n", pCodetable[i].chCode)
		}
	}
}