package cavgorm

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/jinzhu/gorm"
	nd "goAgent"
	"goAgent/module/cavsql"
)

const (
	cavContextKey = "cassionTX"
)

var handle uint64
var bt uint64

func WithContext(ctx context.Context, db *gorm.DB) *gorm.DB {
	return db.Set(cavContextKey, ctx)
}

func scopeContext(scope *gorm.Scope) (context.Context, bool) {
	value, ok := scope.Get(cavContextKey)
	if !ok {
		return nil, false
	}
	ctx, _ := value.(context.Context)
	return ctx, ctx != nil
}

func RegisterCallbacks(db *gorm.DB) {
	registerCallbacks(db, cavsql.DSNInfo{})
}

<<<<<<< HEAD
=======
/*type DSNInfo struct {
	// Address is the database server address specified by the DSN.
	Address string

	// Port is the database server port specified by the DSN.
	Port int

	// Database is the name of the specific database identified by the DSN.
	Database string

	// User is the username that the DSN specifies for authenticating the
	// database connection.
	User string
}*/
>>>>>>> 83364ee12299b7fc92d14c26331a21466b5ef90a

func registerCallbacks(db *gorm.DB, dsnInfo cavsql.DSNInfo) {
	driverName := db.Dialect().GetName()
	switch driverName {
	case "postgres":
		driverName = "postgresql"
	}
	spanTypePrefix := fmt.Sprintf("db.%s.", driverName)
	querySpanType := spanTypePrefix + "query"
	execSpanType := spanTypePrefix + "exec"

	type params struct {
		spanType  string
		processor func() *gorm.CallbackProcessor
	}
	callbacks := map[string]params{
		"gorm:create": {
			spanType:  execSpanType,
			processor: func() *gorm.CallbackProcessor { return db.Callback().Create() },
		},
		"gorm:delete": {
			spanType:  execSpanType,
			processor: func() *gorm.CallbackProcessor { return db.Callback().Delete() },
		},
		"gorm:query": {
			spanType:  querySpanType,
			processor: func() *gorm.CallbackProcessor { return db.Callback().Query() },
		},
		"gorm:update": {
			spanType:  execSpanType,
			processor: func() *gorm.CallbackProcessor { return db.Callback().Update() },
		},
		"gorm:row_query": {
			spanType:  querySpanType,
			processor: func() *gorm.CallbackProcessor { return db.Callback().RowQuery() },
		},
	}
	for name, params := range callbacks {
		const callbackPrefix = "cavisson"
		params.processor().Before(name).Register(
			fmt.Sprintf("%s:before:%s", callbackPrefix, name),
			newBeforeCallback(params.spanType),
		)
		params.processor().After(name).Register(
			fmt.Sprintf("%s:after:%s", callbackPrefix, name),
			newAfterCallback(dsnInfo),
		)
	}

}

//dbcallout begin
func newBeforeCallback(spanType string) func(*gorm.Scope) {
	return func(scope *gorm.Scope) {
		ctx, ok := scopeContext(scope)
		if !ok {
			return
		}
		bt = ctx.Value("CavissonTx").(uint64)
		handle = nd.IP_db_callout_begin(bt, "db_host", spanType)

		scope.Set(cavContextKey, ctx)
	}
}

//dbcallout end
func newAfterCallback(dsnInfo cavsql.DSNInfo) func(*gorm.Scope) {
	return func(scope *gorm.Scope) {
		
		defer nd.IP_db_callout_end(bt, handle)

		for _, err := range scope.DB().GetErrors() {
			if gorm.IsRecordNotFoundError(err) || err == sql.ErrNoRows {
				continue
			}
		
		}
	}
}
