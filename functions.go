package configure

import (
	"fmt"
	"reflect"
)

//Load 读取配置
//confmap 传入待检索对象, providers 传入配置提供器实例
func Load(confmap map[string]interface{}, config interface{}) error {
	//如果这个类型不是引用，那么后继是无法设置值的，报错给用户
	if !reflect.Indirect(reflect.ValueOf(config)).CanAddr() {
		return fmt.Errorf("config: LoadConfiguration(non-pointer %s)", reflect.ValueOf(config).Type().String())
	}
	//获取当前配置器类型
	v := reflect.New(reflect.TypeOf(config).Elem())
	//返回该提供器非指针实例
	v = reflect.Indirect(v)
	//获取当前提供器实例对象（注意，这里如果是引用将会造成无法执行 .NumField 调用）
	t := reflect.TypeOf(v.Interface())
	//如果该实例并不是一个结构，那么是非法的，报告给用户
	if t.Kind() != reflect.Struct {
		return fmt.Errorf("config: LoadConfiguration(non-struct %s)", v.Type().String())
	}
	//遍历所有的字段
	for i := 0; i < t.NumField(); i++ {
		//获取上下文目标字段
		tfield := t.Field(i)
		if tfield.Type.Kind() != reflect.Struct {
			continue
		}
		//检索是否包含 tags.required 标记
		_, nilable := tfield.Tag.Lookup(tags.nilable)
		_, external := tfield.Tag.Lookup(tags.external)
		//检索指定键名
		mapkey := tfield.Tag.Get(tags.key)
		if mapkey == "" {
			//回落到字段名
			mapkey = tfield.Name
		}
		//判断传入的映射中是否包含键对应的字段
		mapval, existed := confmap[mapkey]
		if mapval == nil && external {
			//如果数据是 nil 并且还被标记为导出字段那就不管了
			continue
		}
		if !existed && !nilable {
			//不包含，且 tags.nilable 的字段没不存在，那么需要报错给用户
			return fmt.Errorf("config: the key %s is required but get nil", mapkey)
		}
		//获取字段所对应的反射值信息
		vfield := v.Field(i)
		if !vfield.CanSet() {
			//如果无法设置，应该是传入的实例非引用
			return fmt.Errorf("config: the key %s is irregular", mapkey)
		}
		//如果类型为目标结构，那么需要构建该结构
		if tfield.Type == types.provider || tfield.Type == types.providerptr {
			//通过反射创建该类型的实例
			instance := Field{
				key:     mapkey,
				val:     mapval,
				existed: existed,
				tags:    tfield.Tag,
			}
			if tfield.Type == types.provider {
				vfield.Set(reflect.ValueOf(instance))
			} else {
				vfield.Set(reflect.ValueOf(&instance))
			}
		} else {
			//否则视为子配置信息
			for vfield.Type().Kind() == reflect.Ptr {
				//这里按层次构建所有的 nil 引用实例
				vfield.Set(reflect.New(vfield.Type().Elem()))
				vfield = vfield.Elem()
			}
			//如果对象是 nil，那就 nil 吧，下一位。
			if mapval == nil {
				continue
			}
			//递归执行子配置层次
			err := Load(mapval.(map[string]interface{}), vfield.Addr().Interface())
			if err != nil {
				return err
			}
		}
	}
	reflect.Indirect(reflect.ValueOf(config)).Set(v)
	return nil
}

//Name 字段名
func (provider Field) Name() string {
	return provider.key
}

//Get 获取字段值
func (provider Field) Get(defaultvalue ...interface{}) interface{} {
	if !provider.existed && len(defaultvalue) > 0 {
		return defaultvalue[0]
	}
	return provider.val
}

//Lookup 查找字段值
func (provider Field) Lookup() (interface{}, bool) {
	return provider.val, provider.existed
}

//GetAttribute 获取标签属性
func (provider Field) GetAttribute(name string) string {
	return provider.tags.Get(name)
}

//LookupAttribute 获取是否存在标签
func (provider Field) LookupAttribute(name string) (string, bool) {
	return provider.tags.Lookup(name)
}
