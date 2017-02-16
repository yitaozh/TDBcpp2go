package main

/*
#cgo LDFLAGS: -lTDBAPI
#include "include/TDBAPI.h"
#include "include/TDBAPIStruct.h"
#include <stdlib.h>
#include <string.h>
*/
import "C"
import (
	"fmt"
	"unsafe"
	"time"
	"strconv"
	"log"
	//"io"
	"os"
	"bytes"
	"encoding/binary"
	"github.com/influxdata/influxdb/client/v2"
	"math/rand"
	//"code.google.com/p/mahonia"
)

const (
	MyDB = "square_holes"
	username = "bubba"
	password = "bumblebeetuna"
)

type Define_Tick struct{
	chWindCode [32]byte         	//万得代码(ag1312.SHF)
	chCode [32]byte			//交易所代码(ag1312)
	nDate int32                       //日期（自然日）
	nTime int32                       //时间（HHMMSSmmm）例如94500000 表示 9点45分00秒000毫秒
	nPrice int32                      //成交价((a double number + 0.00005) *10000)
	iVolume int64     	     	//成交量
	iTurover int64                	//成交额(元)
	nMatchItems int32                 //成交笔数
	nInterest int32                   //IOPV(基金)、利息(债券)
	chTradeFlag byte                //成交标志
	chBSFlag byte                 	//BS标志
 	iAccVolume int64                 	//当日累计成交量
    	iAccTurover int64             	//当日成交额(元)
 	nHigh int32                   	//最高((a double number + 0.00005) *10000)
 	nLow int32                      	//最低((a double number + 0.00005) *10000)
    	nOpen int32                       //开盘((a double number + 0.00005) *10000)
    	nPreClose int32                  //前收盘((a double number + 0.00005) *10000)

	//期货字段
 	nSettle int32               	//结算价((a double number + 0.00005) *10000)
 	nPosition int32           	//持仓量
	nCurDelta int32                  	//虚实度
 	nPreSettle int32                	//昨结算((a double number + 0.00005) *10000)
 	nPrePosition int32              	//昨持仓

	//买卖盘字段
    	nAskPrice[10] int32               //叫卖价((a double number + 0.00005) *10000)
 	nAskVolume[10] uint32           	//叫卖量
    	nBidPrice[10] int32               //叫买价((a double number + 0.00005) *10000)
 	nBidVolume[10] uint32          	//叫买量
    	nAskAvPrice int32                 //加权平均叫卖价(上海L2)((a double number + 0.00005) *10000)
    	nBidAvPrice int32                 //加权平均叫买价(上海L2)((a double number + 0.00005) *10000)
  	iTotalAskVolume int64         	//叫卖总量(上海L2)
  	iTotalBidVolume int64         	//叫买总量(上海L2)

	//下面的字段指数使用
        nIndex int32               	//不加权指数
        nStocks int32             	//品种总数
        nUps int32               		//上涨品种数
        nDowns int32               	//下跌品种数
        nHoldLines int32             	//持平品种数

	//保留字段
 	nResv1 int32//保留字段1
 	nResv2 int32//保留字段2
 	nResv3 int32//保留字段3
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

type Define_Order struct{
	chWindCode[32]byte        //万得代码(ag1312.SHF)
	chCode[32]byte            //交易所代码(ag1312)
	nDate int32                 //日期（自然日）格式YYMMDD
	nTime int32                //委托时间（精确到毫秒HHMMSSmmm）
	nIndex int32                //委托编号
	nOrder int32                //交易所委托号
	chOrderKind byte           //委托类别
	chFunctionCode byte        //委托代码, B, S, C
	nOrderPrice int32           //委托价格((a double number + 0.00005) *10000)
	nOrderVolume int32         //委托数量
}

func combineNums(nDate int32, nTime int32) string {
	var str string
	str = strconv.Itoa(int(nDate)) + strconv.Itoa(int(nTime))
	return str
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

func array2str4int(arr [10]int32, len int) string {
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

func array2str4uint(arr [10]uint32, len int) string {
	var str string
	for i:=0; i<len; i++ {
		if i==len-1 {
			str += strconv.FormatUint(uint64(arr[i]), 10) + " "
		}else {
			str += strconv.FormatUint(uint64(arr[i]), 10) + ","
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
func GetKData(hTdb C.THANDLE, szCode string, szMarket string, nBeginDate int, nEndDate int, nCycle int, nUserDef int, nCQFlag int, nAutoComplete int, clnt client.Client) {
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
		windCode := Char2byte(uintptr(unsafe.Pointer(&kL.chWindCode)),unsafe.Sizeof(kL.chWindCode[0]),len(kL.chWindCode))
		code := Char2byte(uintptr(unsafe.Pointer(&kL.chCode)),unsafe.Sizeof(kL.chCode[0]),len(kL.chCode))
		fmt.Printf("WindCode:%s\n Code:%s\n Date:%d\n Time:%d\n Open:%d\n High:%d\n Low:%d\n Close:%v\n Volume:%v\n Turover:%d\n MatchItem:%d\n Interest:%d\n",
			windCode,//kL.chWindCode
			code,//kL.chCode
			kL.nDate, kL.nTime, kL.nOpen, kL.nHigh, kL.nLow, kL.nClose, kL.iVolume, kL.iTurover, kL.nMatchItems, kL.nInterest )
		fmt.Println("--------------------------------------")
		tmpPtr += sizeOf*100
		i += 100
	}
}

//tested good

func GetTickData(hTdb C.THANDLE, szCode string, szMarket string, nDate int, clnt client.Client)  {
	var req C.TDBDefine_ReqTick
	String2char(szCode,uintptr(unsafe.Pointer(&req.chCode)),unsafe.Sizeof(req.chCode[0]))
	String2char(szMarket,uintptr(unsafe.Pointer(&req.chMarketKey)),unsafe.Sizeof(req.chMarketKey[0]))

	req.nDate = C.int(nDate)
	req.nBeginTime = 0
	req.nEndTime = 0

	var pTick *C.TDBDefine_Tick = nil
	var pCount C.int
	C.TDB_GetTick(hTdb,&req,&pTick, &pCount)

	var tick Define_Tick
	fmt.Println("------------------------Tick Data---------------------------")
	fmt.Printf("共收到 %d 条Tick数据， 打印 1/100 条：\n", pCount)

	bp, _ := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  "TDB",
		Precision: "us",
	})

	tmpPtr := uintptr(unsafe.Pointer(pTick))
	sizeOf := unsafe.Sizeof(*pTick)
	for i:=0; i<int(pCount); {
		pT := (*C.TDBDefine_Tick)(unsafe.Pointer(tmpPtr))
		buf := (*[1024]byte)(unsafe.Pointer(pT))
		binary.Read(bytes.NewBuffer(buf[0:32]), binary.LittleEndian, &tick.chWindCode)
		binary.Read(bytes.NewBuffer(buf[32:64]), binary.LittleEndian, &tick.chCode)
		binary.Read(bytes.NewBuffer(buf[64:68]), binary.LittleEndian, &tick.nDate)
		binary.Read(bytes.NewBuffer(buf[68:72]), binary.LittleEndian, &tick.nTime)
		binary.Read(bytes.NewBuffer(buf[72:76]), binary.LittleEndian, &tick.nPrice)
		binary.Read(bytes.NewBuffer(buf[76:84]), binary.LittleEndian, &tick.iVolume)
		binary.Read(bytes.NewBuffer(buf[84:92]), binary.LittleEndian, &tick.iTurover)
		binary.Read(bytes.NewBuffer(buf[92:96]), binary.LittleEndian, &tick.nMatchItems)
		binary.Read(bytes.NewBuffer(buf[96:100]), binary.LittleEndian, &tick.nInterest)
		binary.Read(bytes.NewBuffer(buf[100:101]), binary.LittleEndian, &tick.chTradeFlag)
		binary.Read(bytes.NewBuffer(buf[101:102]), binary.LittleEndian, &tick.chBSFlag)
		binary.Read(bytes.NewBuffer(buf[102:110]), binary.LittleEndian, &tick.iAccVolume)
		binary.Read(bytes.NewBuffer(buf[110:118]), binary.LittleEndian, &tick.iAccTurover)
		binary.Read(bytes.NewBuffer(buf[118:122]), binary.LittleEndian, &tick.nHigh)
		binary.Read(bytes.NewBuffer(buf[122:126]), binary.LittleEndian, &tick.nLow)
		binary.Read(bytes.NewBuffer(buf[126:130]), binary.LittleEndian, &tick.nOpen)
		binary.Read(bytes.NewBuffer(buf[130:134]), binary.LittleEndian, &tick.nPreClose)
		binary.Read(bytes.NewBuffer(buf[134:138]), binary.LittleEndian, &tick.nSettle)
		binary.Read(bytes.NewBuffer(buf[138:142]), binary.LittleEndian, &tick.nPosition)
		binary.Read(bytes.NewBuffer(buf[142:146]), binary.LittleEndian, &tick.nCurDelta)
		binary.Read(bytes.NewBuffer(buf[146:150]), binary.LittleEndian, &tick.nPreSettle)
		binary.Read(bytes.NewBuffer(buf[150:154]), binary.LittleEndian, &tick.nPrePosition)
		binary.Read(bytes.NewBuffer(buf[154:194]), binary.LittleEndian, &tick.nAskPrice)
		binary.Read(bytes.NewBuffer(buf[194:234]), binary.LittleEndian, &tick.nAskVolume)
		binary.Read(bytes.NewBuffer(buf[234:274]), binary.LittleEndian, &tick.nBidPrice)
		binary.Read(bytes.NewBuffer(buf[274:314]), binary.LittleEndian, &tick.nBidVolume)
		binary.Read(bytes.NewBuffer(buf[314:318]), binary.LittleEndian, &tick.nAskAvPrice)
		binary.Read(bytes.NewBuffer(buf[318:322]), binary.LittleEndian, &tick.nBidAvPrice)
		binary.Read(bytes.NewBuffer(buf[322:330]), binary.LittleEndian, &tick.iTotalAskVolume)
		binary.Read(bytes.NewBuffer(buf[330:338]), binary.LittleEndian, &tick.iTotalBidVolume)
		binary.Read(bytes.NewBuffer(buf[338:342]), binary.LittleEndian, &tick.nIndex)
		binary.Read(bytes.NewBuffer(buf[342:346]), binary.LittleEndian, &tick.nStocks)
		binary.Read(bytes.NewBuffer(buf[346:350]), binary.LittleEndian, &tick.nUps)
		binary.Read(bytes.NewBuffer(buf[350:354]), binary.LittleEndian, &tick.nDowns)
		binary.Read(bytes.NewBuffer(buf[354:358]), binary.LittleEndian, &tick.nHoldLines)
		binary.Read(bytes.NewBuffer(buf[358:362]), binary.LittleEndian, &tick.nResv1)
		binary.Read(bytes.NewBuffer(buf[362:366]), binary.LittleEndian, &tick.nResv2)
		binary.Read(bytes.NewBuffer(buf[366:370]), binary.LittleEndian, &tick.nResv3)

		fmt.Printf("万得代码 chWindCode:%s \n", tick.chWindCode)
		fmt.Printf("日期 nDate:%d \n", tick.nDate)
		fmt.Printf("时间 nTime:%d \n", tick.nTime)

		fmt.Printf("成交价 nPrice:%d \n", tick.nPrice)
		fmt.Printf("成交量 iVolume:%d \n", tick.iVolume)
		fmt.Printf("成交额(元) iTurover:%d \n", tick.iTurover)
		fmt.Printf("成交笔数 nMatchItems:%d \n", tick.nMatchItems)
		fmt.Printf(" nInterest:%d \n", tick.nInterest)

		fmt.Printf("成交标志: chTradeFlag:%c \n", tick.chTradeFlag)
		fmt.Printf("BS标志: chBSFlag:%c \n", tick.chBSFlag)
		fmt.Printf("当日成交量: iAccVolume:%d \n", tick.iAccVolume)
		fmt.Printf("当日成交额: iAccTurover:%v \n", tick.iAccTurover)

		fmt.Printf("最高 nHigh:%d \n", tick.nHigh)
		fmt.Printf("最低 nLow:%d \n", tick.nLow)
		fmt.Printf("开盘 nOpen:%d \n",tick.nOpen)
		fmt.Printf("前收盘 nPreClose:%d \n", tick.nPreClose)

		//买卖盘字段
		var strOut string
		strOut = array2str4int(tick.nAskPrice, 10)
		fmt.Printf("叫卖价 nAskPrice:%s \n", strOut)
		strOut = array2str4uint(tick.nAskVolume, 10)
		fmt.Printf("叫卖量 nAskVolume:%s \n", strOut)
		strOut = array2str4int(tick.nBidPrice, 10)
		fmt.Printf("叫买价 nBidPrice:%s \n", strOut)
		strOut = array2str4uint(tick.nBidVolume, 10)
		fmt.Printf("叫买量 nBidVolume:%s \n", strOut)
		fmt.Printf("加权平均叫卖价 nAskAvPrice:%d \n", tick.nAskAvPrice)
		fmt.Printf("加权平均叫买价 nBidAvPrice:%d \n", tick.nBidAvPrice)
		fmt.Printf("叫卖总量 iTotalAskVolume:%v \n", tick.iTotalAskVolume)
		fmt.Printf("叫买总量 iTotalBidVolume:%v \n", tick.iTotalBidVolume)


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
		i += 1000
		tmpPtr += (sizeOf-2)*1000

		tags := map[string]string{
			"Code": string(tick.chCode[:]),
		}
		fields := map[string]interface{}{
			"Time": combineNums(tick.nDate, tick.nTime),
			"Price": tick.nPrice,
			"Volume": tick.iVolume,
			"Turover": tick.iTurover,
			"MatchItems": tick.nMatchItems,
			"Interest": tick.nInterest,
			"TradeFlag": tick.chTradeFlag,
			"BSFlag": tick.chBSFlag,
			"AccVolume": tick.iAccVolume,
			"AccTurover": tick.iAccTurover,
			"High": tick.nHigh,
			"Low": tick.nLow,
			"Open": tick.nOpen,
			"PreClose": tick.nPreClose,
			//期货字段
//			"Settle": tick.nSettle,
//			"Position": tick.nPosition,
//			"CurDelta": tick.nCurDelta,
//			"PreSettle": tick.nPreSettle,
//			"PrePosition": tick.nPrePosition,
			"AskPrice": array2str4int(tick.nAskPrice, 10),
			"AskVolume": array2str4uint(tick.nAskVolume, 10),
			"BidPrice": array2str4int(tick.nBidPrice, 10),
			"BidVolume": array2str4uint(tick.nBidVolume, 10),
			"AskAvPrice": tick.nAskAvPrice,
			"BidAvPrice": tick.nBidAvPrice,
			"TotalAskVolume": tick.iTotalAskVolume,
			"TotalBidVolume": tick.iTotalBidVolume,

		}
		pt, err := client.NewPoint(
			"TDBTransaction",
			tags,
			fields,
			time.Now(),
		)
		if err != nil {
			log.Fatal(err)
		}
		bp.AddPoint(pt)
	}
	if err := clnt.Write(bp); err != nil {
		log.Fatal(err)
	}
}


//tested good
func GetTransaction(hTdb C.THANDLE, szCode string, szMarketKey string, nDate int, clnt client.Client)  {
	var req C.TDBDefine_ReqTransaction
	String2char(szCode,uintptr(unsafe.Pointer(&req.chCode)),unsafe.Sizeof(req.chCode[0]))
	String2char(szMarketKey,uintptr(unsafe.Pointer(&req.chMarketKey)),unsafe.Sizeof(req.chMarketKey[0]))
	req.nDate = C.int(nDate)
	req.nBeginTime = 0
	req.nEndTime = 0

	var pTransaction *C.TDBDefine_Transaction = nil
	var pCount C.int
	C.TDB_GetTransaction(hTdb,&req, &pTransaction, &pCount)

	var transaction Define_Transaction
	fmt.Println("-----------------------Transaction Data----------------------------")
	fmt.Printf("收到 %d 条逐笔成交消息，打印 1/10000 条\n", pCount)

	bp, _ := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  "TDB",
		Precision: "us",
	})

	tmpPtr := uintptr(unsafe.Pointer(pTransaction))
	sizeOf := unsafe.Sizeof(*pTransaction)
	for i:=0; i<int(pCount); {
		pT := (*C.TDBDefine_Transaction)(unsafe.Pointer(tmpPtr))
		buf := (*[1024]byte)(unsafe.Pointer(pT))
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
		fmt.Printf("成交时间(Date): %d \n", transaction.nDate)
		fmt.Printf("成交时间: %d \n", transaction.nTime)
		fmt.Printf("成交代码: %c \n", transaction.chFunctionCode)
		fmt.Printf("委托类别: %c \n", transaction.chOrderKind)
		fmt.Printf("BS标志: %c \n", transaction.chBSFlag)
		fmt.Printf("成交价格: %d \n", transaction.nTradePrice)
		fmt.Printf("成交数量: %d \n", transaction.nTradeVolume)
		fmt.Printf("叫卖序号: %d \n", transaction.nAskOrder)
		fmt.Printf("叫买序号: %d \n", transaction.nBidOrder)
		fmt.Println("---------------------------------------------")
		//fmt.Printf("成交编号: %d \n", pT.nBidOrder)
		i += 10000
		tmpPtr += (sizeOf-1)*10000

		tags := map[string]string{
			"Code": string(transaction.chCode[:]),
		}
		fields := map[string]interface{}{
			"Time": combineNums(transaction.nDate, transaction.nTime),
			"Index": transaction.nIndex,
			"FunctionCode": transaction.chFunctionCode,
			"OrderKind": transaction.chOrderKind,
			"BSFlag": transaction.chBSFlag,
			"TradePrice": transaction.nTradePrice,
			"TradeVolume": transaction.nTradeVolume,
			"AskOrder": transaction.nAskOrder,
			"BidOrder": transaction.nBidOrder,
		}
		pt, err := client.NewPoint(
			"TDBTransaction",
			tags,
			fields,
			time.Now(),
		)
		if err != nil {
			log.Fatal(err)
		}
		bp.AddPoint(pt)
	}
	if err := clnt.Write(bp); err != nil {
		log.Fatal(err)
	}

}

//tested good
func GetOrder(hTdb C.THANDLE, szCode string, szMarketKey string, nDate int, clnt client.Client)  {
	var req C.TDBDefine_ReqOrder
	String2char(szCode,uintptr(unsafe.Pointer(&req.chCode)),unsafe.Sizeof(req.chCode[0]))
	String2char(szMarketKey,uintptr(unsafe.Pointer(&req.chMarketKey)),unsafe.Sizeof(req.chMarketKey[0]))
	req.nDate = C.int(nDate)
	req.nBeginTime = 0
	req.nEndTime = 0

	var pOrder *C.TDBDefine_Order = nil
	var pCount C.int
	C.TDB_GetOrder(hTdb,&req, &pOrder, &pCount)

	var order Define_Order
	fmt.Println("-------------------------Order Data--------------------------")
	fmt.Printf("收到 %d 条逐笔委托消息，打印 1/10000 条\n", pCount)

	bp, _ := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  "TDB",
		Precision: "us",
	})

	tmpPtr := uintptr(unsafe.Pointer(pOrder))
	sizeOf := unsafe.Sizeof(*pOrder)
	for i:=0; i<int(pCount); {
		pO := (*C.TDBDefine_Order)(unsafe.Pointer(tmpPtr))
		buf := (*[1024]byte)(unsafe.Pointer(pO))
		binary.Read(bytes.NewBuffer(buf[0:32]), binary.LittleEndian, &order.chWindCode)
		binary.Read(bytes.NewBuffer(buf[32:64]), binary.LittleEndian, &order.chCode)
		binary.Read(bytes.NewBuffer(buf[64:68]), binary.LittleEndian, &order.nDate)
		binary.Read(bytes.NewBuffer(buf[68:72]), binary.LittleEndian, &order.nTime)
		binary.Read(bytes.NewBuffer(buf[72:76]), binary.LittleEndian, &order.nIndex)
		binary.Read(bytes.NewBuffer(buf[76:80]), binary.LittleEndian, &order.nOrder)
		binary.Read(bytes.NewBuffer(buf[80:81]), binary.LittleEndian, &order.chOrderKind)
		binary.Read(bytes.NewBuffer(buf[81:82]), binary.LittleEndian, &order.chFunctionCode)
		binary.Read(bytes.NewBuffer(buf[82:86]), binary.LittleEndian, &order.nOrderPrice)
		binary.Read(bytes.NewBuffer(buf[86:90]), binary.LittleEndian, &order.nOrderVolume)
		fmt.Printf("订单时间(Date): %d \n", order.nDate)
		fmt.Printf("委托时间(HHMMSSmmm): %d \n", order.nTime)
		fmt.Printf("委托编号Order: %d \n", order.nOrder)
		fmt.Printf("委托类别OrderKind: %c \n", order.chOrderKind)
		fmt.Printf("委托代码FunctionCode: %c \n", order.chFunctionCode)
		fmt.Printf("委托价格OrderPrice: %d \n", order.nOrderPrice)
		fmt.Printf("委托数量OrderVolume: %d \n", order.nOrderVolume)
		fmt.Println("---------------------------------------------")
		//fmt.Println(order)
		i += 10000
		tmpPtr += (sizeOf-2)*10000

		tags := map[string]string{
			"Code": string(order.chCode[:]),
		}
		fields := map[string]interface{}{
			"Time": combineNums(order.nDate, order.nTime),
			"Index": order.nIndex,
			"Order": order.nOrder,
			"OrderKind": order.chOrderKind,
			"FunctionCode": order.chFunctionCode,
			"OrderPrice": order.nOrderPrice,
			"OrderVolume": order.nOrderVolume,
		}
		pt, err := client.NewPoint(
			"TDBOrder",
			tags,
			fields,
			time.Now(),
		)
		if err != nil {
			log.Fatal(err)
		}
		bp.AddPoint(pt)
	}
	if err := clnt.Write(bp); err != nil {
		log.Fatal(err)
	}
}

func GetOrderQueue(hTdb C.THANDLE, szCode string, szMarketKey string, nDate int, clnt client.Client) {
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


//指标公式
func UseEZFFormula(hTdb C.THANDLE) {
	fmt.Println("-------------------UseEZFFormula-------------");
	//公式的编写，请参考<<TRANSEND-TS-M0001 易编公式函数表V1(2).0-20110822.pdf>>
	strName := "KDJ"
	strContent := "INPUT:N(9), M1(3,1,100,2), M2(3);RSV:=(CLOSE-LLV(LOW,N))/(HHV(HIGH,N)-LLV(LOW,N))*100;K:SMA(RSV,M1,1);D:SMA(K,M2,1);J:3*K-2*D;"

	//添加公式到服务器并编译，若不过，会有错误返回
	var addRes *C.TDBDefine_AddFormulaRes = new(C.TDBDefine_AddFormulaRes)
	nErr := C.TDB_AddFormula(hTdb, C.CString(strName), C.CString(strContent),addRes)
	/*chInfo := Char2byte(uintptr(unsafe.Pointer(&addRes.chInfo)),unsafe.Sizeof(addRes.chInfo[0]),len(addRes.chInfo))
	fmt.Printf("Add Formula Result:%s\n",chInfo)

	string_chInfo := string(chInfo)
	enc := mahonia.NewEncoder("UTF-8")
	strr := enc.ConvertString(string_chInfo)
	fmt.Println(strr)*/
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

func queryDB(clnt client.Client, cmd string) (res []client.Result, err error) {
	q := client.Query{
		Command:  cmd,
		Database: MyDB,
	}
	if response, err := clnt.Query(q); err == nil {
		if response.Error() != nil {
			return res, response.Error()
		}
		res = response.Results
	} else {
		return res, err
	}
	return res, nil
}

func writePoints(clnt client.Client) {
	sampleSize := 10
	_, err := queryDB(clnt, fmt.Sprintf("DROP DATABASE %s", "systemstats"))
	if err != nil {
		log.Fatal(err)
	}

	_, err = queryDB(clnt, fmt.Sprintf("CREATE DATABASE %s", "systemstats"))
	if err != nil {
		log.Fatal(err)
	}

	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  "systemstats",
		Precision: "us",
	})
	if err != nil {
		log.Fatal(err)
	}

	rand.Seed(time.Now().UnixNano())
	for i := 0; i < sampleSize; i++ {
		regions := []string{"us-west1", "us-west2", "us-west3", "us-east1"}
		tags := map[string]string{
			"cpu":    "cpu-total",
			"host":   fmt.Sprintf("host%d", rand.Intn(1000)),
			"region": regions[rand.Intn(len(regions))],
		}

		idle := rand.Float64() * 100.0
		fields := map[string]interface{}{
			"idle": idle,
			"busy": 100.0 - idle,
		}

		pt, err := client.NewPoint(
			"cpu_usage",
			tags,
			fields,
			time.Now(),
		)
		if err != nil {
			log.Fatal(err)
		}
		bp.AddPoint(pt)
	}

	if err := clnt.Write(bp); err != nil {
		log.Fatal(err)
	}
}

