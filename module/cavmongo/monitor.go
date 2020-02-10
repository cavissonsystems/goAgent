package cavmongo

import (
        "log"
	"context"
	"sync"
        "fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/bson/bsonrw"
	"go.mongodb.org/mongo-driver/event"
        nd "goAgent"
)
var handle uint64
var bt uint64
var (
	extjPool = bsonrw.NewExtJSONValueWriterPool()
	swPool   = sync.Pool{
		New: func() interface{} {
			return &bsonrw.SliceWriter{}
		},
	}
)
// CommandMonitor returns a new event.CommandMonitor which will report a span
// for each command executed within a context containing a sampled transaction.
func CommandMonitor() *event.CommandMonitor {
	cm := commandMonitor{
		bsonRegistry: bson.DefaultRegistry,
	}
	return &event.CommandMonitor{
		Started:   cm.started,
		Succeeded: cm.succeeded,
		Failed:    cm.failed,
	}
}

type commandMonitor struct {
	// TODO(axw) record number of active commands and report as a
	// metric so users can, for example, identify unclosed cursors.
	bsonRegistry *bsoncodec.Registry

	mu    sync.Mutex
}

type commandKey struct {
	connectionID string
	requestID    int64
}

func (c *commandMonitor) started(ctx context.Context, event *event.CommandStartedEvent) {
	spanName := event.CommandName
	if collectionName, ok := collectionName(event.CommandName, event.Command); ok {
		spanName = collectionName + "." + spanName
	}
                fmt.Println(spanName)
                bt = ctx.Value("CavissonTx").(uint64)
                handle =  nd.IP_db_callout_begin(bt, "db.mongodb.query" , spanName)
                fmt.Println("started called")
	if len(event.Command) > 0 {
		// Encode the command as MongoDB Extended JSON
		// for the "statement" in database span context.
		sw := swPool.Get().(*bsonrw.SliceWriter)
		ejvw := extjPool.Get(sw, false /* non-canonical */, false /* don't escape HTML */)
		ec := bsoncodec.EncodeContext{Registry: c.bsonRegistry}
		if enc, err := bson.NewEncoderWithContext(ec, ejvw); err == nil {
			if err := enc.Encode(event.Command); err != nil {
                   log.Println("Error : come from encode bson")               
			}
		}
		*sw = (*sw)[:0]
		extjPool.Put(ejvw)
		swPool.Put(sw)
	}


	c.mu.Lock()
	c.mu.Unlock()
}

func (c *commandMonitor) succeeded(ctx context.Context, event *event.CommandSucceededEvent) {
	c.finished(ctx, &event.CommandFinishedEvent)
        fmt.Println("suceeded called")
}

func (c *commandMonitor) failed(ctx context.Context, event *event.CommandFailedEvent) {
	c.finished(ctx, &event.CommandFinishedEvent)
}

func (c *commandMonitor) finished(ctx context.Context, event *event.CommandFinishedEvent) {

	c.mu.Lock()
        fmt.Println("finished called")
	c.mu.Unlock()
        bt = ctx.Value("CavissonTx").(uint64)
        nd.IP_db_callout_end(bt , handle )
}

func collectionName(commandName string, command bson.Raw) (string, bool) {
	switch commandName {
	case
		// Aggregation Commands
		"aggregate",
		"count",
		"distinct",
		"mapReduce",

		// Geospatial Commands
		"geoNear",
		"geoSearch",

		// Query and Write Operation Commands
		"delete",
		"find",
		"findAndModify",
		"insert",
		"parallelCollectionScan",
		"update",

		// Administration Commands
		"compact",
		"convertToCapped",
		"create",
		"createIndexes",
		"drop",
		"dropIndexes",
		"killCursors",
		"listIndexes",
		"reIndex",

		// Diagnostic Commands
		"collStats":

		collectionValue := command.Lookup(commandName)
		return collectionValue.StringValueOK()
	case "getMore":
		collectionValue := command.Lookup("collection")
		return collectionValue.StringValueOK()
	}
	return "", false
}

// Option sets options for tracing MongoDB commands.
type Option func(*commandMonitor)
