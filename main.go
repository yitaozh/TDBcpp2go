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
	"log"
	"fmt"
	"unsafe"
	"github.com/influxdata/influxdb/client/v2"
)


func main(){
	var hTdb C.THANDLE = nil

	var settings C.OPEN_SETTINGS

	//================================================
	String2char("114.80.154.34",uintptr(unsafe.Pointer(&settings.szIP)),unsafe.Sizeof(settings.szIP[0]))
	String2char("6261",uintptr(unsafe.Pointer(&settings.szPort)),unsafe.Sizeof(settings.szPort[0]))
	String2char("TD3446699001",uintptr(unsafe.Pointer(&settings.szUser)),unsafe.Sizeof(settings.szUser[0]))
	String2char("43449360",uintptr(unsafe.Pointer(&settings.szPassword)),unsafe.Sizeof(settings.szPassword[0]))
	//================================================
	settings.nRetryCount = 100
	settings.nRetryGap = 100
	settings.nTimeOutVal = 100

	//proxy
/*	var proxy_setting C.TDB_PROXY_SETTING

	proxy_setting.nProxyType = C.TDB_PROXY_HTTP11
	//================================================
	string2char("10.100.3.42",uintptr(unsafe.Pointer(&proxy_setting.szProxyHostIp)),unsafe.Sizeof(proxy_setting.szProxyHostIp[0]))
	string2char("12345",uintptr(unsafe.Pointer(&proxy_setting.szProxyPort)),unsafe.Sizeof(proxy_setting.szProxyPort[0]))
	string2char("1",uintptr(unsafe.Pointer(&proxy_setting.szProxyUser)),unsafe.Sizeof(proxy_setting.szProxyUser[0]))
	string2char("1",uintptr(unsafe.Pointer(&proxy_setting.szProxyPwd)),unsafe.Sizeof(proxy_setting.szProxyPwd[0]))
	//================================================
	*/

	var LoginRes C.TDBDefine_ResLogin
	//TDB_OpenProxy
	//hTdb = C.TDB_OpenProxy(&settings, &proxy_setting, &LoginRes)

	hTdb = C.TDB_Open(&settings, &LoginRes)
	if hTdb == nil {
		fmt.Println("连接失败！")
		return
	}

	//TDB_GetCOdeInfo
	var pCode *C.TDBDefine_Code
	pCode = C.TDB_GetCodeInfo(hTdb, C.CString("000001.SZ"), C.CString("SZ-2-0"))
	fmt.Printf("交易所代码 chWindCode:%s \n", Char2byte(uintptr(unsafe.Pointer(&pCode.chCode)),unsafe.Sizeof(pCode.chCode[0]),len(pCode.chCode)))

	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     "http://114.80.253.159:8086",
		//Addr:     "http://localhost:8086",
		Username: username,
		Password: password,
	})
	if err != nil {
		log.Fatal(err)
	}
	_, err = queryDB(c, fmt.Sprintf("DROP DATABASE %s", "TDB"))
	if err != nil {
		log.Fatal(err)
	}
	_, err = queryDB(c, fmt.Sprintf("CREATE DATABASE %s", "TDB"))
	if err != nil {
		log.Fatal(err)
	}

	//GetKData(hTdb, "000001.SZ", "SZ-2-0", 20170208, 20170208, C.CYC_MINUTE, 0, 0, 1,c);	//autocomplete k-minute
	//GetTickData(hTdb, "000001.SZ", "SZ-2-0", 20161122,c);//带买卖盘的tick					//tick
	//GetTransaction(hTdb, "000001.sz", "SZ-2-0", 20150910, c);					//Transaction
	//GetOrder(hTdb, "112436.SZ", "SZ-2-0", 20170203, c);					//Order
	//GetOrderQueue(hTdb, "000001.sz", "SZ-2-0", 20150910,c);					//OrderQueue
	//UseEZFFormula(hTdb);									//test for formula
	Table, count := GetCodeTable(hTdb, "SZ-2-0")
	for i := 20170203; i<=20170228; i++ {
		for j := 0; j < count; j++ {
			GetKData(hTdb, Table[j].chWindCode, Table[j].chMarket, i, i, C.CYC_MINUTE, 0, 0, 1,c);			//autocomplete k-minute
			GetTickData(hTdb, Table[j].chWindCode, Table[j].chMarket, i,c);//带买卖盘的tick				//tick
			GetTransaction(hTdb, Table[j].chWindCode, Table[j].chMarket, i, c);					//Transaction
			GetOrder(hTdb, Table[j].chWindCode, Table[j].chMarket, i, c);							//Order
			GetOrderQueue(hTdb, Table[j].chWindCode, Table[j].chMarket, i,c);					//OrderQueue
			//UseEZFFormula(hTdb);
			if j%10 == 0{
				fmt.Println(j);
			}
		}
	}
	/*for j := 0; j < count; j=j+10 {
		fmt.Println(Table[j].chWindCode,"\t",Table[j].codeType)
	}*/
	fmt.Println(count)
}
