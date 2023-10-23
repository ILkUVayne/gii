package orm

import (
	_ "github.com/go-sql-driver/mysql"
	"testing"
)

type UserAddr struct {
	Id   int    `orm:"primaryKey;column:id;NOT NULL;AUTO_INCREMENT" json:"id"`
	Addr string `orm:"column:addr;type:varchar(255)" json:"addr"`
	No   int    `orm:"column:no" json:"no"`
}

var (
	addr1 = &UserAddr{Addr: "xxxx路1号", No: 18}
	addr2 = &UserAddr{Addr: "xxxx路2号", No: 25}
	addr3 = &UserAddr{Addr: "xxxx路3号", No: 25}
)

func TestSession_Insert(t *testing.T) {
	engine := NewEngine("mysql", "root:root@tcp(127.0.0.1:3306)/gii")
	defer engine.Close()
	s := engine.NewSession().Model(&UserAddr{})
	s.DropTable()
	s.CreateTable()
	i, err := s.Insert(addr1, addr2, addr3)
	if err != nil {
		t.Error(err)
	}
	if i != 3 {
		t.Error("insert num not 3")
	}
}

func TestSession_All(t *testing.T) {
	engine := NewEngine("mysql", "root:root@tcp(127.0.0.1:3306)/gii")
	defer engine.Close()
	s := engine.NewSession().Model(&UserAddr{})
	var userAddr []UserAddr
	s.Where("id > ?", 1).OrderBy("id desc").Limit(2).All(&userAddr)
	if len(userAddr) != 2 {
		t.Error("failed to query all")
	}
}

func TestSession_Update(t *testing.T) {
	engine := NewEngine("mysql", "root:root@tcp(127.0.0.1:3306)/gii")
	defer engine.Close()
	s := engine.NewSession().Model(&UserAddr{})
	i := s.Where("id  = ?", 2).Update("addr", "ssss路1号", "no", 255)
	if i != 1 {
		t.Errorf("update line failed %d", i)
	}
	s.Where("id  = ?", 3).Update(map[string]interface{}{
		"addr": "qqqq路1号",
		"no":   333,
	})
	var userAddr []UserAddr
	s.Where("id  = ?", 3).All(&userAddr)
	if userAddr[0].Addr != "qqqq路1号" || userAddr[0].No != 333 {
		t.Errorf("update line failed id = %d", i)
	}
}

func TestSession_Count(t *testing.T) {
	engine := NewEngine("mysql", "root:root@tcp(127.0.0.1:3306)/gii")
	defer engine.Close()
	s := engine.NewSession().Model(&UserAddr{})
	if s.Count() != 3 {
		t.Error("get userAddr count failed")
	}
}

func TestSession_First(t *testing.T) {
	engine := NewEngine("mysql", "root:root@tcp(127.0.0.1:3306)/gii")
	defer engine.Close()
	s := engine.NewSession().Model(&UserAddr{})
	var userAddr []UserAddr
	s.Where("addr LIKE '%路%'").First(&userAddr)
	if len(userAddr) != 1 {
		t.Error("get userAddr First failed")
	}
}
