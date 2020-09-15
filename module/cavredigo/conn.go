package cavredigo

import (
	"context"
	"github.com/gomodule/redigo/redis"
	nd "goAgent"
	"strings"
	"time"
)

type Conn interface {
	redis.Conn

	WithContext(ctx context.Context) Conn
}

func Wrap(conn redis.Conn) Conn {
	ctx := context.Background()
	if cwt, ok := conn.(redis.ConnWithTimeout); ok {
		return contextConnWithTimeout{ConnWithTimeout: cwt, ctx: ctx}
	}
	return contextConn{Conn: conn, ctx: ctx}
}

type contextConnWithTimeout struct {
	redis.ConnWithTimeout
	ctx context.Context
}

func (c contextConnWithTimeout) WithContext(ctx context.Context) Conn {
	c.ctx = ctx
	return c
}

func (c contextConnWithTimeout) Do(commandName string, args ...interface{}) (reply interface{}, err error) {
	return Do(c.ctx, c.ConnWithTimeout, commandName, args...)
}

func (c contextConnWithTimeout) DoWithTimeout(timeout time.Duration, commandName string, args ...interface{}) (reply interface{}, err error) {
	return DoWithTimeout(c.ctx, c.ConnWithTimeout, timeout, commandName, args...)
}

type contextConn struct {
	redis.Conn
	ctx context.Context
}

func (c contextConn) WithContext(ctx context.Context) Conn {
	c.ctx = ctx
	return c
}

func (c contextConn) Do(commandName string, args ...interface{}) (reply interface{}, err error) {
	return Do(c.ctx, c.Conn, commandName, args...)
}

func Do(ctx context.Context, conn redis.Conn, commandName string, args ...interface{}) (interface{}, error) {
	spanName := strings.ToUpper(commandName)
	if spanName == "" {
		spanName = "(flush pipeline)"
	}
	bt := ctx.Value("CavissonTx").(uint64)
	db_handle := nd.IP_db_callout_begin(bt, "db.redis", spanName)
	defer nd.IP_db_callout_end(bt, db_handle)
	return conn.Do(commandName, args...)
}

func DoWithTimeout(ctx context.Context, conn redis.Conn, timeout time.Duration, commandName string, args ...interface{}) (interface{}, error) {
	spanName := strings.ToUpper(commandName)
	if spanName == "" {
		spanName = "(flush pipeline)"
	}
	bt := ctx.Value("CavissonTx").(uint64)
	db_handle := nd.IP_db_callout_begin(bt, "db.redis", spanName)
	defer nd.IP_db_callout_end(bt, db_handle)
	return redis.DoWithTimeout(conn, timeout, commandName, args...)
}
