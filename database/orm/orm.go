package orm

import (
	"time"

	_ "github.com/go-sql-driver/mysql" // justifying
	"github.com/go-xorm/xorm"
	"xorm.io/core"

	"github.com/any-lyu/go.library/errors"
	"github.com/any-lyu/go.library/logs"
	xtime "github.com/any-lyu/go.library/time"
)

// Config database config.
type Config struct {
	DSN         string         // data source name.
	Active      int            // pool
	Idle        int            // pool
	IdleTimeout xtime.Duration // connect max life time.
}

type ormLog struct {
	showSQL bool
}

func (l ormLog) Debug(v ...interface{}) {
	logs.Debug(v)
}

func (l ormLog) Debugf(format string, v ...interface{}) {
	logs.Debug(format, v...)
}

func (l ormLog) Error(v ...interface{}) {
	logs.Error(v)
}

func (l ormLog) Errorf(format string, v ...interface{}) {
	logs.Error(format, v...)
}

func (l ormLog) Info(v ...interface{}) {
	logs.Info(v)
}

func (l ormLog) Infof(format string, v ...interface{}) {
	logs.Info(format, v...)
}

func (l ormLog) Warn(v ...interface{}) {
	logs.Warn(v)
}

func (l ormLog) Warnf(format string, v ...interface{}) {
	logs.Warn(format, v...)
}

func (l ormLog) Level() core.LogLevel {
	return core.LogLevel(logs.GetLevel())
}

func (l ormLog) SetLevel(level core.LogLevel) {
	logs.Info("ormLog-SetLevel-Invalid")
}

func (l ormLog) ShowSQL(show ...bool) {
	if len(show) == 0 {
		l.showSQL = true
	} else {
		l.showSQL = show[0]
	}
}

func (l ormLog) IsShowSQL() bool {
	return l.showSQL
}

func init() {
	xorm.ErrNotExist = errors.ErrNotFount
}

// NewMySQL new db and retry connection when has error.
func NewMySQL(c *Config) (db *xorm.Engine) {
	db, err := xorm.NewEngine("database", c.DSN)
	if err != nil {
		logs.Error("db dsn(%s) error: %v", c.DSN, err)
		panic(err)
	}
	db.DB().SetMaxIdleConns(c.Idle)
	db.DB().SetMaxOpenConns(c.Active)
	db.DB().SetConnMaxLifetime(time.Duration(c.IdleTimeout) / time.Second)
	db.SetLogger(ormLog{})
	return
}
