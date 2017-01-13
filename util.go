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
	"strconv"
)

func copy(str string, des uintptr, sizeOf uintptr, len uint) {
	bytes := []byte(str)
	for i:=0; i<len; i++ {
		unit := (*C.char)(unsafe.Pointer(des))
		*unit = C.char(bytes[i])
		des += sizeOf
	}
}

func GetTickCount() int64 {
	return time.Now().Unix()
}

func array2str(arr *[]int, len uint) string {
	var str string
	for i:=0; i<len; i++ {
		if i==len-1 {
			str += strconv.Itoa(arr[i]) + " "
		}else {
			str += strconv.Itoa(arr[i]) + ","
		}
	}
	return str
}

//请求代码表
func GetCodeTable(hTdb C.THANDLE, szMarket string)  {
	var (
		pCodetable *[]C.TDBDefine_Code = nil
		pCount int
		outPutTable bool = true)
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
			fmt.Printf("市场代码 chWindCode:%s \n", pCt.chMarket)
			fmt.Printf("证券中文名称 chWindCode:%s \n", pCt.chCNName)
			fmt.Printf("证券英文名称 chWindCode:%s \n", pCt.chENName)
			fmt.Printf("证券类型 chWindCode:%d \n", pCt.nType)
			fmt.Println("----------------------------------------")
			tmpPtr += sizeOf
		}
	}
	C.TDB_Free(pCodetable)
}

func GetKData(hTdb C.THANDLE, szCode string, szMarket string, nBeginDate int, nEndDate int, nCycle int, nUserDef int, nCQFlag int, nAutoComplete int) {
	var req *C.TDBDefine_ReqKLine = new(C.TDBDefine_ReqKLine)
	copy(szCode, uintptr(unsafe.Pointer(&req.chCode)), unsafe.Sizeof(req.chCode[0]), 32)
	copy(szMarket, uintptr(unsafe.Pointer(&req.chMarketKey)), unsafe.Sizeof(req.chMarketKey[0]), 24)
	req.nCQFlag = (C.REFILLFLAG)nCQFlag  //除权标志，由用户定义
	req.nBeginDate = nBeginDate  //开始日期
	req.nEndDate = nEndDate;//结束日期
	req.nBeginTime = 0;//开始时间
	req.nEndTime = 0;//结束时间

	req.nCycType = (C.CYCTYPE)nCycle;
	req.nCycDef = 0;
	req.nAutoComplete = nAutoComplete;

	var kLine *C.TDBDefine_KLine = nil
	var pCount int

	C.TDB_GetKLine(hTdb,req,&kLine,&pCount)
	req=nil
	fmt.Println("---------------------------K Data--------------------")
	fmt.Printf("数据条数：%d,打印 1/100 条\n\n",pCount)
	tmpPtr := uintptr(unsafe.Pointer(kLine))
	sizeOf := unsafe.Sizeof(*kLine)
	for i:=0; i<pCount;  {
		kL := (*C.TDBDefine_KLine)(unsafe.Pointer(tmpPtr))
		fmt.Printf("WindCode:%s\n Code:%s\n Date:%d\n Time:%d\n Open:%d\n High:%d\n Low:%d\n Close:%d\n Volume:%lld\n Turover:%lld\n MatchItem:%d\n Interest:%d\n",
			kL.chWindCode,kL.chCode,kL.nDate,kL.nTime,kL.nOpen,kL.nHigh,kL.nLow,kL.nClose, kL.iVolume,kL.iTurover,kL.nMatchItems,kL.nInterest)
		tmpPtr += sizeOf*100
		i += 100
	}
}


