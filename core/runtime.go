package core

import (
	"freemasonry.cc/trerr"
	"strings"
	"reflect"
	"runtime"
)

//()
func GetStructFuncName(object interface{}) string {
	structName := reflect.TypeOf(object).String() //
	pc := make([]uintptr, 1)
	runtime.Callers(2, pc)
	f := runtime.FuncForPC(pc[0])
	arry := strings.Split(f.Name(), ".")
	funcName := arry[len(arry)-1] //
	return structName + "." + funcName
}

//()
//
func GetFuncName() string {
	pc := make([]uintptr, 1)
	runtime.Callers(2, pc)
	f := runtime.FuncForPC(pc[0])
	arry := strings.Split(f.Name(), ".")
	return arry[len(arry)-1]
}

//()
//
//  +  
func GetPackageFuncName() string {
	pc := make([]uintptr, 1)
	runtime.Callers(2, pc)
	f := runtime.FuncForPC(pc[0])
	list := strings.Split(f.Name(), "/")
	return list[len(list)-1]
}

///,，*
func GetObjectName(ob interface{}) string {
	pv1 := reflect.ValueOf(ob)
	var object reflect.Type
	if pv1.Kind() == reflect.Ptr {
		object = pv1.Type().Elem()
		//fmt.Println("a:",object,"n:",object.Name())
	} else {
		object = reflect.TypeOf(ob)
		//fmt.Println("b:",object,"n:",object.Name())
	}
	return object.Name()
}

//
func Errformat(err error) error {
	//rpc error: code = InvalidArgument desc = failed to execute message
	if strings.Contains(err.Error(), "failed to execute message; message index: 0") {
		errContent := err.Error()
		errContent = strings.ReplaceAll(errContent, "rpc error: code = InvalidArgument desc = failed to execute message;", "")
		errContent = strings.ReplaceAll(errContent, ": invalid request", "")
		errContent = strings.ReplaceAll(errContent, " message index: 0: ", "")
		return trerr.TransError(errContent)
	}
	return err
}
