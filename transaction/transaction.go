package transaction

import (
	"fmt"
	"reflect"

	"github.com/go-xorm/xorm"
	"github.com/jeebeys/go-transaction/aop"
	"gorm.io/gorm"
)

var methodSessionMap = make(map[string]*joinPointSessionInfo)

type Transactional struct {
	ReadOnly                bool
	RollbackFor             reflect.Type
	RollbackForStructName   reflect.Type
	NoRollbackFor           reflect.Type
	NoRollbackForStructName reflect.Type
	Propagation             Propagation
	Isolation               int
	Timeout                 int
}

func (t *Transactional) Before(point *aop.JoinPoint, methodLocation string) bool {

	if methodSessionMap[methodLocation] != nil {
		t.doSessionBegin(point.Params[methodSessionMap[methodLocation].SessionIndex].Interface())
	} else {
		// cache
		for i, v := range point.Params {
			if t.doSessionBegin(v.Interface()) {
				methodSessionMap[methodLocation] = &joinPointSessionInfo{
					SessionIndex: i,
				}
				break
			}
		}
	}
	fmt.Println("transactional before", methodLocation, methodSessionMap, point.Params)
	return true
}

func (t *Transactional) After(point *aop.JoinPoint, methodLocation string) {
	if methodSessionMap[methodLocation] != nil {
		// TODO 规约：返回值第一个参数为处理结果，类型为布尔型。由此确认是提交还是回滚
		for i, v := range point.Result {
			if i == 0 {
				switch result := v.Interface().(type) {
				case bool:
					if result {
						t.doSessionCommit(point.Params[methodSessionMap[methodLocation].SessionIndex].Interface())
					} else {
						t.doSessionRollback(point.Params[methodSessionMap[methodLocation].SessionIndex].Interface())
					}
				}
			} else {
				break
			}
		}
	}

	fmt.Println("transactional After", methodLocation, methodSessionMap, point.Result)
}

func (t *Transactional) Finally(point *aop.JoinPoint, methodLocation string) {
	if methodSessionMap[methodLocation] != nil {
		t.doSessionClose(point.Params[methodSessionMap[methodLocation].SessionIndex].Interface())
	} else {

	}
	fmt.Println("transactional Finally", methodLocation, methodSessionMap, point.Result)
}

func (t *Transactional) IsMatch(methodLocation string) bool {
	if methodLocationMap[methodLocation] != nil {
		return true
	} else {
		return false
	}
}

func (t *Transactional) doSessionBegin(v interface{}) bool {
	switch ses := v.(type) {
	case *xorm.Session:
		_ = ses.Begin()
		return true
	case *gorm.DB:
		_ = ses.Begin()
		return true
	default:
		return false
	}
}

func (t *Transactional) doSessionCommit(v interface{}) bool {
	switch ses := v.(type) {
	case *xorm.Session:
		_ = ses.Commit()
		return true
	case *gorm.DB:
		_ = ses.Commit()
		return true
	default:
		return false
	}
}

func (t *Transactional) doSessionRollback(v interface{}) bool {
	switch ses := v.(type) {
	case *xorm.Session:
		_ = ses.Rollback()
		return true
	case *gorm.DB:
		_ = ses.Rollback()
		return true
	default:
		return false
	}
}

func (t *Transactional) doSessionClose(v interface{}) bool {
	fmt.Println("close", v)
	switch ses := v.(type) {
	case *xorm.Session:
		ses.Close()
		return true
	case *gorm.DB:
		return true
	default:
		return false
	}
}
