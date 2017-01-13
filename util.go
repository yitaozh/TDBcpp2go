package main

/*
#cgo LDFLAGS: -lTDBAPI
#include "include/TDBAPI.h"
#include "include/TDBAPIStruct.h"
 */
import "C"
import (
	"fmt"
	"time"
	"unsafe"
	"strconv"

)

func main()  {
	fmt.Println("hello world")
}

func copystr(str string, des uintptr, sizeOf uintptr, len int) { //len需要去TDBStruct.h里查看
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

func array2str(arr []int, len int) string {
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

func array2str4C(arr [50]C.int, len C.int) string {
	var str string
	for i:=0; i<int(len); i++ {
		if i==int(len-1) {
			str += strconv.Itoa(int(arr[i])) + " "
		}else {
			str += strconv.Itoa(int(arr[i])) + ","
		}
	}
	return str
}


//请求代码表
func GetCodeTable(hTdb C.THANDLE, szMarket string)  {
	var (
		pCodetable *C.TDBDefine_Code = nil
		pCount C.int
		outPutTable bool = true)
	ret := C.TDB_GetCodeTable(hTdb, C.CString(szMarket), &pCodetable, &pCount)

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
		for i:=0; i<int(pCount); i++ {
			pCt := (*C.TDBDefine_Code)(unsafe.Pointer(tmpPtr))
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
	copystr(szCode, uintptr(unsafe.Pointer(&req.chCode)), unsafe.Sizeof(req.chCode[0]), 32)
	copystr(szMarket, uintptr(unsafe.Pointer(&req.chMarketKey)), unsafe.Sizeof(req.chMarketKey[0]), 24)
	//req.nCQFlag = C.REFILLFLAG(nCQFlag)  //除权标志，由用户定义
	req.nBeginDate = C.int(nBeginDate)  //开始日期
	req.nEndDate = C.int(nEndDate)//结束日期
	req.nBeginTime = 0//开始时间
	req.nEndTime = 0//结束时间

	//req.nCycType = C.CYCTYPE(nCycle)
	req.nCycDef = 0
	req.nAutoComplete = C.int(nAutoComplete)

	var kLine *C.TDBDefine_KLine = nil
	var pCount C.int

	C.TDB_GetKLine(hTdb,req,&kLine,&pCount)
	req=nil
	fmt.Println("---------------------------K Data--------------------")
	fmt.Printf("数据条数：%d,打印 1/100 条\n\n",pCount)
	tmpPtr := uintptr(unsafe.Pointer(kLine))
	sizeOf := unsafe.Sizeof(*kLine)
	for i:=0; i<int(pCount);  {
		kL := (*C.TDBDefine_KLine)(unsafe.Pointer(tmpPtr))
		fmt.Printf("WindCode:%s\n Code:%s\n Date:%d\n Time:%d\n Open:%d\n High:%d\n Low:%d\n Close:%v\n Volume:%v\n Turover:%d\n MatchItem:%d\n Interest:%d\n",
			kL.chWindCode,kL.chCode,kL.nDate,kL.nTime,kL.nOpen,kL.nHigh,kL.nLow,kL.nClose, kL.iVolume,kL.iTurover,kL.nMatchItems,kL.nInterest)
		tmpPtr += sizeOf*100
		i += 100
	}

	C.TDB_Free(kLine)
}


func GetTickData(hTdb C.THANDLE, szCode string, szMarket string, nDate int)  {
	var req C.TDBDefine_ReqTick
	copystr(szCode, uintptr(unsafe.Pointer(&req.chCode)), unsafe.Sizeof(req.chCode[0]), 32)
	copystr(szMarket, uintptr(unsafe.Pointer(&req.chMarketKey)), unsafe.Sizeof(req.chMarketKey[0]), 24)

	req.nDate = C.int(nDate)
	req.nBeginTime = 0
	req.nEndTime = 0

	var pTick *C.TDBDefine_Tick = nil
	var pCount C.int
	C.TDB_GetTick(hTdb,&req,&pTick, &pCount)

	fmt.Println("---------------------------------------Tick Data------------------------------------------")
	fmt.Printf("共收到 %d 条Tick数据， 打印 1/100 条：\n", pCount)

	tmpPtr := uintptr(unsafe.Pointer(pTick))
	sizeOf := unsafe.Sizeof(*pTick)
	for i:=0; i<int(pCount); {
		pT := (*C.TDBDefine_Tick)(unsafe.Pointer(tmpPtr))
		fmt.Printf("万得代码 chWindCode:%s \n", pT.chWindCode)
		fmt.Printf("日期 nDate:%d \n", pT.nDate)
		fmt.Printf("时间 nTime:%d \n", pT.nTime)

		fmt.Printf("成交价 nPrice:%d \n", pT.nPrice)
		fmt.Printf("成交量 iVolume:%v \n", pT.iVolume)
		fmt.Printf("成交额(元) iTurover:%v \n", pT.iTurover)
		fmt.Printf("成交笔数 nMatchItems:%d \n", pT.nMatchItems)
		fmt.Printf(" nInterest:%d \n", pT.nInterest)

		fmt.Printf("成交标志: chTradeFlag:%c \n", pT.chTradeFlag)
		fmt.Printf("BS标志: chBSFlag:%c \n", pT.chBSFlag)
		fmt.Printf("当日成交量: iAccVolume:%v \n", pT.iAccVolume)
		fmt.Printf("当日成交额: iAccTurover:%v \n", pT.iAccTurover)

		fmt.Printf("最高 nHigh:%d \n", pT.nHigh)
		fmt.Printf("最低 nLow:%d \n", pT.nLow)
		fmt.Printf("开盘 nOpen:%d \n", pT.nOpen)
		fmt.Printf("前收盘 nPreClose:%d \n", pT.nPreClose)

		//买卖盘字段
		var strOut string
		strOut = array2str(pT.nAskPrice, 10)
		fmt.Printf("叫卖价 nAskPrice:%s \n", strOut)
		strOut = array2str(pT.nAskVolume, 10)
		fmt.Printf("叫卖量 nAskVolume:%s \n", strOut)
		strOut = array2str(pT.nBidPrice, 10)
		fmt.Printf("叫买价 nBidPrice:%s \n", strOut)
		strOut = array2str(pT.nBidVolume, 10)
		fmt.Printf("叫买量 nBidVolume:%s \n", strOut)
		fmt.Printf("加权平均叫卖价 nAskAvPrice:%d \n", pT.nAskAvPrice)
		fmt.Printf("加权平均叫买价 nBidAvPrice:%d \n", pT.nBidAvPrice)
		fmt.Printf("叫卖总量 iTotalAskVolume:%v \n", pT.iTotalAskVolume)
		fmt.Printf("叫买总量 iTotalBidVolume:%v \n", pT.iTotalBidVolume)

/*
		//期货字段
		fmt.Printf("结算价 nSettle:%d \n", pT.nSettle)
		fmt.Printf("持仓量 nPosition:%d \n", pT.nPosition)
		fmt.Printf("虚实度 nCurDelta:%d \n", pT.nCurDelta)
		fmt.Printf("昨结算 nPreSettle:%d \n", pT.nPreSettle)
		fmt.Printf("昨持仓 nPrePosition:%d \n", pT.nPrePosition)

		//指数
		fmt.Printf("不加权指数 nIndex:%d \n", pT.nIndex)
		fmt.Printf("品种总数 nStocks:%d \n", pT.nStocks)
		fmt.Printf("上涨品种数 nUps:%d \n", pT.nUps)
		fmt.Printf("下跌品种数 nDowns:%d \n", pT.nDowns)
		fmt.Printf("持平品种数 nHoldLines:%d \n", pT.nHoldLines)
		*/

		fmt.Println("-------------------------------")
		i += 100
		tmpPtr += sizeOf*100
	}
	C.TDB_Free(pTick) //释放
}


func GetTransaction(hTdb C.THANDLE, szCode string, szMarketKey string, nDate int)  {
	var req C.TDBDefine_ReqTransaction
	copystr(szCode, uintptr(unsafe.Pointer(&req.chCode)), unsafe.Sizeof(req.chCode[0]), 32)
	copystr(szMarketKey, uintptr(unsafe.Pointer(&req.chMarketKey)), unsafe.Sizeof(req.chMarketKey[0]), 24)
	req.nDate = C.int(nDate)
	req.nBeginTime = 0
	req.nEndTime = 0

	var pTransaction *C.TDBDefine_Transaction = nil
	var pCount C.int
	C.TDB_GetTransaction(hTdb,&req, &pTransaction, &pCount)

	fmt.Println("---------------------------------------Transaction Data------------------------------------------")
	fmt.Printf("收到 %d 条逐笔成交消息，打印 1/10000 条\n", pCount)
	tmpPtr := uintptr(unsafe.Pointer(pTransaction))
	sizeOf := unsafe.Sizeof(*pTransaction)
	for i:=0; i<int(pCount); {
		pT := (*C.TDBDefine_Transaction)(unsafe.Pointer(tmpPtr))
		fmt.Printf("成交时间(Date): %d \n", pT.nDate)
		fmt.Printf("成交时间: %d \n", pT.nTime)
		fmt.Printf("成交代码: %c \n", pT.chFunctionCode)
		fmt.Printf("委托类别: %c \n", pT.chOrderKind)
		fmt.Printf("BS标志: %c \n", pT.chBSFlag)
		fmt.Printf("成交价格: %d \n", pT.nTradePrice)
		fmt.Printf("成交数量: %d \n", pT.nTradeVolume)
		fmt.Printf("叫卖序号: %d \n", pT.nAskOrder)
		fmt.Printf("叫买序号: %d \n", pT.nBidOrder)
		fmt.Println("---------------------------------------------")
		//fmt.Printf("成交编号: %d \n", pT.nBidOrder)
		i += 10000
		tmpPtr += sizeOf*10000
	}
	C.TDB_Free(pTransaction)
}
/*

func getOrder(hTdb C.THANDLE, szCode string, szMarketKey string, nDate int)  {
	var req C.TDBDefine_ReqOrder
	copystr(szCode, uintptr(unsafe.Pointer(&req.chCode)), unsafe.Sizeof(req.chCode[0]), 32)
	copystr(szMarketKey, uintptr(unsafe.Pointer(&req.chMarketKey)), unsafe.Sizeof(req.chMarketKey[0]), 24)
	req.nDate = C.int(nDate)
	req.nBeginTime = 0
	req.nEndTime = 0

	var pOrder *C.TDBDefine_Order = nil
	var pCount C.int
	C.TDB_GetOrder(hTdb,&req, &pOrder, &pCount)

	fmt.Println("---------------------------------------Transaction Data------------------------------------------")
	fmt.Printf("收到 %d 条逐笔委托消息，打印 1/10000 条\n", pCount)
	tmpPtr := uintptr(unsafe.Pointer(pOrder))
	sizeOf := unsafe.Sizeof(*pOrder)
	for i:=0; i<int(pCount); {
		pO := (*C.TDBDefine_Order)(unsafe.Pointer(tmpPtr))
		fmt.Printf("订单时间(Date): %d \n", pO.nDate)
		fmt.Printf("委托时间(HHMMSSmmm): %d \n", pO.nTime)
		fmt.Printf("委托编号: %d \n", pO.nOrder)
		fmt.Printf("委托类别: %c \n", pO.chOrderKind)
		fmt.Printf("委托代码: %c \n", pO.chFunctionCode)
		fmt.Printf("委托价格: %d \n", pO.nOrderPrice)
		fmt.Printf("委托数量: %d \n", pO.nOrderVolume)
		fmt.Println("---------------------------------------------")

		i += 10000
		tmpPtr += sizeOf*10000
	}
	C.TDB_Free(pOrder)
}

func GetOrderQueue(hTdb C.THANDLE, szCode string, szMarketKey string, nDate int)  {
	var req C.TDBDefine_ReqOrderQueue
	copystr(szCode, uintptr(unsafe.Pointer(&req.chCode)), unsafe.Sizeof(req.chCode[0]), 32)
	copystr(szMarketKey, uintptr(unsafe.Pointer(&req.chMarketKey)), unsafe.Sizeof(req.chMarketKey[0]), 24)
	req.nDate = C.int(nDate)
	req.nBeginTime = 0
	req.nEndTime = 0

	var pOrderQueue *C.TDBDefine_OrderQueue = nil
	var pCount C.int
	C.TDB_GetOrderQueue(hTdb,&req, &pOrderQueue, &pCount)

	fmt.Println("-------------------OrderQueue Data-------------");
	fmt.Printf("收到 %d 条委托队列消息，打印 1/1000 条\n", pCount);
	tmpPtr := uintptr(unsafe.Pointer(pOrderQueue))
	sizeOf := unsafe.Sizeof(*pOrderQueue)
	for i:=0; i<int(pCount); {
		pOQ := (*C.TDBDefine_OrderQueue)(unsafe.Pointer(tmpPtr))
		fmt.Printf("订单时间(Date): %d \n", pOQ.nDate)
		fmt.Printf("订单时间(HHMMSS): %d \n", pOQ.nTime)
		fmt.Printf("买卖方向('B':Bid 'A':Ask): %c \n", pOQ.nSide)
		fmt.Printf("成交价格: %d \n", pOQ.nPrice)
		fmt.Printf("订单数量: %d \n", pOQ.nOrderItems)
		fmt.Printf("明细个数: %d \n", pOQ.nABItems)
		fmt.Printf("订单明细: %s \n", array2str4C(pOQ.nABVolume, pOQ.nABItems))
		fmt.Println("---------------------------------------------")
		i += 10000
		tmpPtr += sizeOf*10000
	}
	C.TDB_Free(pOrderQueue)
}
*/


