package models

type RowScanner interface {
	Scan(dest ...interface{}) error
}
