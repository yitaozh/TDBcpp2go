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

/*
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
*/

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

func UseEZFFormula(hTdb C.THANDLE){
	//公式的编写，请参考<<TRANSEND-TS-M0001 易编公式函数表V1(2).0-20110822.pdf>>;
	strName := "KDJ"
	strContent := "INPUT:N(9), M1(3,1,100,2), M2(3);"
	"RSV:=(CLOSE-LLV(LOW,N))/(HHV(HIGH,N)-LLV(LOW,N))*100;"
	"K:SMA(RSV,M1,1);"
	"D:SMA(K,M2,1);"
	"J:3*K-2*D;"

	//添加公式到服务器并编译，若不过，会有错误返回
	TDBDefine_AddFormulaRes* addRes = new TDBDefine_AddFormulaRes
	nErr := C.TDB_AddFormula(hTdb, C.CString(strName), C.CString(strContent) ,addRes)
	fmt.Printf("Add Formula Result:%s",addRes.chInfo)

	//查询服务器上的公式，能看到我们刚才上传的"KDJ"
	var pEZFItem *C.TDBDefine_FormulaItem = nil
	nItems := 0
	//名字为空表示查询服务器上所有的公式
	nErr = C.TDB_GetFormula(hTdb, C.NULL, &pEZFItem, &nItems);

	for i:=0; i<nItems; i++{
		std::string strNameInner(pEZFItem[i].chFormulaName, 0, sizeof(pEZFItem[i].chFormulaName))
		std::string strParam(pEZFItem[i].chParam, 0, sizeof(pEZFItem[i].chParam))
		printf("公式名称：%s, 参数:%s \n", strNameInner.c_str(), strParam.c_str())
	}

	type EZFCycDefine struct{
		char chName[8]
		int  nCyc
		int  nCyc1
	}
	EZFCyc[5]={
		{"日线", 2, 0},
		{"30分", 0, 30},
		{"5分钟", 0, 5},
		{"1分钟", 0, 1},
		{"15秒", 11, 15}
	}

	//获取公式的计算结果
	TDBDefine_ReqCalcFormula reqCalc = {0}
	strncpy(reqCalc.chFormulaName, "KDJ", sizeof(reqCalc.chFormulaName))
	strncpy(reqCalc.chParam, "N=9,M1=3,M2=3", sizeof(reqCalc.chParam))
	strncpy(reqCalc.chCode, "000001.SZ", sizeof(reqCalc.chCode))
	strncpy(reqCalc.chMarketKey, "SZ-2-0", sizeof(reqCalc.chMarketKey))
	reqCalc.nCycType = (CYCTYPE)(EZFCyc[0].nCyc) 	//0表示日线
	reqCalc.nCycDef = EZFCyc[0].nCyc1
	reqCalc.nCQFlag = (REFILLFLAG)0			//除权标志
	reqCalc.nCalcMaxItems = 4000 			//计算的最大数据量
	reqCalc.nResultMaxItems = 100			//传送的结果的最大数据量

	TDBDefine_CalcFormulaRes* pResult = new TDBDefine_CalcFormulaRes
	nErr = TDB_CalcFormula(hTdb, &reqCalc, pResult)
	//判断错误代码

	printf("计算结果有: %d 条:\n", pResult->nRecordCount)
	char szLineBuf[1024] = {0}
	//输出字段名
	for j:=0; j<pResult->nFieldCount;j++{
		std::cout << pResult->chFieldName[j] << "  "
	}
	std::cout << endl << endl
	//输出数据
	for (int i=0; i<pResult->nRecordCount; i+=100){
		for (int j=0; j<pResult->nFieldCount;j++){
			std::cout << (pResult->dataFileds)[j][i] << "  "}
		std::cout << endl
	}

	//删除之前上传的公式指标
	TDBDefine_DelFormulaRes pDel = {0}
	nErr = TDB_DeleteFormula(hTdb, "KDJ", &pDel)
	printf("删除指标信息:%s", pDel.chInfo)
	//释放内存
	delete pEZFItem
	TDB_ReleaseCalcFormula(pResult)
}

