package db

import (
	"fmt"
	_ "github.com/mattn/go-sqlite3"
)

type RombaseEnum struct {
	Type string // 类型
	Name string // 名称
	Sort uint32 // 排序
}

func (*RombaseEnum) TableName() string {
	return "rombase_enum"
}

// 批量写入
func (m *RombaseEnum) BatchAdd(romlist []*RombaseEnum) error {

	if len(romlist) == 0 {
		return nil
	}

	tx := getDb().Begin()
	for _, v := range romlist {
		tx.Create(&v)
	}
	tx.Commit()
	return nil
}

// 读取全部数据
func (*RombaseEnum) GetAll() ([]*RombaseEnum, error) {
	volist := []*RombaseEnum{}
	result := getDb().Order("sort Asc").Find(&volist)
	if result.Error != nil {
		fmt.Println(result.Error.Error())
	}
	return volist, result.Error
}

// 根据类型读取数据
func (*RombaseEnum) GetByType(t string) ([]*RombaseEnum, error) {
	volist := []*RombaseEnum{}
	result := getDb().Where("type=?", t).Order("sort Asc").Find(&volist)
	if result.Error != nil {
		fmt.Println(result.Error.Error())
	}
	return volist, result.Error
}

// 删除一个类型
func (m *RombaseEnum) DeleteByType() error {
	result := getDb().Where("type=? ", m.Type).Delete(&m)
	return result.Error
}
