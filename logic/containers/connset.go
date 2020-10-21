package containers

//// ConnSet 保存连接的集合
//// 包括了普通的连接，还有注册为服务的连接的映射
//type ConnSet struct {
//	connLock     sync.RWMutex
//	indexLock    sync.RWMutex
//	indexReLock  sync.RWMutex
//	connMap      map[string]conn.Conn
//	serviceIndex map[string]string
//
//	// 服务映射的反转，删除的时候用的，用空间换时间
//	serviceIndexReverse map[string][]string
//}
//
//func NewConnSet() *ConnSet {
//	return &ConnSet{
//		connLock:            sync.RWMutex{},
//		indexLock:           sync.RWMutex{},
//		indexReLock:         sync.RWMutex{},
//		connMap:             make(map[string]conn.Conn),
//		serviceIndex:        make(map[string]string),
//		serviceIndexReverse: make(map[string][]string),
//	}
//}
//
//func (s *ConnSet) GetByUUID(uuid string) (c conn.Conn, exist bool) {
//	s.connLock.RLock()
//
//	c, exist = s.connMap[uuid]
//
//	s.connLock.RUnlock()
//
//	return
//}
//
//func (s *ConnSet) GetByServiceName(serviceName string) (c conn.Conn, exist bool) {
//	s.indexLock.RLock()
//	index, exist := s.serviceIndex[serviceName]
//	if !exist {
//		return
//	}
//	s.indexLock.RUnlock()
//	s.connLock.RLock()
//	c, exist = s.connMap[index]
//	if !exist {
//		panic("bug!!!!!!!!!!!!!!!")
//	}
//	s.connLock.RUnlock()
//	return
//}
//
//func (s *ConnSet) Add(c conn.Conn) {
//	s.connLock.Lock()
//	defer s.connLock.Unlock()
//
//	s.connMap[c.Uuid()] = c
//}
//
//func (s *ConnSet) register(uuid, serviceName string) (err error) {
//	s.connLock.RLock()
//
//	_, exist := s.connMap[uuid]
//	if !exist {
//		err = fmt.Errorf("connect uuid [%s] not found", uuid)
//		return
//	}
//
//	s.connLock.RUnlock()
//
//	s.indexLock.Lock()
//
//	if _, exist := s.serviceIndex[serviceName]; exist {
//		err = fmt.Errorf("service [%s] exist!", serviceName)
//		s.indexLock.Unlock()
//		return
//	}
//	s.serviceIndex[serviceName] = uuid
//
//	s.indexLock.Unlock()
//
//	s.indexReLock.Lock()
//
//	l, ok := s.serviceIndexReverse[uuid]
//	if ok {
//		s.serviceIndexReverse[uuid] = append(l, serviceName)
//	} else {
//		s.serviceIndexReverse[uuid] = []string{serviceName}
//	}
//
//	s.indexReLock.Unlock()
//
//	return
//}
//
//// UnSignByServiceName 注销连接的服务映射，但是连接仍然存在
//func (s *ConnSet) UnSign(serviceName string) (err error) {
//	s.indexLock.Lock()
//
//	uuid, ok := s.serviceIndex[serviceName]
//	if !ok {
//		err = fmt.Errorf("service [%s] not exist!", serviceName)
//		s.indexLock.Unlock()
//		return
//	}
//	delete(s.serviceIndex, serviceName)
//
//	s.indexLock.Unlock()
//
//	s.indexReLock.Lock()
//	// 只要服务有，那么这个肯定有
//	l, ok := s.serviceIndexReverse[uuid]
//	if !ok {
//		panic("bug!!!!!!!!")
//	}
//
//	if len(l) > 1 {
//		var pos int
//		flag := false
//		for pos = 0; pos < len(l); pos++ {
//			if l[pos] == serviceName {
//				flag = true
//				break
//			}
//		}
//
//		if !flag {
//			panic("bug!!!!!!!")
//		}
//
//		s.serviceIndexReverse[uuid] = append(l[:pos], l[pos+1:]...)
//		// 该连接下注册了多个服务
//	} else {
//		// 该连接只有一个服务
//		delete(s.serviceIndexReverse, uuid)
//	}
//	s.indexReLock.Unlock()
//
//	return
//}
//
//// Remove 根据uuid 删除掉连接
//func (s *ConnSet) Remove(uuid string) (err error) {
//
//	// 删除连接
//	s.connLock.Lock()
//
//	delete(s.connMap, uuid)
//
//	s.connLock.Unlock()
//
//	// 删除服务
//	s.indexReLock.RLock()
//	l, exist := s.serviceIndexReverse[uuid]
//	s.indexReLock.RUnlock()
//
//	if exist {
//		for _, service := range l {
//			if err = s.UnSign(service); err != nil {
//				return
//			}
//		}
//	}
//
//	return
//}
