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
	"unsafe"
)

func main(){
	var hTdb C.THANDLE = nil


	var settings C.OPEN_SETTINGS

	//================================================
	settings_bytes1 := []byte("114.80.154.34")
	for i:=0; i<len(settings_bytes1); i++{
		settings.szIP[i]=C.char(settings_bytes1[i])
	}
	//================================================
	settings_bytes2 := []byte("6261")
	for i:=0 ;i<len(settings_bytes2) ;i++{
		settings.szPort[i]=C.char(settings_bytes2[i])
	}
	//================================================
	settings_bytes3 := []byte("TD3446699001")
	for i:=0 ;i<len(settings_bytes3) ;i++{
		settings.szUser[i]=C.char(settings_bytes3[i])
	}
	//================================================
	settings_bytes4 := []byte("43449360")
	for i:=0 ;i<len(settings_bytes4) ;i++{
		settings.szPassword[i]=C.char(settings_bytes4[i])
	}
	//================================================

	settings.nRetryCount = 15
	settings.nRetryGap = 1
	settings.nTimeOutVal = 1


	//proxy
/*
	var proxy_setting C.TDB_PROXY_SETTING
	proxy_setting.nProxyType = C.TDB_PROXY_HTTP11

	//================================================
	proxy_bytes1 := []byte("10.100.3.42")
	for i:=0 ;i<len(proxy_bytes1) ;i++{
		proxy_setting.szProxyHostIp[i]=C.char(proxy_bytes1[i])
	}
	//================================================
	proxy_bytes2 := []byte("12345")
	for i:=0 ;i<len(proxy_bytes2) ;i++{
		proxy_setting.szProxyPort[i]=C.char(proxy_bytes2[i])
	}
	//================================================
	proxy_bytes3 := []byte("1")
	for i:=0 ;i<len(proxy_bytes3) ;i++{
		proxy_setting.szProxyUser[i]=C.char(proxy_bytes3[i])
	}
	//================================================
	proxy_bytes4 := []byte("1")
	for i:=0 ;i<len(proxy_bytes4) ;i++{
		proxy_setting.szProxyPwd[i]=C.char(proxy_bytes4[i])
	}
	//================================================
	*/
	var LoginRes C.TDBDefine_ResLogin
	//TDB_OpenProxy
	//hTdb = C.TDB_OpenProxy(&settings, &proxy_setting, &LoginRes)

	hTdb = C.TDB_Open(&settings, &LoginRes)
	//fmt.Println("aaa")
	if hTdb == nil {
		fmt.Println("连接失败！")
		return
	}

	//TDB_GetCOdeInfo
	var pCode *C.TDBDefine_Code
	pCode = C.TDB_GetCodeInfo(hTdb, C.CString("000001.SZ"), C.CString("SZ-2-0"))
	var code_byte = [32]byte{}
	for i:=0; i<len(pCode.chCode); i++ {
		code_byte[i] = byte(pCode.chCode[i])
	}
	fmt.Printf("交易所代码 chWindCode:%s \n", code_byte)


	var pCount C.int = 0
	C.TDB_GetCodeTable(hTdb,C.CString("SZ"),&pCode,&pCount);
	tmpPtr := uintptr(unsafe.Pointer(pCode))
	sizeOf := unsafe.Sizeof(*pCode)
	if pCount!=0 && pCode!=nil{
		for i := 0; i < 2; i++{
		pC := (*C.TDBDefine_Code)(unsafe.Pointer(tmpPtr))
		fmt.Println("-------------code table ----------------------------");
		fmt.Printf("chWindCode:%s \n", pC.chCode);
		fmt.Printf("chWindCode:%s \n", pC.chMarket);
		fmt.Printf("chWindCode:%s \n", pC.chCNName);
		fmt.Printf("chWindCode:%s \n", pC.chENName);
		fmt.Printf("chWindCode:%s \n", pC.nType);
		tmpPtr += sizeOf
		}
	}
}
