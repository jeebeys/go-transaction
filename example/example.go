package example

import (
	"fmt"

	"github.com/go-xorm/xorm"
)

type ExampleDao struct {
}

func (e *ExampleDao) Select() (bool, error) {
	return true, nil
}

// Update
// @Transactional
//
//go:noinline
func (d *ExampleDao) Update( /*@Header()*/ s *xorm.Session, param string) (bool, error) {
	fmt.Println("update param:", param)
	return true, nil
}

// Delete
// @Transactional
func (d *ExampleDao) Delete(s *xorm.Session) (bool, error) {
	return false, nil
}
