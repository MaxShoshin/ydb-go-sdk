package generator

import "time"

type RowID = uint64

type Row struct {
	Hash             uint64     `gorm:"column:hash;primarykey;autoIncrement:false" xorm:"pk 'hash'"`
	ID               RowID      `gorm:"column:id;primarykey;autoIncrement:false" xorm:"pk 'id'"`
	PayloadStr       *string    `gorm:"column:payload_str" xorm:"'payload_str'"`
	PayloadDouble    *float64   `gorm:"column:payload_double" xorm:"'payload_double'"`
	PayloadTimestamp *time.Time `gorm:"column:payload_timestamp" xorm:"'payload_timestamp'"`
	PayloadHash      uint64     `gorm:"column:payload_hash" xorm:"'payload_hash'"`
}
