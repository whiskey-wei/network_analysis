package handlepacket

import "sync"

type BaseMap struct {
	sync.Map //使用sync.Map而不是Map保证线程安全
}

func (m *BaseMap) Add(info *PacketInfo) {
	key := GetHashKey(info)
	val, ok := m.Load(key)
	if !ok {
		//没有这条记录
		mapVal := MapValue{}
		mapVal.List.InsertBefore(info, mapVal.List.Front())
		mapVal.DataSize += info.DataSize
		m.Store(key, mapVal)
	} else {
		mapVal, ok := val.(MapValue)
		if !ok {
			return
		}
		mapVal.List.InsertBefore(info, mapVal.List.Front())
		mapVal.DataSize += info.DataSize
	}
}

func (m *BaseMap) Del(info *PacketInfo) {
	m.Delete(*info)
}

var TcpMap BaseMap
var UdpMap BaseMap
