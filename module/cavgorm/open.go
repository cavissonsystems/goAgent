package cavgorm

import (
	"goAgent/module/cavsql"

	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

func Open(dialect string, args ...interface{}) (*gorm.DB, error) {
	var driverName, dsn string
	switch len(args) {
	case 1:
		switch arg0 := args[0].(type) {
		case string:
			driverName = dialect
			dsn = arg0
		}
	case 2:
		driverName, _ = args[0].(string)
		dsn, _ = args[1].(string)
	}
	db, err := gorm.Open(dialect, args...)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	registerCallbacks(db, cavsql.DriverDSNParser(driverName)(dsn))
	return db, nil
}
