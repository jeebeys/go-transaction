package example

import (
	"github.com/go-xorm/xorm"
	"github.com/jeebeys/go-transaction/transaction"
	"testing"
)

func TestDao(t *testing.T) {
	scanPath := `D:\src\workspace.golang.project\go-transaction\example`
	transaction.NewTransactionManager(transaction.TransactionConfig{ScanPath: scanPath}).RegisterDao(new(ExampleDao))

	dao := new(ExampleDao)
	dao.Select()
	dao.Update(new(xorm.Session), "") // auto commit
	dao.Delete(new(xorm.Session))     // handle fail and auto rollback
}
