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
	settings_bytes1 := []byte("114.80.154.34")
	for i:=0 ;i<len(settings_bytes1) ;i++{
		settings.szIP[i]=C.char(settings_bytes1[i])
	}
	//================================================
	settings_bytes2 := []byte("6271")
	for i:=0 ;i<len(settings_bytes2) ;i++{
		settings.szPort[i]=C.char(settings_bytes2[i])
	}
	//================================================
	settings_bytes3 := []byte("TD3446699201")
	for i:=0 ;i<len(settings_bytes3) ;i++{
		settings.szUser[i]=C.char(settings_bytes3[i])
	}
	//================================================
	settings_bytes4 := []byte("43449560")
	for i:=0 ;i<len(settings_bytes4) ;i++{
		settings.szPassword[i]=C.char(settings_bytes4[i])
	}
	//================================================

	settings.nRetryCount = 15
	settings.nRetryGap = 1
	settings.nTimeOutVal = 1


	//proxy

	var proxy_setting C.TDB_PROXY_SETTING
	proxy_setting.nProxyType = C.TDB_PROXY_HTTP11
	fmt.Println("aaa")
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
	fmt.Println("aaa")
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
	fmt.Printf("交易所代码 chWindCode:%s \n", pCode.chCode)
}
