package orm

import (
	"gii/orm/session"
	"github.com/ILkUVayne/utlis-go/v2/ulog"
	_ "github.com/go-sql-driver/mysql"
	"testing"
)

type UserAddr struct {
	Id     int    `orm:"primaryKey;column:id;NOT NULL;AUTO_INCREMENT" json:"id"`
	Addr   string `orm:"column:addr;type:varchar(255)" json:"addr"`
	No     int    `orm:"column:no" json:"no"`
	IdCard string `orm:"column:id_card;type:varchar(70)" json:"id_card"`
}

func (a *UserAddr) BeforeInsert(s *session.Session) error {
	ulog.Info("before inert", a)
	a.No = 888
	return nil
}

func (a *UserAddr) AfterInsert(s *session.Session) error {
	ulog.Info("after inert", a)
	return nil
}

func (a *UserAddr) BeforeQuery(s *session.Session) error {
	ulog.Info("before query", a)
	return nil
}

func (a *UserAddr) AfterQuery(s *session.Session) error {
	ulog.Info("after query", a)
	a.IdCard = "************"
	return nil
}

func (a *UserAddr) BeforeUpdate(s *session.Session) error {
	ulog.Info("before update", a)
	a.Addr = "修改啊啊啊"
	return nil
}

func (a *UserAddr) AfterUpdate(s *session.Session) error {
	ulog.Info("after update", a)
	return nil
}

var (
	addr1 = &UserAddr{Addr: "xxxx路1号", IdCard: "11111", No: 18}
	addr2 = &UserAddr{Addr: "xxxx路2号", IdCard: "22222", No: 25}
	addr3 = &UserAddr{Addr: "xxxx路3号", IdCard: "33333", No: 25}
	addr4 = &UserAddr{Addr: "xxxx路3号", IdCard: "d1", No: 25}
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
	s.Where("id  = ?", 3).Update(map[string]any{
		"addr": "qqqq路1号",
		"no":   333,
	})
	var userAddr []UserAddr
	s.Where("id  = ?", 3).All(&userAddr)
	if userAddr[0].Addr != "qqqq路1号" || userAddr[0].No != 333 {
		t.Errorf("update line failed id = %d", i)
	}
}

func TestSession_Delete(t *testing.T) {
	engine := NewEngine("mysql", "root:root@tcp(127.0.0.1:3306)/gii")
	defer engine.Close()
	s := engine.NewSession().Model(&UserAddr{})
	_, err := s.Insert(addr4)
	if err != nil {
		t.Error(err)
	}

	i := s.Where("id_card = ?", "d1").Delete()
	if i != 1 {
		t.Error("delete failed")
	}
	var userAddr []UserAddr
	s.Where("id_card = ?", "d1").First(&userAddr)
	if len(userAddr) == 1 {
		t.Error("delete failed")
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
	s.Where("addr LIKE '%路%'").OrderBy("id desc").First(&userAddr)
	if len(userAddr) != 1 {
		t.Error("get userAddr First failed")
	}
}
