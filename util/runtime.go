package util

import (
	"freemasonry.cc/trerr"
	"reflect"
	"runtime"
	"strings"
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
	if strings.Contains(err.Error(), "failed to execute message; message index: 0") {
		errArry := strings.Split(err.Error(), ":")
		if len(errArry) > 3 {
			return trerr.TransError(errArry[2])
		}
	}
	return err
}
