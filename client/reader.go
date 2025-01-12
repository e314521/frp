package client

/*
#include <stdlib.h>
#include <stdio.h>
#include "SynReader.h"

#cgo amd64 LDFLAGS: -L ./libx64  -lSynReader64  -lwlt -lusb-1.0
#cgo arm LDFLAGS: -L ./libArm  -lSynReaderArm -lwlt -lusb-1.0

*/
import "C"

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/fatedier/frp/pkg/util/log"
	"net/http"
	"os"
	"sync"
	"unsafe"
	"net/http/httputil"
	"net/url"
	"github.com/guanguans/id-validator"
)

var mu sync.RWMutex

type User struct {
	CardType      string `json:"card_type"`       //I为外国人居住证，J 为港澳台居住证，空格(0x20)为普通身份证
	Name          string `json:"name"`            //姓名
	EngName       string `json:"eng_name"`        //英文名(外国人居住证)
	Sex           string `json:"sex"`             //性别
	Nation        string `json:"nation"`          //民族或国籍(外国人居住证)
	Birthday      string `json:"birthday"`        //出生日期
	Address       string `json:"address"`         //住址
	IDCardNo      string `json:"id_card_no"`      //身份证号或外国人居住证号(外国人居住证)
	GrantDept     string `json:"grant_dept"`      //发证机关
	UserLifeBegin string `json:"user_life_begin"` //有效开始日期
	UserLifeEnd   string `json:"user_life_end"`   //有效截止日期
	HeadImage     string `json:"head_image"`      //证件照
	BirthPlace    string `json:"birth_place"`      //籍贯
}
func GetDate(date string) string{
	if(len(date) != 8){
		return date
	}
	dateRune := []rune(date)
	return string(dateRune[:4]) + "-" + string(dateRune[4:6]) + "-" + string(dateRune[6:8])

}
func Reader(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Access-Control-Allow-Methods", "*")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	mu.Lock()
	defer mu.Unlock()
	nReader := C.OpenUsbComm()
	defer C.CloseComm()
	if nReader != 0 {
		log.Info("OpenUsbComm %d", nReader)
		w.Write([]byte(`{"success":false, "msg":"身份证阅读器未连接", "data":{}}`))
		return
	}

	var ret C.int
	var stIDCardDataUTF8 C.St_IDCardDataUTF8
	ret = C.getIDcard(&stIDCardDataUTF8)
	if ret != 0 {
		log.Info("getIDcard %d", ret)
		w.Write([]byte(`{"success":false, "msg":"请放入身份证", "data":{}}`))
		return
	}
	go_path := fmt.Sprintf("/oem/IDCard/%s.bmp", C.GoString((*C.char)(unsafe.Pointer(&stIDCardDataUTF8.IDCardNo))))
	path := C.CString(go_path)
	defer C.free(unsafe.Pointer(path))
	ret = C.saveWlt2Bmp((*C.char)(unsafe.Pointer(&stIDCardDataUTF8.wlt)), path)

	if ret != 1 {
		log.Info("Wlt2Bmp %d", ret)
		w.Write([]byte(`{"success":false, "msg":"身份照解密失败", "data":{}}`))
		return
	}
	bmp, err := os.ReadFile(go_path)
	if err != nil {
		w.Write([]byte(`{"success":false, "msg":"身份照读取失败", "data":{}}`))
		return
	}
	if C.GoString((*C.char)(unsafe.Pointer(&stIDCardDataUTF8.CardType))) != " " {
		w.Write([]byte(`{"success":false, "msg":"当前只支持大陆身份证", "data":{}}`))
		return
	}
	headImage := base64.StdEncoding.EncodeToString(bmp)
	user := User{
		C.GoString((*C.char)(unsafe.Pointer(&stIDCardDataUTF8.CardType))),
		C.GoString((*C.char)(unsafe.Pointer(&stIDCardDataUTF8.Name))),
		C.GoString((*C.char)(unsafe.Pointer(&stIDCardDataUTF8.EngName))),
		C.GoString((*C.char)(unsafe.Pointer(&stIDCardDataUTF8.Sex))),
		C.GoString((*C.char)(unsafe.Pointer(&stIDCardDataUTF8.Nation))),
		GetDate(C.GoString((*C.char)(unsafe.Pointer(&stIDCardDataUTF8.Birthday)))),
		C.GoString((*C.char)(unsafe.Pointer(&stIDCardDataUTF8.Address))),
		C.GoString((*C.char)(unsafe.Pointer(&stIDCardDataUTF8.IDCardNo))),
		C.GoString((*C.char)(unsafe.Pointer(&stIDCardDataUTF8.GrantDept))),
		GetDate(C.GoString((*C.char)(unsafe.Pointer(&stIDCardDataUTF8.UserLifeBegin)))),
		GetDate(C.GoString((*C.char)(unsafe.Pointer(&stIDCardDataUTF8.UserLifeEnd)))),
		headImage,
		"",
	}
	IdInfo, err :=idvalidator.GetInfo(user.IDCardNo, false)
	if(err == nil){
		user.BirthPlace = IdInfo.AddressTree[0] + IdInfo.AddressTree[1]
	}
	data, err := json.Marshal(&user)
	if err != nil {
		log.Info("序列化错误 err=%v\n", err)
		w.Write([]byte(`{"success":false, "msg":"序列化失败", "data":{}}`))
		return
	}
	w.Write([]byte(`{"success":true, "msg":"获取成功", "data":` + string(data) + `}`))

}

func Reverse(w http.ResponseWriter, r *http.Request) {
	if r.RequestURI == "/favicon.ico" {
        	//io.WriteString(w, "Request path Error")
        	return
    	}
	remote, err := url.Parse("http://139.170.150.17:8080" )
	if err != nil {
		panic(err)
		return
	}
	proxy := httputil.NewSingleHostReverseProxy(remote)
	//w.Header().Set("Access-Control-Allow-Headers", "*")
	//w.Header().Set("Access-Control-Allow-Methods", "*")
	//w.Header().Set("Access-Control-Allow-Origin", "*")
	//w.Header().Set("Access-Control-Allow-Credentials", "true")
	proxy.ServeHTTP(w, r)
	w.Header().Add("Access-Control-Allow-Methods1", "*")
}

func (svr *Service) RunReaderServer(address string) (err error) {
	folderPath := "/oem/IDCard/"
	if _, err := os.Stat(folderPath); os.IsNotExist(err) {
		os.Mkdir(folderPath, 0777)
	}
	http.HandleFunc("/", Reverse)
	http.HandleFunc("/getIDcard", Reader)
	go http.ListenAndServe(address, nil)
	go http.ListenAndServeTLS(":443", "serve.pem","serve.key", nil)
	return
}
