package main

/*
#cgo LDFLAGS: -lTDBAPI
#include "include/TDBAPI.h"
#include "include/TDBAPIStruct.h"
#include <stdlib.h>
 */
import "C"
import (
	"fmt"
	"time"
	"unsafe"
	"strconv"

)


func GetTickCount() int64 {
	return time.Now().Unix()
}

func array2str(arr [10]C.int, len int) string {
	var str string
	for i:=0; i<len; i++ {
		if i==len-1 {
			str += strconv.Itoa(int(arr[i])) + " "
		}else {
			str += strconv.Itoa(int(arr[i])) + ","
		}
	}
	return str
}

func array2str4uint(arr [10]C.uint, len int) string {
	var str string
	for i:=0; i<len; i++ {
		if i==len-1 {
			str += strconv.Itoa(int(arr[i])) + " "
		}else {
			str += strconv.Itoa(int(arr[i])) + ","
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
			fmt.Printf("万得代码 chWindCode:%s \n", Char2byte(uintptr(unsafe.Pointer(&pCt.chWindCode)),unsafe.Sizeof(pCt.chWindCode[0]),len(pCt.chWindCode)))
			fmt.Printf("交易所代码 chWindCode:%s \n", Char2byte(uintptr(unsafe.Pointer(&pCt.chCode)),unsafe.Sizeof(pCt.chCode[0]),len(pCt.chCode)))
			fmt.Printf("市场代码 chMarket:%s \n", Char2byte(uintptr(unsafe.Pointer(&pCt.chMarket)),unsafe.Sizeof(pCt.chMarket[0]),len(pCt.chMarket)))
			fmt.Printf("证券中文名称 chCNName:%s \n", Char2byte(uintptr(unsafe.Pointer(&pCt.chCNName)),unsafe.Sizeof(pCt.chCNName[0]),len(pCt.chCNName)))
			fmt.Printf("证券英文名称 chENName:%s \n", Char2byte(uintptr(unsafe.Pointer(&pCt.chENName)),unsafe.Sizeof(pCt.chENName[0]),len(pCt.chENName)))
			fmt.Printf("证券类型 nType:%d \n", pCt.nType)
			fmt.Println("----------------------------------------")
			tmpPtr += sizeOf
		}
	}

}
//tested good
func GetKData(hTdb C.THANDLE, szCode string, szMarket string, nBeginDate int, nEndDate int, nCycle int, nUserDef int, nCQFlag int, nAutoComplete int) {
	var req *C.TDBDefine_ReqKLine = new(C.TDBDefine_ReqKLine)
	String2char(szCode,uintptr(unsafe.Pointer(&req.chCode)),unsafe.Sizeof(req.chCode[0]))
	String2char(szMarket,uintptr(unsafe.Pointer(&req.chMarketKey)),unsafe.Sizeof(req.chMarketKey[0]))
	req.nCQFlag = C.REFILLFLAG(nCQFlag)  //除权标志，由用户定义
	req.nBeginDate = C.int(nBeginDate)  //开始日期
	req.nEndDate = C.int(nEndDate)//结束日期
	req.nBeginTime = 0//开始时间
	req.nEndTime = 0//结束时间

	req.nCycType = C.CYCTYPE(nCycle)
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
			Char2byte(uintptr(unsafe.Pointer(&kL.chWindCode)),unsafe.Sizeof(kL.chWindCode[0]),len(kL.chWindCode)),//kL.chWindCode
			Char2byte(uintptr(unsafe.Pointer(&kL.chCode)),unsafe.Sizeof(kL.chCode[0]),len(kL.chCode)),//kL.chCode
			kL.nDate, kL.nTime, kL.nOpen, kL.nHigh, kL.nLow, kL.nClose, kL.iVolume, kL.iTurover, kL.nMatchItems, kL.nInterest )
		tmpPtr += sizeOf*100
		i += 100
	}
}

//tested good
func GetTickData(hTdb C.THANDLE, szCode string, szMarket string, nDate int)  {
	var req C.TDBDefine_ReqTick
	String2char(szCode,uintptr(unsafe.Pointer(&req.chCode)),unsafe.Sizeof(req.chCode[0]))
	String2char(szMarket,uintptr(unsafe.Pointer(&req.chMarketKey)),unsafe.Sizeof(req.chMarketKey[0]))

	req.nDate = C.int(nDate)
	req.nBeginTime = 0
	req.nEndTime = 0

	var pTick *C.TDBDefine_Tick = nil
	var pCount C.int
	C.TDB_GetTick(hTdb,&req,&pTick, &pCount)

	fmt.Println("------------------------Tick Data---------------------------")
	fmt.Printf("共收到 %d 条Tick数据， 打印 1/100 条：\n", pCount)

	tmpPtr := uintptr(unsafe.Pointer(pTick))
	sizeOf := unsafe.Sizeof(*pTick)
	for i:=0; i<int(pCount); {
		pT := (*C.TDBDefine_Tick)(unsafe.Pointer(tmpPtr))
		fmt.Printf("万得代码 chWindCode:%s \n", Char2byte(uintptr(unsafe.Pointer(&pT.chWindCode)),unsafe.Sizeof(pT.chWindCode[0]),len(pT.chWindCode)))
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
		strOut = array2str4uint(pT.nAskVolume, 10)
		fmt.Printf("叫卖量 nAskVolume:%s \n", strOut)
		strOut = array2str(pT.nBidPrice, 10)
		fmt.Printf("叫买价 nBidPrice:%s \n", strOut)
		strOut = array2str4uint(pT.nBidVolume, 10)
		fmt.Printf("叫买量 nBidVolume:%s \n", strOut)
		fmt.Printf("加权平均叫卖价 nAskAvPrice:%d \n", pT.nAskAvPrice)
		fmt.Printf("加权平均叫买价 nBidAvPrice:%d \n", pT.nBidAvPrice)
		fmt.Printf("叫卖总量 iTotalAskVolume:%v \n", pT.iTotalAskVolume)
		fmt.Printf("叫买总量 iTotalBidVolume:%v \n", pT.iTotalBidVolume)


		//期货字段
//		fmt.Printf("结算价 nSettle:%d \n", pT.nSettle)
//		fmt.Printf("持仓量 nPosition:%d \n", pT.nPosition)
//		fmt.Printf("虚实度 nCurDelta:%d \n", pT.nCurDelta)
//		fmt.Printf("昨结算 nPreSettle:%d \n", pT.nPreSettle)
//		fmt.Printf("昨持仓 nPrePosition:%d \n", pT.nPrePosition)

		//指数
//		fmt.Printf("不加权指数 nIndex:%d \n", pT.nIndex)
//		fmt.Printf("品种总数 nStocks:%d \n", pT.nStocks)
//		fmt.Printf("上涨品种数 nUps:%d \n", pT.nUps)
//		fmt.Printf("下跌品种数 nDowns:%d \n", pT.nDowns)
//		fmt.Printf("持平品种数 nHoldLines:%d \n", pT.nHoldLines)


		fmt.Println("--------------------------------------")
		i += 100
		tmpPtr += sizeOf*100
	}
}

//tested good
func GetTransaction(hTdb C.THANDLE, szCode string, szMarketKey string, nDate int)  {
	var req C.TDBDefine_ReqTransaction
	String2char(szCode,uintptr(unsafe.Pointer(&req.chCode)),unsafe.Sizeof(req.chCode[0]))
	String2char(szMarketKey,uintptr(unsafe.Pointer(&req.chMarketKey)),unsafe.Sizeof(req.chMarketKey[0]))
	req.nDate = C.int(nDate)
	req.nBeginTime = 0
	req.nEndTime = 0

	var pTransaction *C.TDBDefine_Transaction = nil
	var pCount C.int
	C.TDB_GetTransaction(hTdb,&req, &pTransaction, &pCount)

	fmt.Println("-----------------------Transaction Data----------------------------")
	fmt.Printf("收到 %d 条逐笔成交消息，打印 1/10000 条\n", pCount)
	tmpPtr := uintptr(unsafe.Pointer(pTransaction))
	sizeOf := unsafe.Sizeof(*pTransaction)
	for i:=0; i<int(pCount); {
		pT := (*C.TDBDefine_Transaction)(unsafe.Pointer(tmpPtr))
		fmt.Printf("成交时间(Date): %d \n", pT.nDate)
		fmt.Printf("成交时间: %d \n", pT.nTime)
		fmt.Printf("成交代码: %c \n", byte(pT.chFunctionCode))
		fmt.Printf("委托类别: %c \n", byte(pT.chOrderKind))
		fmt.Printf("BS标志: %c \n", byte(pT.chBSFlag))
		fmt.Printf("成交价格: %d \n", pT.nTradePrice)
		fmt.Printf("成交数量: %d \n", pT.nTradeVolume)
		fmt.Printf("叫卖序号: %d \n", pT.nAskOrder)
		fmt.Printf("叫买序号: %d \n", pT.nBidOrder)
		fmt.Println("---------------------------------------------")
		//fmt.Printf("成交编号: %d \n", pT.nBidOrder)
		i += 10000
		tmpPtr += sizeOf*10000
	}
}

//tested good
func GetOrder(hTdb C.THANDLE, szCode string, szMarketKey string, nDate int)  {
	var req C.TDBDefine_ReqOrder
	String2char(szCode,uintptr(unsafe.Pointer(&req.chCode)),unsafe.Sizeof(req.chCode[0]))
	String2char(szMarketKey,uintptr(unsafe.Pointer(&req.chMarketKey)),unsafe.Sizeof(req.chMarketKey[0]))
	req.nDate = C.int(nDate)
	req.nBeginTime = 0
	req.nEndTime = 0

	var pOrder *C.TDBDefine_Order = nil
	var pCount C.int
	C.TDB_GetOrder(hTdb,&req, &pOrder, &pCount)

	fmt.Println("-------------------------Transaction Data--------------------------")
	fmt.Printf("收到 %d 条逐笔委托消息，打印 1/10000 条\n", pCount)
	tmpPtr := uintptr(unsafe.Pointer(pOrder))
	sizeOf := unsafe.Sizeof(*pOrder)
	for i:=0; i<int(pCount); {
		pO := (*C.TDBDefine_Order)(unsafe.Pointer(tmpPtr))
		fmt.Printf("订单时间(Date): %d \n", pO.nDate)
		fmt.Printf("委托时间(HHMMSSmmm): %d \n", (*C.int)(unsafe.Pointer(uintptr(unsafe.Pointer(pOrder)) + unsafe.Offsetof(pOrder.nTime))))
		fmt.Printf("委托编号: %d \n", pO.nOrder)
		fmt.Printf("委托类别: %c \n", pO.chOrderKind)
		fmt.Printf("委托代码: %c \n", pO.chFunctionCode)
		fmt.Printf("委托价格: %d \n", pO.nOrderPrice)
		fmt.Printf("委托数量: %d \n", pO.nOrderVolume)
		fmt.Println("---------------------------------------------")
		i += 10000
		tmpPtr += sizeOf*10000
	}

}
//tested
func GetOrderQueue(hTdb C.THANDLE, szCode string, szMarketKey string, nDate int)  {
	var req C.TDBDefine_ReqOrderQueue
	String2char(szCode,uintptr(unsafe.Pointer(&req.chCode)),unsafe.Sizeof(req.chCode[0]))
	String2char(szMarketKey,uintptr(unsafe.Pointer(&req.chMarketKey)),unsafe.Sizeof(req.chMarketKey[0]))
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
}

//指标公式
func UseEZFFormula(hTdb C.THANDLE) {
	//公式的编写，请参考<<TRANSEND-TS-M0001 易编公式函数表V1(2).0-20110822.pdf>>
	strName := "KDJ"
	strContent := "INPUT:N(9), M1(3,1,100,2), M2(3);RSV:=(CLOSE-LLV(LOW,N))/(HHV(HIGH,N)-LLV(LOW,N))*100;K:SMA(RSV,M1,1);D:SMA(K,M2,1);J:3*K-2*D;"

	//添加公式到服务器并编译，若不过，会有错误返回
	var addRes *C.TDBDefine_AddFormulaRes = new(C.TDBDefine_AddFormulaRes)
	nErr := C.TDB_AddFormula(hTdb, C.CString(strName), C.CString(strContent),addRes)
	fmt.Printf("Add Formula Result:%s\n",Char2byte(uintptr(unsafe.Pointer(&addRes.chInfo)),unsafe.Sizeof(addRes.chInfo[0]),len(addRes.chInfo)))

	//查询服务器上的公式，能看到我们刚才上传的"KDJ"
	var pEZFItem *C.TDBDefine_FormulaItem = nil
	var nItems C.int = 0
	//名字为空表示查询服务器上所有的公式
	nErr = C.TDB_GetFormula(hTdb, nil, &pEZFItem, &nItems)
	fmt.Println(nErr)
	tmpPtr := uintptr(unsafe.Pointer(pEZFItem))
	sizeOf := unsafe.Sizeof(*pEZFItem)
	for i:=0; i<int(nItems); i++{
		pEZF := (*C.TDBDefine_FormulaItem)(unsafe.Pointer(tmpPtr))
		fmt.Printf("公式名称：%s, 参数:%s \n",
			Char2byte(uintptr(unsafe.Pointer(&pEZF.chFormulaName)),unsafe.Sizeof(pEZF.chFormulaName[0]),len(pEZF.chFormulaName)),
			Char2byte(uintptr(unsafe.Pointer(&pEZF.chParam)),unsafe.Sizeof(pEZF.chParam[0]),len(pEZF.chParam)),
			)
		tmpPtr += sizeOf
	}

	type EZFCycDefine struct {
		chName string
		nCyc   int
		nCyc1  int
	}
	var EZFCyc[5] EZFCycDefine
	EZFCyc[0] = EZFCycDefine{"日线", 2, 0}
	EZFCyc[1] = EZFCycDefine{"30分", 0, 30}
	EZFCyc[2] = EZFCycDefine{"5分钟", 0, 5}
	EZFCyc[3] = EZFCycDefine{"1分钟", 0, 1}
	EZFCyc[4] = EZFCycDefine{"15秒", 11, 15}

	//获取公式的计算结果
	var reqCalc C.TDBDefine_ReqCalcFormula
	fmt.Println(len(reqCalc.chFormulaName))
	tmpPtr_reqCalc := uintptr(unsafe.Pointer(&reqCalc.chFormulaName))
	sizeOf_reqCalc := unsafe.Sizeof(reqCalc.chFormulaName)

	String2char("KDJ", tmpPtr_reqCalc, sizeOf_reqCalc)
	tmpPtr_chParam := uintptr(unsafe.Pointer(&reqCalc.chParam))
	sizeOf_chParam := unsafe.Sizeof(reqCalc.chParam)
	String2char("N=9,M1=3,M2=3", tmpPtr_chParam, sizeOf_chParam)
	tmpPtr_chCode := uintptr(unsafe.Pointer(&reqCalc.chCode))
	sizeOf_chCode := unsafe.Sizeof(reqCalc.chCode)
	String2char("000001.SZ", tmpPtr_chCode, sizeOf_chCode)
	tmpPtr_chMarketKey := uintptr(unsafe.Pointer(&reqCalc.chMarketKey))
	sizeOf_chMarketKey := unsafe.Sizeof(reqCalc.chMarketKey)
	String2char("SZ-2-0", tmpPtr_chMarketKey, sizeOf_chMarketKey)

	reqCalc.nCycType = C.CYCTYPE(EZFCyc[0].nCyc)		//0表示日线
	reqCalc.nCycDef = C.int(EZFCyc[0].nCyc1)
	reqCalc.nCQFlag = 0				//除权标志
	reqCalc.nCalcMaxItems = 4000 			//计算的最大数据量
	reqCalc.nResultMaxItems = 100			//传送的结果的最大数据量

	var pResult *C.TDBDefine_CalcFormulaRes = new(C.TDBDefine_CalcFormulaRes)
	nErr = C.TDB_CalcFormula(hTdb, &reqCalc, pResult)

	//判断错误代码
	fmt.Printf("计算结果有: %d 条:\n", pResult.nRecordCount)

	var i C.int = 0
	var j C.int = 0
	for j=0; j<pResult.nFieldCount;j++{
		tmpPtr_chFieldName := uintptr(unsafe.Pointer(&pResult.chFieldName[j]))
		sizeOf_chFieldName := unsafe.Sizeof(pResult.chFieldName[j])
		fmt.Printf("%s  ",Char2byte(tmpPtr_chFieldName,sizeOf_chFieldName,len(pResult.chFieldName[j])))
	}
	fmt.Println();
	fmt.Println();
	//输出数据
	for i=0; i<pResult.nRecordCount; i+=100{
		for j=0; j<pResult.nFieldCount;j++{
			fmt.Printf("%c  ", pResult.dataFileds[i*pResult.nFieldCount+j])
		}
		fmt.Println()
	}

	//删除之前上传的公式指标
	var pDel C.TDBDefine_DelFormulaRes

	nErr = C.TDB_DeleteFormula(hTdb, C.CString("KDJ"), &pDel)
	fmt.Printf("删除指标信息:%s\n", Char2byte(uintptr(unsafe.Pointer(&pDel.chInfo)),unsafe.Sizeof(pDel.chInfo[1]),len(pDel.chInfo)))
	fmt.Printf("Error:%d\n", int(nErr))

	C.TDB_ReleaseCalcFormula(pResult)
}



