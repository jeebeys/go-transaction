package transaction

import (
	"fmt"
	"github.com/jeebeys/go-transaction/aop"
	"github.com/jeebeys/go-transaction/ast"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"unicode"
)

var config TransactionConfig
var methodLocationMap = make(map[string]*struct{})

type joinPointSessionInfo struct {
	SessionIndex int // orm会话对象在方法入参列表的位置
}

type TransactionConfig struct {
	ScanPath string
}

type TransactionManager struct {
}

func NewTransactionManager(cfg TransactionConfig) *TransactionManager {
	if err := cfg.check(); err != nil {
		panic(err)
	}
	config = cfg
	scanGoFile()
	aop.RegisterAspect(new(Transactional))
	return new(TransactionManager)
}

func (c *TransactionConfig) check() error {
	_, err := os.Stat(c.ScanPath)
	return err
}

// TODO
func (c *TransactionConfig) Reload() error {
	return nil
}

func (t *TransactionManager) Register(daoLs ...interface{}) (tm *TransactionManager) {
	tm = t
	for _, v := range daoLs {
		aop.RegisterPoint(reflect.TypeOf(v))
	}
	return
}

func scanGoFile() {
	_ = filepath.Walk(config.ScanPath, walkFunc)
}

func walkFunc(fullPath string, info os.FileInfo, err error) error {
	if info == nil {
		return err
	}
	if info.IsDir() {
		return nil
	} else {
		if GO_FILE_SUFIX == path.Ext(fullPath) {
			// is go file
			cacheMethodLocationMap(fullPath)
		}
		return nil
	}
}

func cacheMethodLocationMap(fullPath string) {
	bt, err := os.ReadFile(fullPath)
	if err != nil {
		return
	}
	result := ast.ScanFuncDeclByComment("", string(bt), COMMENT_NAME)
	if result == nil {
		return
	}
	if result.RecvMethods != nil {
		for _, l := range result.RecvMethods {
			for _, v := range l {
				for i, s := range v.MethodName {
					// skip the private method
					if i == 0 && unicode.IsUpper(s) {
						methodLocationMap[fmt.Sprintf("%s.%s.%s", v.PkgName, v.RecvName, v.MethodName)] = new(struct{})
					}
					break
				}
			}
		}
	}
}
