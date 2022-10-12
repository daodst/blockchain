package keeper

import (
	"errors"
	"freemasonry.cc/blockchain/util"
	"github.com/cosmos/cosmos-sdk/types"
	"reflect"
	"strings"
)

var HelpTypeError = errors.New("HelpTypeError")
var HelpNilUnmarshal = errors.New("HelpNilUnmarshal")

type storeHelper struct {
	types.KVStore
}

func (s storeHelper) Has(key string) bool {
	return s.KVStore.Has([]byte(strings.ToLower(key)))
}
func (s storeHelper) Set(key string, value interface{}) error {
	valueType := reflect.TypeOf(value)
	var rawvalue []byte
	switch valueType.Kind() {
	//
	case reflect.String:
		rawvalue = []byte(value.(string))
	//
	case reflect.Struct:
		b, err := util.Json.Marshal(value)
		if err != nil {
			return err
		}
		rawvalue = b
	//
	case reflect.Slice:
		switch value.(type) {
		case []byte:
			rawvalue = value.([]byte)
		default:
			b, err := util.Json.Marshal(value)
			if err != nil {
				return err
			}
			rawvalue = b
		}
	case reflect.Map:
		b, err := util.Json.Marshal(value)
		if err != nil {
			return err
		}
		rawvalue = b
	default:
		return HelpTypeError
	}
	s.KVStore.Set([]byte(strings.ToLower(key)), rawvalue)
	return nil
}

func (s storeHelper) Get(key string) []byte {
	return s.KVStore.Get([]byte(strings.ToLower(key)))
}
func (s storeHelper) GetUnmarshal(key string, v interface{}) error {
	b := s.KVStore.Get([]byte(strings.ToLower(key)))
	if b == nil {
		return HelpNilUnmarshal
	}
	return util.Json.Unmarshal(b, v)
}
func (s storeHelper) Delete(key string) {
	s.KVStore.Delete([]byte(strings.ToLower(key)))
}

func (s storeHelper) KVStorePrefixIterator(key string) types.Iterator {
	return types.KVStorePrefixIterator(s.KVStore, []byte(strings.ToLower(key)))
}
