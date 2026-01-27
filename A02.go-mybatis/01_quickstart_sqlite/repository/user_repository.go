package repository

import (
	"quickstart_sqlite/model"
)

type UserMapper struct {
	SelectById func(id int64) (model.User, error) `mapperParams:"id"`
	Insert     func(params map[string]interface{}) (int64, error)
	SelectAll  func() ([]model.User, error)
}