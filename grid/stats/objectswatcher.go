package stats

import (
	"fmt"
	"sync"
	"time"
)

type Obj struct {
	objType string
	objId   string
}

type ObjInfo struct {
	objId      string
	statistics string

	objCreateTime time.Time
}

func (c *ObjInfo) ObjId() string {
	return c.objId
}

func (c *ObjInfo) Statistics() string {
	return c.statistics
}

func (c *Obj) SetStatistics(statistics string) {
	return
	if oType, ok := objWatcher.objOfType[c.objType]; ok {
		if objInfo, ok := oType.objects[c.objId]; ok {
			objInfo.statistics = statistics
		}
	}
}

type objectsOfType struct {
	objects map[string]*ObjInfo
}

type ObjectsWatcher struct {
	mtx       sync.Mutex
	objOfType map[string]*objectsOfType
}

var objWatcher ObjectsWatcher

func init() {
	objWatcher.objOfType = make(map[string]*objectsOfType)
}

func (c *Obj) InitObj(objType string, objId string) {
	objWatcher.mtx.Lock()
	defer objWatcher.mtx.Unlock()

	c.objId = objId
	c.objType = objType

	var info ObjInfo
	info.objId = objId
	info.objCreateTime = time.Now()

	if _, ok := objWatcher.objOfType[objType]; !ok {
		var objOfType objectsOfType
		objOfType.objects = make(map[string]*ObjInfo)
		objWatcher.objOfType[objType] = &objOfType
	}

	if v, ok := objWatcher.objOfType[objType]; ok {
		v.objects[objId] = &info
	}
}

func (c *Obj) UninitObj() {
	objWatcher.mtx.Lock()
	defer objWatcher.mtx.Unlock()

	if t, ok := objWatcher.objOfType[c.objType]; ok {
		delete(t.objects, c.objId)
	}
}

func ObjectTypes() []string {
	objWatcher.mtx.Lock()
	defer objWatcher.mtx.Unlock()

	result := make([]string, 0)
	for key := range objWatcher.objOfType {
		result = append(result, key)
	}

	return result
}

func ObjectsOfTypes(objType string) []*ObjInfo {
	objWatcher.mtx.Lock()
	defer objWatcher.mtx.Unlock()

	result := make([]*ObjInfo, 0)

	if v, ok := objWatcher.objOfType[objType]; ok {
		for _, o := range v.objects {
			result = append(result, o)
		}
	}
	return result
}

func Dump() {
	fmt.Println("Memory dump:")
	types := ObjectTypes()
	found := false
	for _, t := range types {
		objects := ObjectsOfTypes(t)
		if len(objects) > 0 {
			fmt.Println("Type ", t, " count ", len(objects))
			for _, o := range objects {
				fmt.Println("\tId: ", o.objId, "\t ", o.statistics)
				found = true
			}
		}
	}

	if !found {
		fmt.Println("\tNo objects")
	}
}

func AllObjects() []*ObjInfo {
	result := make([]*ObjInfo, 0)
	types := ObjectTypes()
	for _, t := range types {
		objects := ObjectsOfTypes(t)
		if len(objects) > 0 {
			for _, o := range objects {
				result = append(result, o)
			}
		}
	}
	return result
}
