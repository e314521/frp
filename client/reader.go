package client
/*
#include <stdlib.h>
#include "SynReader.h"
#cgo LDFLAGS: -L ./libx64 -lSynReader64

const char* os = "windows";
const char* getchar(char a[] ){
	return os;
};
*/
import "C" 

import (
	"fmt"
	"unsafe"
	"os"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"github.com/fatedier/frp/pkg/util/log"
)
type User struct {
	CardType 	string		//I为外国人居住证，J 为港澳台居住证，空格(0x20)为普通身份证
	Name 		string		//姓名 
	EngName 	string		//英文名(外国人居住证)
	Sex 		string		//性别
	Nation		string		//民族或国籍(外国人居住证)     
	Birthday	string		//出生日期
	Address		string		//住址
	IDCardNo	string		//身份证号或外国人居住证号(外国人居住证)
	GrantDept	string		//发证机关
	UserLifeBegin	string		//有效开始日期
	UserLifeEnd	string		//有效截止日期
	PassID		string		//通行证号码(港澳台)
	IssuesTimes	string		//签发次数(港澳台)
	CertVol		string		//证件版本号(外国人居住证)
	HeadImage	string		//证件照
}

var nReader C.int

func Reader(w http.ResponseWriter, r *http.Request) {
	if(nReader != 0){
		nReader  = C.OpenUsbComm()
	}
	if(nReader != 0){
		log.Info("OpenUsbComm %d", nReader)
		w.Write([]byte(`{"success":false, "msg":"身份在阅读器未连接", "data":{}}`))
		return
	}
	
	var ret C.int
	var stIDCardDataUTF8 C.St_IDCardDataUTF8
	ret=C.getIDcard(&stIDCardDataUTF8)
	if(ret != 0){
		log.Info("getIDcard %d", ret)
		w.Write([]byte(`{"success":false, "msg":"请放入身份证", "data":{}}`))
		return
	}
	go_path := fmt.Sprintf("/oem/IDCard/%s.bmp", C.GoString((*C.char)(unsafe.Pointer(&stIDCardDataUTF8.IDCardNo))))
	path := C.CString(go_path)
	defer C.free(unsafe.Pointer(path))
	ret = C.saveWlt2Bmp((*C.char)(unsafe.Pointer(&stIDCardDataUTF8.wlt)),path)
	
	if(ret != 1){
		log.Info("Wlt2Bmp %d", ret)
		w.Write([]byte(`{"success":false, "msg":"身份照解密失败", "data":{}}`))
		return
	}
	bmp, err := os.ReadFile(go_path)
	if err != nil {
		w.Write([]byte(`{"success":false, "msg":"身份照读取失败", "data":{}}`))
	}
	headImage := base64.StdEncoding.EncodeToString(bmp)
	user := User{
		C.GoString((*C.char)(unsafe.Pointer(&stIDCardDataUTF8.CardType))), 
		C.GoString((*C.char)(unsafe.Pointer(&stIDCardDataUTF8.Name))),
		C.GoString((*C.char)(unsafe.Pointer(&stIDCardDataUTF8.EngName))),
		C.GoString((*C.char)(unsafe.Pointer(&stIDCardDataUTF8.Sex))),
		C.GoString((*C.char)(unsafe.Pointer(&stIDCardDataUTF8.Nation))),
		C.GoString((*C.char)(unsafe.Pointer(&stIDCardDataUTF8.Birthday))),
		C.GoString((*C.char)(unsafe.Pointer(&stIDCardDataUTF8.Address))),
		C.GoString((*C.char)(unsafe.Pointer(&stIDCardDataUTF8.IDCardNo))),
		C.GoString((*C.char)(unsafe.Pointer(&stIDCardDataUTF8.GrantDept))),
		C.GoString((*C.char)(unsafe.Pointer(&stIDCardDataUTF8.UserLifeBegin))),
		C.GoString((*C.char)(unsafe.Pointer(&stIDCardDataUTF8.UserLifeEnd))),
		C.GoString((*C.char)(unsafe.Pointer(&stIDCardDataUTF8.PassID))),
		C.GoString((*C.char)(unsafe.Pointer(&stIDCardDataUTF8.IssuesTimes))),
		C.GoString((*C.char)(unsafe.Pointer(&stIDCardDataUTF8.CertVol))),
		headImage,
		
	}
	data, err := json.Marshal(&user)
	if err != nil {
		log.Info("序列化错误 err=%v\n", err)
		w.Write([]byte(`{"success":false, "msg":"序列化失败", "data":{}}`))
		return
	}
	w.Write([]byte(`{"success":true, "msg":"获取成功", "data":`+ string(data) +`}`))
	
}

func (svr *Service)RunReaderServer(address string) (err error) {
	
	nReader  = C.OpenUsbComm()
	http.HandleFunc("/getIDcard", Reader)
	go http.ListenAndServe(address, nil)
	return

	var nRet C.int
	nRet  = C.OpenUsbComm()
	if(nRet != 0){
		log.Info("OpenUsbComm %d", nRet)
	}

	log.Info("admin server listen on %d", nRet)
	return
}



