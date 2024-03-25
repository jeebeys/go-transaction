package example

import (
	"github.com/go-xorm/xorm"
	"github.com/jeebeys/go-transaction/transaction"
	"testing"
)

func TestDao(t *testing.T) {
	scanPath := `D:\src\workspace.golang.project\go-transaction\example`
	transaction.NewTransactionManager(transaction.TransactionConfig{ScanPath: scanPath}).Register(new(ExampleDao))

	dao := new(ExampleDao)
	_, _ = dao.Select()
	_, _ = dao.Update(new(xorm.Session), "") // auto commit
	_, _ = dao.Delete(new(xorm.Session))     // handle fail and auto rollback
}
