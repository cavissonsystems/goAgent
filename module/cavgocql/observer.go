+
package cavgocql

import (
	"context"
       nd "goAgent"
	"github.com/gocql/gocql"

)

type Observer struct {
}

func NewObserver() *Observer {
	return &Observer{}
}

func (o *Observer) ObserveBatch(ctx context.Context, batch gocql.ObservedBatch) {

	for _, statement := range batch.Statements {
               bt := ctx.Value("CavissonTx").(uint64)
               db_handle :=  nd.IP_db_callout_begin(bt ,"db.cassendra",  querySignature(statement))
               defer nd.IP_db_callout_end(bt , db_handle)

	}

}

func (o *Observer) ObserveQuery(ctx context.Context, query gocql.ObservedQuery) {
               bt := ctx.Value("CavissonTx").(uint64)
               db_handle :=  nd.IP_db_callout_begin(bt ,"db.redis",  querySignature(query.Statement))
               defer nd.IP_db_callout_end(bt , db_handle)

}

