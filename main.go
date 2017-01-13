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

func string2char(str string, des uintptr, sizeOf uintptr){
	bytes := []byte(str)
	for i:=0; i<len(bytes); i++{
		unit := (*C.char)(unsafe.Pointer(des))
		*unit = C.char(bytes[i])
		des += sizeOf
	}
}

func char2byte(src uintptr, sizeOf uintptr) [128]byte {
	var des [128]byte
	for i:=0; i<32;i++ {
		unit := (*byte)(unsafe.Pointer(src))
		des[i] = *unit
		src += sizeOf
	}
	return des
}

func main(){
	var hTdb C.THANDLE = nil

	var settings C.OPEN_SETTINGS
	//================================================
	string2char("114.80.154.34",uintptr(unsafe.Pointer(&settings.szIP)),unsafe.Sizeof(settings.szIP[0]))
	string2char("6261",uintptr(unsafe.Pointer(&settings.szPort)),unsafe.Sizeof(settings.szPort[0]))
	string2char("TD3446699001",uintptr(unsafe.Pointer(&settings.szUser)),unsafe.Sizeof(settings.szUser[0]))
	string2char("43449360",uintptr(unsafe.Pointer(&settings.szPassword)),unsafe.Sizeof(settings.szPassword[0]))
	//================================================
	settings.nRetryCount = 15
	settings.nRetryGap = 1
	settings.nTimeOutVal = 1

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
	/*var code_byte = [32]byte{}
	for i:=0; i<len(pCode.chCode); i++ {
		code_byte[i] = byte(pCode.chCode[i])
	}*/

	fmt.Printf("交易所代码 chWindCode:%s \n", char2byte(uintptr(unsafe.Pointer(&pCode.chCode)),unsafe.Sizeof(pCode.chCode[0])))

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
