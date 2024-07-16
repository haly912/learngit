// model/time.go
package model

import (
	"database/sql/driver"
	"fmt"
	"time"
)

// timeFormat 定义了时间的格式，符合 Go 的时间格式化规则。
const timeFormat = "2006-01-02 15:04:05"

// timezone 定义了时区，这里设置为东八区（北京时间）。
const timezone = "Asia/Shanghai"

// Time 是对 time.Time 的一个封装，使得它可以自定义 JSON 编码和解码行为。
type Time time.Time

// MarshalJSON 实现了 Time 类型的 JSON 编码方法。
func (t Time) MarshalJSON() ([]byte, error) {
	b := make([]byte, 0, len(timeFormat)+2)      // 分配足够的字节用于格式化时间
	b = append(b, '"')                           // 添加 JSON 字符串开始标记
	b = time.Time(t).AppendFormat(b, timeFormat) // 格式化时间
	b = append(b, '"')                           // 添加 JSON 字符串结束标记
	return b, nil                                // 返回编码后的字节和 nil 错误
}

// UnmarshalJSON 实现了 Time 类型的 JSON 解码方法。
func (t *Time) UnmarshalJSON(data []byte) (err error) {
	// 解码 JSON 字符串到 time.Time，忽略时区
	now, _ := time.ParseInLocation(`"`+timeFormat+`"`, string(data), time.Local)
	*t = Time(now) // 更新 Time 类型的值
	return         // 返回 nil 错误
}

// String 方法返回 Time 类型的时间字符串表示。
func (t Time) String() string {
	return time.Time(t).Format(timeFormat) // 使用定义的时间格式
}

// local 方法将 UTC 时间转换为指定的时区时间。
func (t Time) local() time.Time {
	loc, _ := time.LoadLocation(timezone) // 加载时区
	return time.Time(t).In(loc)           // 转换为本地时间
}

// Value 方法实现了数据库驱动的 Valuer 接口，用于数据库交互。
func (t Time) Value() (driver.Value, error) {
	var zeroTime time.Time // 零时间
	var ti = time.Time(t)  // 转换为 time.Time 类型
	if ti.UnixNano() == zeroTime.UnixNano() {
		return nil, nil // 如果是零时间，则返回 nil
	}
	return ti, nil // 返回 time.Time 类型的时间和 nil 错误
}

// Scan 方法实现了数据库驱动的 Scanner 接口，用于从数据库扫描时间值。
func (t *Time) Scan(v interface{}) error {
	value, ok := v.(time.Time) // 断言接口为 time.Time 类型
	if ok {
		*t = Time(value) // 更新 Time 类型的值
		return nil       // 返回 nil 错误
	}
	return fmt.Errorf("can not convert %v to timestamp", v) // 返回错误信息
}
