package main

/*
#cgo LDFLAGS: -lTDBAPI
#include "TDAPI.h"
 */
import "C"
import (
	"fmt"
	"time"
	"unsafe"
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
	tmpPtr := uintptr(unsafe.Pointer(pCodetable))
	sizeOf := unsafe.Sizeof(*pCodetable)
	if outPutTable {
		for i:=0; i<pCount; i++ {
			pCt := (*C.struct_TDBDefine_Code)(unsafe.Pointer(tmpPtr))
			fmt.Printf("交易所代码 chWindCode:%s \n", pCt.chCode)

			tmpPtr += sizeOf

		}
	}
}

func GetKData(hTdb C.THANDLE, szCode string, szMarket string, nBeginDate int, nEndDate int, nCycle int, nUserDef int, nCQFlag int, nAutoComplete int) {

}

func UseEZFFormula(hTdb C.THANDLE)  {
	strName := "KDJ"
	strContent := "INPUT:N(9), M1(3,1,100,2), M2(3);"
	"RSV:=(CLOSE-LLV(LOW,N))/(HHV(HIGH,N)-LLV(LOW,N))*100;"
	"K:SMA(RSV,M1,1);"
	"D:SMA(K,M2,1);"
	"J:3*K-2*D;"

	var addRes *C.TDBDefine_AddFormulaRes = new(C.TDBDefine_AddFormulaRes)
	var nErr int = C.TDB_AddFormula(hTdb, strName.c_str(), strContent.c_str(),addRes)
}

