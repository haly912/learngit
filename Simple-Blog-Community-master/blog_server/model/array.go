// model/array.go
package model

import (
	"database/sql/driver"
	"strings"
)

// Array 是一个字符串切片的别名，用于表示一个字符串数组。
type Array []string

// Scan 方法用于从数据库读取数据后，对其进行处理，获得Go类型的变量。
// 它接受一个interface{}类型的参数val，这通常是数据库驱动返回的原始数据。
// 此+。
func (m *Array) Scan(val interface{}) error {
	// 断言val为[]byte类型
	s := val.([]uint8)
	// 将字节切片转换为字符串，并使用"|"作为分隔符来分割字符串
	ss := strings.Split(string(s), "|")
	// 将解析后的字符串数组赋值给Array类型的实例
	*m = ss
	// 返回nil表示没有错误发生
	return nil
}

// Value 方法用于将数据存到数据库时，对数据进行处理，获得数据库支持的类型。
// 此方法将Array类型的所有字符串元素使用"|"作为分隔符连接成一个单一的字符串。
// 返回连接后的字符串和nil错误，表示数据可以被存储到数据库中。
func (m Array) Value() (driver.Value, error) {
	// 使用"|"作为分隔符，将Array类型的所有字符串元素连接成一个单一的字符串
	str := strings.Join(m, "|")
	// 返回连接后的字符串和nil错误
	return str, nil
}
