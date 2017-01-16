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
	//"time"
	"unsafe"
	//"strconv"

)


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
/*
func cpy(str string, des uintptr, sizeOf uintptr, len int) {
	bytes := []byte(str)
	for i:=0; i<len; i++ {
		unit := (*C.char)(unsafe.Pointer(des))
		*unit = C.char(bytes[i])
		des += sizeOf
	}
}
*/


/*func GetTickCount() int64 {
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
}*/

/*
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

func GetOrderQueue(hTdb C.THANDLE, szCode *C.char, szMarketKey *C.char, nDate int){
	//请求
	var req C.TDBDefine_ReqOrderQueue = {0}
	strncpy(req.chCode, szCode, sizeof(req.chCode)) //代码写成你想获取的股票代码
	strncpy(req.chMarketKey, szMarketKey, sizeof(req.chMarketKey))
	req.nDate = nDate
	req.nBeginTime = 0
	req.nEndTime = 0

	var pOrderQueue *C.TDBDefine_OrderQueue = nil
	var pCount int
	C.TDB_GetOrderQueue(hTdb,&req, &pOrderQueue, &pCount)

	fmt.Println("-------------------OrderQueue Data-------------")
	fmt.Printf("收到 %d 条委托队列消息，打印 1/1000 条\n", pCount);

	for i:=0; i<pCount ;{
		const TDBDefine_OrderQueue& que = pOrderQueue[i];
		printf("订单时间(Date): %d \n", que.nDate);
		printf("订单时间(HHMMSS): %d \n", que.nTime);
		printf("买卖方向('B':Bid 'A':Ask): %c \n", que.nSide);
		printf("成交价格: %d \n", que.nPrice);
		printf("订单数量: %d \n", que.nOrderItems);
		printf("明细个数: %d \n", que.nABItems);
		printf("订单明细: %s \n", array2str(que.nABVolume, que.nABItems).c_str());
		printf("-------------------------------\n");
		i += 1000;
	}
	//释放
	C.TDB_Free(pOrderQueue);
}
*/
func UseEZFFormula(hTdb C.THANDLE){
	//公式的编写，请参考<<TRANSEND-TS-M0001 易编公式函数表V1(2).0-20110822.pdf>>;
	strName := "KDJ"
	strContent := "INPUT:N(9), M1(3,1,100,2), M2(3);RSV:=(CLOSE-LLV(LOW,N))/(HHV(HIGH,N)-LLV(LOW,N))*100;K:SMA(RSV,M1,1);D:SMA(K,M2,1);J:3*K-2*D;"
	var addRes *C.TDBDefine_AddFormulaRes = new(C.TDBDefine_AddFormulaRes)

	nErr := C.TDB_AddFormula(hTdb, C.CString(strName), C.CString(strContent), addRes)

	fmt.Printf("Add Formula Result:%s\n",Char2byte(uintptr(unsafe.Pointer(&addRes.chInfo)),unsafe.Sizeof(addRes.chInfo[0]),len(addRes.chInfo)))

	//查询服务器上的公式，能看到我们刚才上传的"KDJ"
	var pEZFItem *C.TDBDefine_FormulaItem = nil
	var nItems C.int = 0

	//名字为空表示查询服务器上所有的公式
	nErr = C.TDB_GetFormula(hTdb, nil, &pEZFItem, &nItems)
	tmpPtr_pEZFItem := uintptr(unsafe.Pointer(pEZFItem))
	sizeOf_pEZFItem := unsafe.Sizeof(*pEZFItem)
	//fmt.Println(nItems)
	var i C.int = 0
	var j C.int = 0
	for i=0; i<nItems; i++{
		tmpEZFItem := (*C.TDBDefine_FormulaItem)(unsafe.Pointer(tmpPtr_pEZFItem))
		strNameInner := Char2byte(uintptr(unsafe.Pointer(&tmpEZFItem.chFormulaName)),unsafe.Sizeof(tmpEZFItem.chFormulaName),len(tmpEZFItem.chFormulaName))
		strParam := Char2byte(uintptr(unsafe.Pointer(&tmpEZFItem.chParam)),unsafe.Sizeof(tmpEZFItem.chParam),len(tmpEZFItem.chParam))
		fmt.Printf("公式名称：%s, 参数:%s \n", strNameInner, strParam)
		tmpPtr_pEZFItem += sizeOf_pEZFItem
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
	//strncpy(reqCalc.chFormulaName, "KDJ", sizeof(reqCalc.chFormulaName));
	tmpPtr_chParam := uintptr(unsafe.Pointer(&reqCalc.chParam))
	sizeOf_chParam := unsafe.Sizeof(reqCalc.chParam)
	String2char("N=9,M1=3,M2=3", tmpPtr_chParam, sizeOf_chParam)
	//strncpy(reqCalc.chParam, "N=9,M1=3,M2=3", sizeof(reqCalc.chParam))
	tmpPtr_chCode := uintptr(unsafe.Pointer(&reqCalc.chCode))
	sizeOf_chCode := unsafe.Sizeof(reqCalc.chCode)
	String2char("000001.SZ", tmpPtr_chCode, sizeOf_chCode)
	//strncpy(reqCalc.chCode, "000001.SZ", sizeof(reqCalc.chCode))
	tmpPtr_chMarketKey := uintptr(unsafe.Pointer(&reqCalc.chMarketKey))
	sizeOf_chMarketKey := unsafe.Sizeof(reqCalc.chMarketKey)
	String2char("SZ-2-0", tmpPtr_chMarketKey, sizeOf_chMarketKey)
	//strncpy(reqCalc.chMarketKey, "SZ-2-0", sizeof(reqCalc.chMarketKey))

	reqCalc.nCycType = uint32(EZFCyc[0].nCyc)		//0表示日线
	reqCalc.nCycDef = C.int(EZFCyc[0].nCyc1)
	reqCalc.nCQFlag = 0				//除权标志
	reqCalc.nCalcMaxItems = 4000 			//计算的最大数据量
	reqCalc.nResultMaxItems = 100			//传送的结果的最大数据量

	var pResult *C.TDBDefine_CalcFormulaRes = new(C.TDBDefine_CalcFormulaRes)
	nErr = C.TDB_CalcFormula(hTdb, &reqCalc, pResult)
	//判断错误代码

	fmt.Printf("计算结果有: %d 条:\n", pResult.nRecordCount)
	//var szLineBuf[1024]byte
	//输出字段名

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
	//释放内存
	//C.delete(pEZFItem)
	C.TDB_ReleaseCalcFormula(pResult)
}