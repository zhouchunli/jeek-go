package main

import (
	"context"
	"database/sql"
	"github.com/pkg/errors"
)

type demoDao struct {
	db *sql.DB
}

func NewDemoDao() *demoDao {
	return &demoDao{

	}
}
//对查询的结果不关心的
func (t *demoDao) demoIgnore(ctx context.Context, uid int64) error {
	var info interface{}
	var err error
	sqlStr := "select * from `user` where uid = ? and status = 1 limit 1"
	_, e := t.db.Query(sqlStr, &info, uid)
	if e == sql.ErrNoRows {
		return nil
	} else if e != nil {
		err = errors.Wrapf(e, "demoDao: get err in demoIgnore")
		return err
	}
	return nil
}


//关心结果不存在，比如配置信息
func (t *demoDao) demoErr(ctx context.Context, uid int64) error {
	var info interface{}
	var err error
	sqlStr := "select * from `config` where uid = ? and status = 1 limit 1"
	_, e := t.db.Query(sqlStr, &info, uid)
	if e != nil {
		err = errors.Wrapf(e, "demoDao: get err in demoErr")
		return err
	}
	return nil
}

