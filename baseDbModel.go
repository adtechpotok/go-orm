package go_orm

import (
	"sync"
)

const StatusActive = "Active"

type BaseDbModel struct {
	cache      map[int]*PotokOrm
	cacheMutex *sync.RWMutex
	isInit     bool

	altIntIndex map[string]AltIntIndex
	altIntCache map[string]map[int]*PotokOrm

	altStringIndex map[string]AltStringIndex
	altStringCache map[string]map[string]*PotokOrm
}

func (m *BaseDbModel) GetCache() interface{} {
	m.init()
	return m.cache
}

func (m *BaseDbModel) init() {
	if m.isInit != true {
		m.cacheMutex = &sync.RWMutex{}
		m.cache = make(map[int]*PotokOrm)
		m.isInit = true
		m.altIntIndex = make(map[string]AltIntIndex)
		m.altIntCache = make(map[string]map[int]*PotokOrm)

		m.altStringIndex = make(map[string]AltStringIndex)
		m.altStringCache = make(map[string]map[string]*PotokOrm)
	}
}

func (m *BaseDbModel) RegisterIntIndex(string string, AltIndex AltIntIndex) {
	m.init()
	m.altIntIndex[string] = AltIndex
}

func (m *BaseDbModel) RegisterStringIndex(string string, AltIndex AltStringIndex) {
	m.init()
	m.altStringIndex[string] = AltIndex
}

func (m *BaseDbModel) FindInCache(id int) PotokOrm {
	m.init()
	m.cacheMutex.Lock()
	defer m.cacheMutex.Unlock()

	if resultVal, ok := m.cache[id]; ok {
		result := (*resultVal).(PotokOrm)
		if result.IsActive() {
			return result
		}

	}

	return nil

}

func (m *BaseDbModel) FindIndex(index string, id interface{}, onlyActive bool) PotokOrm {
	m.init()
	m.cacheMutex.Lock()
	defer m.cacheMutex.Unlock()

	switch id.(type) {
	case int:
		if _, ok := m.altIntIndex[index]; ok {
			if resultVal, ok := m.altIntCache[index][id.(int)]; ok {
				result := *resultVal
				if !onlyActive || result.IsActive() {
					return result
				}
			}
		}
	case string:
		if _, ok := m.altStringIndex[index]; ok {
			if resultVal, ok := m.altStringCache[index][id.(string)]; ok {
				result := *resultVal
				if !onlyActive || result.IsActive() {
					return result
				}
			}
		}
	default:
		panic("Only string or int key are available.")
	}

	return nil

}

func (m *BaseDbModel) AddToCache(v PotokOrm) {
	m.init()
	m.cacheMutex.Lock()
	defer m.cacheMutex.Unlock()

	if oldVal, ok := m.cache[v.GetId()]; ok {
		for id, alterIndexFunc := range m.altIntIndex {
			intIdToDelete := alterIndexFunc(*oldVal)
			delete(m.altIntCache[id], intIdToDelete)
		}

		for id, alterIndexFunc := range m.altStringIndex {
			stringIdToDelete := alterIndexFunc(*oldVal)
			delete(m.altStringCache[id], stringIdToDelete)
		}
	}

	m.cache[v.GetId()] = &v

	for id, val := range m.altIntIndex {
		if _, ok := m.altIntCache[id]; !ok {
			m.altIntCache[id] = make(map[int]*PotokOrm)
		}
		m.altIntCache[id][val(v)] = &v
	}

	for id, val := range m.altStringIndex {
		if _, ok := m.altIntCache[id]; !ok {
			m.altStringCache[id] = make(map[string]*PotokOrm)
		}
		m.altStringCache[id][val(v)] = &v
	}

}

func (m *BaseDbModel) ClearCache() {
	m.init()
	m.cacheMutex.Lock()
	defer m.cacheMutex.Unlock()

	m.cache = make(map[int]*PotokOrm)
	m.altIntCache = make(map[string]map[int]*PotokOrm)
	m.altStringCache = make(map[string]map[string]*PotokOrm)
}

func (m *BaseDbModel) Len() int {
	return len(m.cache)
}

type PotokOrm interface {
	IsActive() bool
	GetId() int
}

type AltIntIndex func(orm PotokOrm) int
type AltStringIndex func(orm PotokOrm) string
