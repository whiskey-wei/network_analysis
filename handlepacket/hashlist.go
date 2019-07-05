package handlepacket

import (
	"sync"
)

var HashMux sync.Mutex

type BaseMap struct {
	sync.Map //使用sync.Map而不是Map保证线程安全
}

func (m *BaseMap) Add(info *PacketInfo) {
	key := GetHashKey(info)
	val, ok := m.Load(key)
	if !ok {
		//没有这条记录
		//fmt.Println(*info)
		HashMux.Lock()
		mapVal := &MapValue{}
		mapVal.List.PushFront(info)
		mapVal.NowSize += info.DataSize
		mapVal.NowCount++
		HashMux.Unlock()
		m.Store(key, mapVal)
	} else {
		mapVal, ok := val.(*MapValue)
		if !ok {
			return
		}
		HashMux.Lock()
		mapVal.List.PushFront(info)
		mapVal.NowSize += info.DataSize
		mapVal.NowCount++
		HashMux.Unlock()
	}
}

func (m *BaseMap) Del(info PacketInfo) {
	key := GetHashKey(&info)
	m.Delete(key)
}

var TcpMap BaseMap
var UdpMap BaseMap
var SumMap BaseMap
