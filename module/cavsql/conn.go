package cavsql

import (
	"context"
	"database/sql/driver"
	"errors"
         "fmt"
	nd "goAgent"
       logger "goAgent/logger"
)

func newConn(in driver.Conn, d *tracingDriver, dsnInfo DSNInfo) driver.Conn {
	conn := &conn{Conn: in, driver: d}
	conn.dsnInfo = dsnInfo
	conn.namedValueChecker, _ = in.(namedValueChecker)
	conn.pinger, _ = in.(driver.Pinger)
	conn.queryer, _ = in.(driver.Queryer)
	conn.queryerContext, _ = in.(driver.QueryerContext)
	conn.connPrepareContext, _ = in.(driver.ConnPrepareContext)
	conn.execer, _ = in.(driver.Execer)
	conn.execerContext, _ = in.(driver.ExecerContext)
	conn.connBeginTx, _ = in.(driver.ConnBeginTx)
	if in, ok := in.(driver.ConnBeginTx); ok {
		return &connBeginTx{conn, in}
	}
	return conn
}

type conn struct {
	driver.Conn
	driver  *tracingDriver
	dsnInfo DSNInfo

	namedValueChecker  namedValueChecker
	pinger             driver.Pinger
	queryer            driver.Queryer
	queryerContext     driver.QueryerContext
	connPrepareContext driver.ConnPrepareContext
	execer             driver.Execer
	execerContext      driver.ExecerContext
	connBeginTx        driver.ConnBeginTx
}

func namedValueToValue(named []driver.NamedValue) ([]driver.Value, error) {
	dargs := make([]driver.Value, len(named))
	for n, param := range named {
		if len(param.Name) > 0 {
		      return nil, errors.New("sql: driver does not support the use of Named Parameters")
		}
		dargs[n] = param.Value
	}
	return dargs, nil
}

type namedValueChecker interface {
     CheckNamedValue(*driver.NamedValue) error
}

func checkNamedValue(nv *driver.NamedValue, next namedValueChecker) error {
	if next != nil {
	   return next.CheckNamedValue(nv)
	}
	return driver.ErrSkip
}

func (c *conn) Ping(ctx context.Context) (resultError error) {
      if c.pinger == nil {
	    logger.ErrorPrint("Error : request not found")
	}
	return c.pinger.Ping(ctx)
}

func (c *conn) QueryContext(ctx context.Context, query string, args []driver.NamedValue) (_ driver.Rows, resultError error) {
       fmt.Println("query context")
	if c.queryerContext == nil && c.queryer == nil {
	  return nil, driver.ErrSkip
	}
	bt := ctx.Value("CavissonTx").(uint64)
        fmt.Println(bt)
        handle := nd.IP_db_callout_begin(bt, "db_host_1", "query")
        fmt.Println(handle)
        defer nd.IP_db_callout_end(bt, handle)

	if c.queryerContext != nil {
	   return c.queryerContext.QueryContext(ctx, query, args)
	}

	dargs, err := namedValueToValue(args)
	if err != nil {
	   return nil, err
	}
	select {
	default:
	case <-ctx.Done():
	     return nil, ctx.Err()
	}
	return c.queryer.Query(query, dargs)
}
func (c *conn) PrepareContext(ctx context.Context, query string) (_ driver.Stmt, resultError error) {
	var stmt driver.Stmt
	var err error
	if c.connPrepareContext != nil {
	   stmt, err = c.connPrepareContext.PrepareContext(ctx, query)
	} else {
		stmt, err = c.Prepare(query)
		if err == nil {
			select {
			default:
			case <-ctx.Done():
				stmt.Close()
				return nil, ctx.Err()
			}
		}
	}
	if stmt != nil {
	   stmt = newStmt(stmt, c, query)
	}
	return stmt, err
}

func (c *conn) ExecContext(ctx context.Context, query string, args []driver.NamedValue) (result driver.Result, resultError error) {
	if c.execerContext == nil && c.execer == nil {
	   return nil, driver.ErrSkip
	}
	if c.execerContext != nil {
	   return c.execerContext.ExecContext(ctx, query, args)
	}
	dargs, err := namedValueToValue(args)
	if err != nil {
	   return nil, err
	}
	select {
	default:
	case <-ctx.Done():
	     return nil, ctx.Err()
	}
	return c.execer.Exec(query, dargs)
}

func (*conn) Exec(query string, args []driver.Value) (driver.Result, error) {
     return nil, errors.New("Exec should never be called")
}

func (c *conn) CheckNamedValue(nv *driver.NamedValue) error {
     return checkNamedValue(nv, c.namedValueChecker)
}

type connBeginTx struct {
	*conn
	connBeginTx driver.ConnBeginTx
}

func (c *connBeginTx) BeginTx(ctx context.Context, opts driver.TxOptions) (driver.Tx, error) {
     return c.connBeginTx.BeginTx(ctx, opts)
}
