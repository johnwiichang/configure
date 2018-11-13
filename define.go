package configure

import (
	"encoding/json"
	"fmt"
	"reflect"
)

var (
	//tags 标签枚举器
	tags = struct {
		nilable,
		key,
		external string
	}{
		"nilable",  //可空标签
		"key",      //陈述键
		"external", //外部标签
	}

	types = struct {
		provider,
		providerptr reflect.Type
	}{
		reflect.TypeOf(Field{}),  //提供器类型
		reflect.TypeOf(&Field{}), //提供器引用类型
	}
)

type (
	//Field 提供器
	Field struct {
		key     string            //键名
		val     interface{}       //值信息
		existed bool              //存在标签
		tags    reflect.StructTag //所有标签
	}
)

func (f Field) String() string {
	return fmt.Sprintf("%v", f.val)
}

//MarshalJSON JSON 序列化接口实现
func (f Field) MarshalJSON() (dst []byte, err error) {
	return json.Marshal(f.val)
}
