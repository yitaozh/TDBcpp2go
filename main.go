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
)

func main(){
	var hTdb C.THANDLE = nil

	var settings C.OPEN_SETTINGS

	//================================================
	bytes1 := []byte("114.80.154.34")
	for i:=0 ;i<len(bytes1) ;i++{
		settings.szIP[i]=C.char(bytes1[i])
	}
	//================================================
	bytes2 := []byte("6271")
	for i:=0 ;i<len(bytes2) ;i++{
		settings.szPort[i]=C.char(bytes2[i])
	}
	//================================================
	bytes3 := []byte("TD3446699201")
	for i:=0 ;i<len(bytes3) ;i++{
		settings.szUser[i]=C.char(bytes3[i])
	}
	//================================================
	bytes4 := []byte("43449560")
	for i:=0 ;i<len(bytes4) ;i++{
		settings.szPassword[i]=C.char(bytes4[i])
	}
	//================================================

	settings.nRetryCount = 15
	settings.nRetryGap = 1
	settings.nTimeOutVal = 1

	//proxy

	/*var proxy_setting C.TDB_PROXY_SETTING
	proxy_setting.nProxyType = TDB_PROXY_TYPE.TDB_PROXY_HTTP11
	C.strcpy(proxy_setting.szProxyHostIp, C.CString("10.100.3.42"))
	C.strcpy(proxy_setting.szProxyPort, C.CString("12345"))
	C.strcpy(proxy_setting.szProxyUser, C.CString("1"))
	C.strcpy(proxy_setting.szProxyPwd, C.CString("1"))
	*/
	var LoginRes C.TDBDefine_ResLogin
	//TDB_OpenProxy
	//hTdb = TDB_OpenProxy(&settings, &proxy_settings, &LoginRes)

	hTdb = C.TDB_Open(&settings, &LoginRes)

	if hTdb == nil {
		fmt.Println("连接失败！")
	}

	//TDB_GetCOdeInfo
	var pCode *C.TDBDefine_Code
	pCode = C.TDB_GetCodeInfo(hTdb, C.CString("000001.SZ"), C.CString("SZ-2-0"))
	fmt.Printf("交易所代码 chWindCode:%s \n", pCode.chCode)
}
