package sql

import (
	"context"
	"database/sql"
	"github.com/any-lyu/go.library/tool"
	"time"

	"github.com/any-lyu/go.library/logs"
	"github.com/any-lyu/go.library/stat"
	"github.com/opentracing/opentracing-go"
)

const (
	defaultTimeout = 500 * time.Millisecond // 默认的超时时间
	minTimeout     = time.Millisecond       // 最小超时时间
)

// Row 对应 *sql.Row.
type Row interface {
	Scan(dest ...interface{}) error
}

// Rows 对应 *sql.Rows.
type Rows interface {
	Close() error
	Columns() ([]string, error)
	ColumnTypes() ([]*sql.ColumnType, error)
	Err() error
	Next() bool
	NextResultSet() bool
	Scan(dest ...interface{}) error
}

// Stmt 是数据库 Statement 接口, 可以通过这个接口操作数据库.
type Stmt interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) Row

	// WithTimeout 返回一个新的 Stmt, 这个新的 Stmt 的超时时间是 timeout.
	//
	// 请注意当前的 Stmt 的不受影响.
	WithTimeout(timeout time.Duration) Stmt
}

// DB 是数据库接口, 对应 *sql.DB.
type DB interface {
	Stmt
	BeginTx(ctx context.Context, opts *sql.TxOptions) (Tx, error)
	Close() error
	Conn(ctx context.Context) (Conn, error)
	PingContext(ctx context.Context) error
}

var defaultOptions = options{
	timeout:       defaultTimeout,
	enableTracing: true,
	stat:          stat.DB,
}

type options struct {
	timeout       time.Duration // sql 操作超时时间
	enableTracing bool          // 是否启用 tracing 功能, 默认开启
	stat          stat.Stat     //  监控
}

// Option 表示一些可选的参数.
type Option func(*options)

// WithTimeout 设置超时时间.
func WithTimeout(timeout time.Duration) Option {
	return func(o *options) {
		if timeout < minTimeout {
			timeout = minTimeout
		}
		o.timeout = timeout
	}
}

// WithTracing 启用 tracing 功能, 默认开启 tracing 功能.
func WithTracing() Option {
	return func(o *options) {
		o.enableTracing = true
	}
}

// WithoutTracing 禁用 tracing 功能, 默认开启 tracing 功能.
func WithoutTracing() Option {
	return func(o *options) {
		o.enableTracing = false
	}
}

// WithStat 设置自定义监控默认是 metrics.DBDefault
func WithStat(stat stat.Stat) Option {
	return func(o *options) {
		o.stat = stat
	}
}

// WithOutStat 关闭监控
func WithOutStat() Option {
	return func(o *options) {
		o.stat = nil
	}
}

// WrapDB 包装 *sql.DB 返回 DB 接口.
func WrapDB(db *sql.DB, opts ...Option) DB {
	if db == nil {
		return nil
	}
	innerOptions := defaultOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt(&innerOptions)
	}
	return &dbWrapper{db: db, options: innerOptions}
}

var _ DB = (*dbWrapper)(nil)

type dbWrapper struct {
	db      *sql.DB
	options options
}

func (db *dbWrapper) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	var span opentracing.Span
	if db.options.enableTracing {
		span = startTracingSpan(ctx, "exec_context", query)
	}

	logs.Info("database-exec-log", "query", query, "args", tool.D2S(args))
	if d, ok := ctx.Deadline(); !ok || d.After(time.Now().Add(db.options.timeout)) {
		var cancel context.CancelFunc
		ctx, cancel = context.WithCancel(ctx)
		timer := time.AfterFunc(db.options.timeout, cancel)
		defer timer.Stop()
	}
	now := time.Now()
	result, err := db.db.ExecContext(ctx, query, args...)
	statistics(db.options.stat, "sql:dbWrapper:ExecContext", now, err)
	if span != nil {
		finishTracingSpan(span, err)
	}

	return result, err
}

func (db *dbWrapper) QueryContext(ctx context.Context, query string, args ...interface{}) (Rows, error) {
	var span opentracing.Span
	if db.options.enableTracing {
		span = startTracingSpan(ctx, "query_context", query)
	}

	if d, ok := ctx.Deadline(); !ok || d.After(time.Now().Add(db.options.timeout)) {
		var cancel context.CancelFunc
		ctx, cancel = context.WithCancel(ctx)
		timer := time.AfterFunc(db.options.timeout, cancel)
		defer timer.Stop()
	}
	now := time.Now()
	result, err := db.db.QueryContext(ctx, query, args...)
	statistics(db.options.stat, "sql:dbWrapper:QueryContext", now, err)
	if span != nil {
		finishTracingSpan(span, err)
	}

	return result, err
}

func (db *dbWrapper) QueryRowContext(ctx context.Context, query string, args ...interface{}) Row {
	var span opentracing.Span
	if db.options.enableTracing {
		span = startTracingSpan(ctx, "query_row_context", query)
	}

	if d, ok := ctx.Deadline(); !ok || d.After(time.Now().Add(db.options.timeout)) {
		var cancel context.CancelFunc
		ctx, cancel = context.WithCancel(ctx)
		timer := time.AfterFunc(db.options.timeout, cancel)
		defer timer.Stop()
	}
	now := time.Now()
	result := db.db.QueryRowContext(ctx, query, args...)
	statistics(db.options.stat, "sql:dbWrapper:QueryRowContext", now, nil)
	if span != nil {
		finishTracingSpan(span, nil)
	}

	return result
}

func (db *dbWrapper) WithTimeout(timeout time.Duration) Stmt {
	if timeout < minTimeout {
		timeout = minTimeout
	}
	options := db.options
	options.timeout = timeout
	return &dbWrapper{db: db.db, options: options}
}

func (db *dbWrapper) BeginTx(ctx context.Context, opts *sql.TxOptions) (Tx, error) {
	tx, err := db.db.BeginTx(ctx, opts)
	if err != nil {
		return nil, err
	}
	return &txWrapper{tx: tx, options: db.options}, nil
}

func (db *dbWrapper) Close() error { return db.db.Close() }

func (db *dbWrapper) Conn(ctx context.Context) (Conn, error) {
	conn, err := db.db.Conn(ctx)
	if err != nil {
		return nil, err
	}
	return &connWrapper{conn: conn, options: db.options}, nil
}

func (db *dbWrapper) PingContext(ctx context.Context) error {
	if d, ok := ctx.Deadline(); !ok || d.After(time.Now().Add(db.options.timeout)) {
		var cancel context.CancelFunc
		ctx, cancel = context.WithCancel(ctx)
		timer := time.AfterFunc(db.options.timeout, cancel)
		defer timer.Stop()
	}
	return db.db.PingContext(ctx)
}

// Conn 是数据库连接接口, 对应 *sql.Conn.
type Conn interface {
	Stmt
	BeginTx(ctx context.Context, opts *sql.TxOptions) (Tx, error)
	Close() error
	PingContext(ctx context.Context) error
}

var _ Conn = (*connWrapper)(nil)

type connWrapper struct {
	conn    *sql.Conn
	options options
}

func (conn *connWrapper) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	var span opentracing.Span
	if conn.options.enableTracing {
		span = startTracingSpan(ctx, "exec_context", query)
	}

	logs.Info("database-exec-log", "query", query, "args", tool.D2S(args))
	if d, ok := ctx.Deadline(); !ok || d.After(time.Now().Add(conn.options.timeout)) {
		var cancel context.CancelFunc
		ctx, cancel = context.WithCancel(ctx)
		timer := time.AfterFunc(conn.options.timeout, cancel)
		defer timer.Stop()
	}
	now := time.Now()
	result, err := conn.conn.ExecContext(ctx, query, args...)
	statistics(conn.options.stat, "sql:connWrapper:ExecContext", now, err)
	if span != nil {
		finishTracingSpan(span, err)
	}

	return result, err
}

func (conn *connWrapper) QueryContext(ctx context.Context, query string, args ...interface{}) (Rows, error) {
	var span opentracing.Span
	if conn.options.enableTracing {
		span = startTracingSpan(ctx, "query_context", query)
	}

	if d, ok := ctx.Deadline(); !ok || d.After(time.Now().Add(conn.options.timeout)) {
		var cancel context.CancelFunc
		ctx, cancel = context.WithCancel(ctx)
		timer := time.AfterFunc(conn.options.timeout, cancel)
		defer timer.Stop()
	}
	now := time.Now()
	result, err := conn.conn.QueryContext(ctx, query, args...)
	statistics(conn.options.stat, "sql:connWrapper:QueryContext", now, err)
	if span != nil {
		finishTracingSpan(span, err)
	}

	return result, err
}

func (conn *connWrapper) QueryRowContext(ctx context.Context, query string, args ...interface{}) Row {
	var span opentracing.Span
	if conn.options.enableTracing {
		span = startTracingSpan(ctx, "query_row_context", query)
	}

	if d, ok := ctx.Deadline(); !ok || d.After(time.Now().Add(conn.options.timeout)) {
		var cancel context.CancelFunc
		ctx, cancel = context.WithCancel(ctx)
		timer := time.AfterFunc(conn.options.timeout, cancel)
		defer timer.Stop()
	}
	now := time.Now()
	result := conn.conn.QueryRowContext(ctx, query, args...)
	statistics(conn.options.stat, "sql:connWrapper:QueryRowContext", now, nil)
	if span != nil {
		finishTracingSpan(span, nil)
	}

	return result
}

func (conn *connWrapper) WithTimeout(timeout time.Duration) Stmt {
	if timeout < minTimeout {
		timeout = minTimeout
	}
	options := conn.options
	options.timeout = timeout
	return &connWrapper{conn: conn.conn, options: options}
}

func (conn *connWrapper) BeginTx(ctx context.Context, opts *sql.TxOptions) (Tx, error) {
	tx, err := conn.conn.BeginTx(ctx, opts)
	if err != nil {
		return nil, err
	}
	return &txWrapper{tx: tx, options: conn.options}, nil
}
func (conn *connWrapper) Close() error { return conn.conn.Close() }
func (conn *connWrapper) PingContext(ctx context.Context) error {
	if d, ok := ctx.Deadline(); !ok || d.After(time.Now().Add(conn.options.timeout)) {
		var cancel context.CancelFunc
		ctx, cancel = context.WithCancel(ctx)
		timer := time.AfterFunc(conn.options.timeout, cancel)
		defer timer.Stop()
	}
	return conn.conn.PingContext(ctx)
}

// Tx 是数据库事务接口, 对应 *sql.Tx.
type Tx interface {
	Stmt
	Commit() error
	Rollback() error
}

var _ Tx = (*txWrapper)(nil)

type txWrapper struct {
	tx      *sql.Tx
	options options
}

func (tx *txWrapper) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	var span opentracing.Span
	if tx.options.enableTracing {
		span = startTracingSpan(ctx, "exec_context", query)
	}

	logs.Info("database-exec-log", "query", query, "args", tool.D2S(args))
	if d, ok := ctx.Deadline(); !ok || d.After(time.Now().Add(tx.options.timeout)) {
		var cancel context.CancelFunc
		ctx, cancel = context.WithCancel(ctx)
		timer := time.AfterFunc(tx.options.timeout, cancel)
		defer timer.Stop()
	}
	now := time.Now()
	result, err := tx.tx.ExecContext(ctx, query, args...)
	statistics(tx.options.stat, "sql:txWrapper:ExecContext", now, err)
	if span != nil {
		finishTracingSpan(span, err)
	}

	return result, err
}

func (tx *txWrapper) QueryContext(ctx context.Context, query string, args ...interface{}) (Rows, error) {
	var span opentracing.Span
	if tx.options.enableTracing {
		span = startTracingSpan(ctx, "query_context", query)
	}

	if d, ok := ctx.Deadline(); !ok || d.After(time.Now().Add(tx.options.timeout)) {
		var cancel context.CancelFunc
		ctx, cancel = context.WithCancel(ctx)
		timer := time.AfterFunc(tx.options.timeout, cancel)
		defer timer.Stop()
	}
	now := time.Now()
	result, err := tx.tx.QueryContext(ctx, query, args...)
	statistics(tx.options.stat, "sql:txWrapper:QueryContext", now, err)
	if span != nil {
		finishTracingSpan(span, err)
	}

	return result, err
}

func (tx *txWrapper) QueryRowContext(ctx context.Context, query string, args ...interface{}) Row {
	var span opentracing.Span
	if tx.options.enableTracing {
		span = startTracingSpan(ctx, "query_row_context", query)
	}

	if d, ok := ctx.Deadline(); !ok || d.After(time.Now().Add(tx.options.timeout)) {
		var cancel context.CancelFunc
		ctx, cancel = context.WithCancel(ctx)
		timer := time.AfterFunc(tx.options.timeout, cancel)
		defer timer.Stop()
	}
	now := time.Now()
	result := tx.tx.QueryRowContext(ctx, query, args...)
	statistics(tx.options.stat, "sql:connWrapper:QueryRowContext", now, nil)
	if span != nil {
		finishTracingSpan(span, nil)
	}

	return result
}

func (tx *txWrapper) WithTimeout(timeout time.Duration) Stmt {
	if timeout < minTimeout {
		timeout = minTimeout
	}
	options := tx.options
	options.timeout = timeout
	return &txWrapper{tx: tx.tx, options: options}
}

func (tx *txWrapper) Commit() error   { return tx.tx.Commit() }
func (tx *txWrapper) Rollback() error { return tx.tx.Rollback() }
