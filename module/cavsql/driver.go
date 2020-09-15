package cavsql

import (
	"database/sql"
	"database/sql/driver"
	_ "github.com/go-sql-driver/mysql"
	"goAgent/internal/sqlutil"
	"sync"
)

const DriverPrefix = "cav/"

var (
	driversMu sync.RWMutex
	drivers   = make(map[string]*tracingDriver)
)

func Register(name string, driver driver.Driver) {
	driversMu.Lock()
	defer driversMu.Unlock()
	wrapped := newTracingDriver(driver)
	sql.Register(DriverPrefix+name, wrapped)
	drivers[name] = wrapped
}

func Open(driverName, dataSourceName string) (*sql.DB, error) {
	return sql.Open(DriverPrefix+driverName, dataSourceName)
}

func Wrap(driver driver.Driver, opts ...WrapOption) driver.Driver {
	return newTracingDriver(driver, opts...)
}

func newTracingDriver(driver driver.Driver, opts ...WrapOption) *tracingDriver {
	d := &tracingDriver{
		Driver: driver,
	}
	for _, opt := range opts {
		opt(d)
	}
	if d.driverName == "" {
		d.driverName = sqlutil.DriverName(driver)
	}
	if d.dsnParser == nil {
		d.dsnParser = genericDSNParser
	}
	return d
}

func DriverDSNParser(driverName string) DSNParserFunc {
	driversMu.RLock()
	driver := drivers[driverName]
	defer driversMu.RUnlock()

	if driver == nil {
		return genericDSNParser
	}
	return driver.dsnParser
}

type DSNInfo struct {
	Address  string
	Port     int
	Database string
	User     string
}

type DSNParserFunc func(dsn string) DSNInfo

func genericDSNParser(string) DSNInfo {
	return DSNInfo{}
}

type WrapOption func(*tracingDriver)

func WithDriverName(name string) WrapOption {
	return func(d *tracingDriver) {
		d.driverName = name
	}
}

type tracingDriver struct {
	driver.Driver
	driverName      string
	dsnParser       DSNParserFunc
	connectSpanType string
	execSpanType    string
	pingSpanType    string
	prepareSpanType string
	querySpanType   string
}

func (d *tracingDriver) Open(name string) (driver.Conn, error) {
	conn, err := d.Driver.Open(name)
	if err != nil {
		return nil, err
	}
	return newConn(conn, d, d.dsnParser(name)), nil
}
