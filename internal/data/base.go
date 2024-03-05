package data

import "time"

type BaseFields struct {
	Id        int64     `gorm:"primaryKey;autoIncrement;type:int;comment:主键id"`
	CreatedAt time.Time `gorm:"not null;type:timestamp;default:CURRENT_TIMESTAMP;comment:创建时间"`
	UpdatedAt time.Time `gorm:"not null;type:timestamp;default:CURRENT_TIMESTAMP;comment:修改时间"`
}
