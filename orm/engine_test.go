package orm

import (
	"errors"
	"gii/orm/session"
	"testing"
)

type Users struct {
	Name string `orm:"primaryKey;type:varchar(70);column:name" json:"name"`
	Age  int    `orm:"column:age" json:"age"`
}

type Address struct {
	Name string `orm:"primaryKey;type:varchar(70);column:name" json:"name"`
	Age  int    `orm:"column:age" json:"age"`
}

func TestEngine_Transaction(t *testing.T) {
	t.Run("rollback", func(t *testing.T) {
		transactionRollback(t)
	})
	t.Run("commit", func(t *testing.T) {
		transactionCommit(t)
	})
}
func transactionRollback(t *testing.T) {
	engine := NewEngine("mysql", "root:root@tcp(127.0.0.1:3306)/gii")
	defer engine.Close()
	s := engine.NewSession()
	s.Model(&Users{}).DropTable()
	s.Model(&Users{}).CreateTable()
	s.Model(&Address{}).DropTable()
	s.Model(&Address{}).CreateTable()
	_, err := engine.Transaction(func(s *session.Session) (result interface{}, err error) {
		// mysql DDL语句执行时会隐式执行commit,导致rollback失败
		//s.Model(&Users{}).CreateTable()
		_, err = s.Insert(&Users{"Tom", 18})
		_, err = s.Insert(&Address{"Jek", 19})
		return nil, errors.New("Error")
	})
	if err == nil {
		t.Fatal("failed to rollback")
	}
}
func transactionCommit(t *testing.T) {
	engine := NewEngine("mysql", "root:root@tcp(127.0.0.1:3306)/gii")
	defer engine.Close()
	s := engine.NewSession()
	s.Model(&Users{}).DropTable()
	s.Model(&Users{}).CreateTable()
	s.Model(&Address{}).DropTable()
	s.Model(&Address{}).CreateTable()
	_, err := engine.Transaction(func(s *session.Session) (result interface{}, err error) {
		_, err = s.Insert(&Users{"Tom", 18})
		_, err = s.Insert(&Address{"Jek", 19})
		return
	})
	var u []Users
	s.First(&u)
	if err != nil || u[0].Name != "Tom" {
		t.Fatal("failed to commit")
	}
}
