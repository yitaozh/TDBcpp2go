package main

/*
#cgo LDFLAGS: -lTDBAPI -lstdc++
#include "include/TDBAPI.h"
#include "include/TDBAPIStruct.h"
#include <stdlib.h>
#include <string.h>
*/
import "C"
import (
	"fmt"
	"time"
	"unsafe"
	"strconv"
	//"io"
	"os"
	"bytes"
	"encoding/binary"
)

type TDBDefine_ReqTick struct{
	chWindCode [32]byte         	//万得代码(ag1312.SHF)
	chCode [32]byte			//交易所代码(ag1312)
	nDate int                       //日期（自然日）
	nTime int                       //时间（HHMMSSmmm）例如94500000 表示 9点45分00秒000毫秒
	nPrice int                      //成交价((a double number + 0.00005) *10000)
	iVolume int     	     	//成交量
	iTurover int                	//成交额(元)
	nMatchItems int                 //成交笔数
	nInterest int                   //IOPV(基金)、利息(债券)
	chTradeFlag byte                //成交标志
	chBSFlag byte                 	//BS标志
 	iAccVolume int                 	//当日累计成交量
    	iAccTurover int             	//当日成交额(元)
 	nHigh int                   	//最高((a double number + 0.00005) *10000)
 	nLow int                      	//最低((a double number + 0.00005) *10000)
    	nOpen int                       //开盘((a double number + 0.00005) *10000)
    	nPreClose int                   //前收盘((a double number + 0.00005) *10000)

	//期货字段
 	nSettle int               	//结算价((a double number + 0.00005) *10000)
 	nPosition int           	//持仓量
	nCurDelta int                  	//虚实度
 	nPreSettle int                	//昨结算((a double number + 0.00005) *10000)
 	nPrePosition int              	//昨持仓

	//买卖盘字段
    	nAskPrice[10] int               //叫卖价((a double number + 0.00005) *10000)
 	nAskVolume[10] uint           	//叫卖量
    	nBidPrice[10] int               //叫买价((a double number + 0.00005) *10000)
 	nBidVolume[10] uint          	//叫买量
    	nAskAvPrice int                 //加权平均叫卖价(上海L2)((a double number + 0.00005) *10000)
    	nBidAvPrice int                 //加权平均叫买价(上海L2)((a double number + 0.00005) *10000)
  	iTotalAskVolume int         	//叫卖总量(上海L2)
  	iTotalBidVolume int         	//叫买总量(上海L2)

	//下面的字段指数使用
        nIndex int               	//不加权指数
        nStocks int             	//品种总数
        nUps int               		//上涨品种数
        nDowns int               	//下跌品种数
        nHoldLines int             	//持平品种数

	//保留字段
 	nResv1 int//保留字段1
 	nResv2 int//保留字段2
 	nResv3 int//保留字段3
}

type Define_Transaction struct{
	chWindCode[32]byte	//万得代码(ag1312.SHF)
    	chCode[32]byte        	//交易所代码(ag1312)
     	nDate int32             //日期（自然日）格式:YYMMDD
     	nTime int32             //成交时间(精确到毫秒HHMMSSmmm)
     	nIndex int32            //成交编号
    	chFunctionCode byte     //成交代码: 'C', 0
    	chOrderKind byte        //委托类别
    	chBSFlag byte           //BS标志
     	nTradePrice int32       //成交价格((a double number + 0.00005) *10000)
     	nTradeVolume int32      //成交数量
     	nAskOrder int32         //叫卖序号
     	nBidOrder int32         //叫买序号
}

func check(e error)  {
	if e!=nil{
		panic(e)
	}
}

func checkFilesExist(filename string)(bool){
	var exist = true
	if _,err := os.Stat(filename); os.IsNotExist(err){
		exist = false
	}
	return exist
}

func String2char(str string, des uintptr, sizeOf uintptr){
	bytes := []byte(str)
	for i:=0; i<len(bytes); i++{
		unit := (*C.char)(unsafe.Pointer(des))
		*unit = C.char(bytes[i])
		des += sizeOf
	}
}

func Char2byte(des uintptr, sizeOf uintptr, leng int)[256]byte{
	var bytes [256]byte
	for i:=0; i < leng; i++ {
		unit := (*C.char)(unsafe.Pointer(des))
		bytes[i] = byte(*unit)
		des += sizeOf
	}
	return bytes
}

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
/*
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
		fmt.Printf("成交量 iVolume:%d \n", pT.iVolume)
		fmt.Printf("成交额(元) iTurover:%d \n", pT.iTurover)
		fmt.Printf("成交笔数 nMatchItems:%d \n", pT.nMatchItems)
		fmt.Printf(" nInterest:%d \n", pT.nInterest)

		fmt.Printf("成交标志: chTradeFlag:%c \n", pT.chTradeFlag)
		fmt.Printf("BS标志: chBSFlag:%c \n", pT.chBSFlag)
		fmt.Printf("当日成交量: iAccVolume:%d \n", pT.iAccVolume)
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
*/

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
	//================================================================================
	l := unsafe.Sizeof(*pTransaction)
	buf := (*[1024]byte)(unsafe.Pointer(pTransaction))
	fmt.Println("Struct:", *pTransaction)
	fmt.Println("Bytes:", (*buf)[:l])
	fmt.Println("Length:", l)
	var transaction Define_Transaction
	//transaction.chWindCode = buf[0:32]
	//transaction.chCode = string(buf[32:64])
	binary.Read(bytes.NewBuffer(buf[0:32]), binary.LittleEndian, &transaction.chWindCode)
	binary.Read(bytes.NewBuffer(buf[32:64]), binary.LittleEndian, &transaction.chCode)
	binary.Read(bytes.NewBuffer(buf[64:68]), binary.LittleEndian, &transaction.nDate)
	binary.Read(bytes.NewBuffer(buf[68:72]), binary.LittleEndian, &transaction.nTime)
	binary.Read(bytes.NewBuffer(buf[72:76]), binary.LittleEndian, &transaction.nIndex)
	binary.Read(bytes.NewBuffer(buf[76:77]), binary.LittleEndian, &transaction.chFunctionCode)
	binary.Read(bytes.NewBuffer(buf[77:78]), binary.LittleEndian, &transaction.chOrderKind)
	binary.Read(bytes.NewBuffer(buf[78:79]), binary.LittleEndian, &transaction.chBSFlag)
	binary.Read(bytes.NewBuffer(buf[79:83]), binary.LittleEndian, &transaction.nTradePrice)
	binary.Read(bytes.NewBuffer(buf[83:87]), binary.LittleEndian, &transaction.nTradeVolume)
	binary.Read(bytes.NewBuffer(buf[87:91]), binary.LittleEndian, &transaction.nAskOrder)
	binary.Read(bytes.NewBuffer(buf[91:95]), binary.LittleEndian, &transaction.nBidOrder)
	fmt.Println(transaction)
	//================================================================================
	/*fmt.Println("-----------------------Transaction Data----------------------------")
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
	}*/
	//================================================================================
}

//tested good
/*
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
	//tmpPtr += 64+4+4+4+2
	fmt.Printf("委托类别: %c \n", tmpPtr)
	for i:=0; i<int(pCount); {
		for i:=0; i<4 ; {
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
			tmpPtr += sizeOf * 10000
		}
	}

}
*/

//tested
/*
func GetOrderQueue(hTdb C.THANDLE, szCode string, szMarketKey string, nDate int) {
	var req C.TDBDefine_ReqOrderQueue
	String2char(szCode, uintptr(unsafe.Pointer(&req.chCode)), unsafe.Sizeof(req.chCode[0]))
	String2char(szMarketKey, uintptr(unsafe.Pointer(&req.chMarketKey)), unsafe.Sizeof(req.chMarketKey[0]))
	req.nDate = C.int(nDate)
	req.nBeginTime = 0
	req.nEndTime = 0

	var pOrderQueue *C.TDBDefine_OrderQueue = nil
	var pCount C.int
	C.TDB_GetOrderQueue(hTdb, &req, &pOrderQueue, &pCount)

	fmt.Println("-------------------OrderQueue Data-------------");
	fmt.Printf("收到 %d 条委托队列消息，打印 1/1000 条\n", pCount);
	tmpPtr := uintptr(unsafe.Pointer(pOrderQueue))
	sizeOf := unsafe.Sizeof(*pOrderQueue)

	for i := 0; i < int(pCount); {
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
		tmpPtr += sizeOf * 10000

	}
}
*/

//指标公式

func UseEZFFormula(hTdb C.THANDLE) {
	fmt.Println("-------------------UseEZFFormula-------------");
	//公式的编写，请参考<<TRANSEND-TS-M0001 易编公式函数表V1(2).0-20110822.pdf>>
	strName := "KDJ"
	strContent := "INPUT:N(9), M1(3,1,100,2), M2(3);RSV:=(CLOSE-LLV(LOW,N))/(HHV(HIGH,N)-LLV(LOW,N))*100;K:SMA(RSV,M1,1);D:SMA(K,M2,1);J:3*K-2*D;"

	//添加公式到服务器并编译，若不过，会有错误返回
	var addRes *C.TDBDefine_AddFormulaRes = new(C.TDBDefine_AddFormulaRes)
	nErr := C.TDB_AddFormula(hTdb, C.CString(strName), C.CString(strContent),addRes)
	fmt.Printf("Add Formula Result:%s\n",Char2byte(uintptr(unsafe.Pointer(&addRes.chInfo)),unsafe.Sizeof(addRes.chInfo[0]),len(addRes.chInfo)))
//======================================================================================================
/*	var filename string = "./output1.txt"
	var f *os.File
	var err1 error
	if checkFilesExist(filename){
		f, err1 = os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
		fmt.Println("file exist")
	}else{
		f, err1 = os.Create(filename)
		fmt.Println("file not exist")
	}
	check(err1)
	bytes := Char2byte(uintptr(unsafe.Pointer(&addRes.chInfo)),unsafe.Sizeof(addRes.chInfo[0]),len(addRes.chInfo))
	str:= string(bytes[:])
	n, err1 := io.WriteString(f, str)
	check(err1)
	fmt.Printf("write %d bit\n", n)*//*
*/
//======================================================================================================
	//查询服务器上的公式，能看到我们刚才上传的"KDJ"
	var pEZFItem *C.TDBDefine_FormulaItem = nil
	var nItems C.int = 0
	//名字为空表示查�HB�j�RX 询服务器上所有的公式
	nErr = C.TDB_GetFormula(hTdb, nil, &pEZFItem, &nItems)
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
	tmpPtr_reqCalc := uintptr(unsafe.Pointer(&reqCalc.chFormulaName))
	sizeOf_reqCalc := unsafe.Sizeof(reqCalc.chFormulaName[0])
	String2char("KDJ", tmpPtr_reqCalc, sizeOf_reqCalc)

	tmpPtr_chParam := uintptr(unsafe.Pointer(&reqCalc.chParam))
	sizeOf_chParam := unsafe.Sizeof(reqCalc.chParam[0])
	String2char("N=9,M1=3,M2=3", tmpPtr_chParam, sizeOf_chParam)

	tmpPtr_chCode := uintptr(unsafe.Pointer(&reqCalc.chCode))
	sizeOf_chCode := unsafe.Sizeof(reqCalc.chCode[0])
	String2char("000001.SZ", tmpPtr_chCode, sizeOf_chCode)

	tmpPtr_chMarketKey := uintptr(unsafe.Pointer(&reqCalc.chMarketKey))
	sizeOf_chMarketKey := unsafe.Sizeof(reqCalc.chMarketKey[0])
	String2char("SZ-2-0", tmpPtr_chMarketKey, sizeOf_chMarketKey)

	reqCalc.nCycType = C.CYCTYPE(EZFCyc[0].nCyc)		//0表示日线
	reqCalc.nCycDef = C.int(EZFCyc[0].nCyc1)
	reqCalc.nCQFlag = C.REFILLFLAG(0)			//除权标志
	reqCalc.nCalcMaxItems = 4000 			//计算的最大数据量
	reqCalc.nResultMaxItems = 100			//传送的结果的最大数据量

	var pResult *C.TDBDefine_CalcFormulaRes = new(C.TDBDefine_CalcFormulaRes)
	nErr = C.TDB_CalcFormula(hTdb, &reqCalc, pResult)

	//判断错误代码
	fmt.Printf("计算结果有: %d 条\n", pResult.nRecordCount)

	var i C.int = 0
	var j C.int = 0
	//nFieldCount = 4
	for j=0; j<pResult.nFieldCount;j++{
		tmpPtr_chFieldName := uintptr(unsafe.Pointer(&pResult.chFieldName[j]))
		sizeOf_chFieldName := unsafe.Sizeof(pResult.chFieldName[j][1])
		fmt.Printf("%s  ",Char2byte(tmpPtr_chFieldName,sizeOf_chFieldName,len(pResult.chFieldName[j])))
	}
	fmt.Println()

	//输出数据
	for i=0; i<pResult.nRecordCount; i+=100{
		for j=0; j<pResult.nFieldCount;j++{
			fmt.Printf("%d  ", *pResult.dataFileds[j])
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

