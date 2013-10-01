package core

import (
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"xuncache/types"
	"xuncache/util"
)

//存储类型
var Store_list = &List{name: make(map[string][]*Store)}

//队列类型
type List struct {
	name map[string][]*Store
	time uint64
	lock sync.RWMutex
}

//原子类型
type Store struct {
	id    uint64
	value map[string]interface{}
	lock  sync.RWMutex
}

//错误类型
var err error

//基本协议类型
type Basic struct {
	Protocol string
	Password string
	Sources  *util.Json
}

//返回类型
type Result struct {
	Errors    bool
	Resources map[string]interface{}
}

//错误类型
type Error struct {
	Errors bool
	Point  string
}

//定位索引类型 Locate

//初始化控制器
func Init(Receive *util.Json) *Basic {
	Protocol, _ := Receive.Get("Protocol").String()
	Password, _ := Receive.Get("Pass").String()
	Basic_model := &Basic{
		Protocol: Protocol,
		Password: Password,
		Sources:  Receive,
	}
	return Basic_model
}

//调度 节省内存开销
func (Receive *Basic) Handle() (result []byte) {

	switch Receive.Protocol {
	case `push`:
		return Receive.push()
	case `save`:
		return Receive.save()
	case `find`:
		return Receive.find()
	case `query`:
		return Receive.query()
	case `incr`:
		return Receive.incr()
	case `decr`:
		return Receive.decr()
	case `del`:
		return Receive.del()
	case `delete`:
		return Receive.delete()
	default:
		return Errors(errors.New("Protocol error!"))
		break
	}
	return Errors(errors.New("Core Control error!"))

}

func (Receive *Basic) push() (result []byte) {
	//接收
	Push := &types.Push{}
	//key类型
	Push.Key, err = Receive.Sources.Get("key").String()
	if err != nil {
		return Errors(errors.New("key does not exist!"))
	}
	Push.Index = Receive.Sources.Get("index").Uint64()
	//初始化类型
	stroe := &Store{
		value: make(map[string]interface{}),
	}
	//获取数据
	stroe.value, err = Receive.Sources.Get("data").Map()
	if err != nil {
		return Errors(errors.New("key does not exist!"))
	}
	//上锁
	Store_list.lock.Lock()
	defer Store_list.lock.Unlock()
	//获取客户端自定义_id(唯一标识符)用于mysql等数据库的数据导入--mysql自增主键
	max_id := Store_list.List_max_Id(Push.Key)
	if max_id < Push.Index {
		max_id = Push.Index
	}
	stroe.id = max_id
	Store_list.name[Push.Key] = append(Store_list.name[Push.Key], stroe)
	//返回结果
	Back := &types.Back_Push{
		Errors: false,
		Id:     stroe.id,
	}
	result, _ = json.Marshal(Back)
	return
}

func (Receive *Basic) save() (result []byte) {
	return []byte("test")
}
func (Receive *Basic) find() (result []byte) {
	return []byte("test")
}
func (Receive *Basic) query() (result []byte) {
	return []byte("test")
}
func (Receive *Basic) incr() (result []byte) {
	return []byte("test")
}
func (Receive *Basic) decr() (result []byte) {
	return []byte("test")
}
func (Receive *Basic) del() (result []byte) {
	return []byte("test")
}
func (Receive *Basic) delete() (result []byte) {
	return []byte("test")
}

//获取最大Id
func (Lists *List) List_max_Id(key string) uint64 {
	max_nums := len(Lists.name[key]) - 1
	if max_nums < 0 {
		return uint64(1)
	}
	return Lists.name[key][max_nums].id + 1
}

//输出错误类型
func Errors(Error_msg error) (result []byte) {
	errors := &Error{
		Errors: true,
		Point:  fmt.Sprint(Error_msg),
	}
	result, _ = json.Marshal(errors)
	return
}

/*
//二分法定位 主键(_id) --可以识别不存在
func (Lists *List) Dichotomy(key string, min, max, index int) (slice_index int, err bool) {
	if index < 1 || max < 1 {
		return -1, true
	}

	i := int((min + max) / 2)
	if Store_list[key][i]["_id"].(int) == index {
		slice_index = i
	} else {
		if Store_list[key][i]["_id"].(int) > index {
			if max == i {
				return -1, true
			}
			slice_index, _ = Lists.Dichotomy(key, min, i, index)
		} else {
			if min == i {
				return -1, true
			}
			slice_index, _ = Lists.Dichotomy(key, i, max, index)
		}
	}
	if slice_index < 0 {
		return -1, true
	}
	return slice_index, false
}
*/
