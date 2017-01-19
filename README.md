#TDBcpp2go：解决TDB中接口用Go语言调用的问题

### go不能遍历c结构体数组的问题：
```go
tmpPtr := uintptr(unsafe.Pointer(pCode))//已有pCode*数组
sizeOf := unsafe.Sizeof(*pCode)
pC := (*C.TDBDefine_Code)(unsafe.Pointer(tmpPtr))

fmt.Println("-------------code table ----------------------------");
fmt.Printf("chWindCode:%s \n", pC.chCode);
fmt.Printf("chWindCode:%s \n", pC.chMarket);
fmt.Printf("chWindCode:%s \n", pC.chCNName);
fmt.Printf("chWindCode:%s \n", pC.chENName);
fmt.Printf("chWindCode:%s \n", pC.nType);

tmpPtr += sizeOf
```
解答来源于http://studygolang.com/topics/594

### go没有char[]数组，可以先将string转换为byte数组，然后再每个元素分别转换为C.char型
```go
settings_bytes1 := []byte("114.80.154.34")
	for i:=0; i<len(settings_bytes1); i++{
		settings.szIP[i]=C.char(settings_bytes1[i])
	}
```

### 在cgo中，如果c结构体整形、字符类型相间定义，由于字段对齐规则不同，无法对结构体所有字段直接赋值，可以通过切片的方式逐片拼装
```go
l := unsafe.Sizeof(*pTransaction)
buf := (*[1024]byte)(unsafe.Pointer(pTransaction))

var transaction Define_Transaction
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
```
事实上也可以不用这么麻烦，只要知道了下一个结构体的地址在当前结构体地址+size然后-1或者-2后，可以直接用unsafe.Pointer去访问输出。

### 数字到string的转换，有strconv包里的几个函数,其中strconv.Itoa(int)可以把一个int型数转换为string，strconv.FormatUint(uint64, base)可以把一个uint64型数转换成base进制后再变成string

