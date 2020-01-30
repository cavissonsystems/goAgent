package cavsql

import (
	"database/sql/driver"
	"goAgent/internal/sqlutil"
	"sync"
)

//const DriverPrefix = "cav/"

var (
	driversMu sync.RWMutex
	drivers   = make(map[string]*tracingDriver)
)

/*func Register(name string, driver driver.Driver, opts ...WrapOption) {
	driversMu.Lock()
	defer driversMu.Unlock()

	wrapped := newTracingDriver(driver, opts...)
	sql.Register(DriverPrefix+name, wrapped)
	drivers[name] = wrapped
}*/

/*func Open(driverName, dataSourceName string) (*sql.DB, error) {
	return sql.Open(DriverPrefix+driverName, dataSourceName)
}*/

/*func Wrap(driver driver.Driver, opts ...WrapOption) driver.Driver {
	return newTracingDriver(driver, opts...)
}*/

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
	//if d.dsnParser == nil {
	//	d.dsnParser = genericDSNParser
	//}

	// store span types to avoid repeat allocations
	//d.connectSpanType = d.formatSpanType("connect")
	//d.pingSpanType = d.formatSpanType("ping")
	//d.prepareSpanType = d.formatSpanType("prepare")
	//d.querySpanType = d.formatSpanType("query")
	//d.execSpanType = d.formatSpanType("exec")
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

// DSNInfo contains information from a database-specific data source name.
type DSNInfo struct {
	// Address is the database server address specified by the DSN.
	Address string

	// Port is the database server port specified by the DSN.
	Port int

	// Database is the name of the specific database identified by the DSN.
	Database string

	// User is the username that the DSN specifies for authenticating the
	// database connection.
	User string
}

// DSNParserFunc is the type of a function that can be used for parsing a
// data source name, and returning the corresponding Info.
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

//func WithDSNParser(f DSNParserFunc) WrapOption {
//return func(d *tracingDriver) {
//	d.dsnParser = f
//}
//}

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

//func (d *tracingDriver) formatSpanType(suffix string) string {
//return fmt.Sprintf("db.%s.%s", d.driverName, suffix)
//}

func (d *tracingDriver) Open(name string) (driver.Conn, error) {
	conn, err := d.Driver.Open(name)
	if err != nil {
		return nil, err
	}
	return newConn(conn, d, d.dsnParser(name)), nil
}
