package cavgorm

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jinzhu/gorm"
	//"go.elastic.co/apm"

	nd "goAgent"
	"goAgent/module/cavsql"
)

const (
	cavContextKey = "cassionTX"
)

var handle uint64
var bt uint64

// WithContext returns a copy of db with ctx recorded for use by
// the callbacks registered via RegisterCallbacks.
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

// RegisterCallbacks registers callbacks on db for reporting spans
// to Elastic APM. This is called automatically by apmgorm.Open;
// it is provided for cases where a *gorm.DB is acquired by other
// means.
func RegisterCallbacks(db *gorm.DB) {
	registerCallbacks(db, cavsql.DSNInfo{})
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
		//span, ctx := apm.StartSpan(ctx, "", spanType)
		//if span.Dropped() {
		//	span.End()
		//	ctx = nil
		//}
		scope.Set(cavContextKey, ctx)
	}
}

//dbcallout end
func newAfterCallback(dsnInfo cavsql.DSNInfo) func(*gorm.Scope) {
	return func(scope *gorm.Scope) {
		//ctx, ok := scopeContext(scope)
		//	if !ok {
		//	return
	        //	span := apm.SpanFromContext(ctx)
		//	if span == nil {
		//	return
		//	}
		//span.Name = cavsql.QuerySignature(scope.SQL)
		//	span.Context.SetDestinationAddress(dsnInfo.Address, dsnInfo.Port)
		//	span.Context.SetDatabase(cav.DatabaseSpanContext{
		//	Instance:  dsnInfo.Database,
		//	Statement: scope.SQL,
		//Type:      "sql",
		//		User:      dsnInfo.User,
		//	})
		defer nd.IP_db_callout_end(bt, handle)

		// Capture errors, except for "record not found", which may be expected.
		for _, err := range scope.DB().GetErrors() {
			if gorm.IsRecordNotFoundError(err) || err == sql.ErrNoRows {
				continue
			}
			//		if e := cav.CaptureError(ctx, err); e != nil {
			//		e.Send()
			//	}
		}
	}
}
