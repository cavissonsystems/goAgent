package cavmongo

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/bson/bsonrw"
	"go.mongodb.org/mongo-driver/event"
	nd "goAgent"
	logger "goAgent/logger"
	"sync"
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
	bsonRegistry *bsoncodec.Registry

	mu sync.Mutex
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
	bt = ctx.Value("CavissonTx").(uint64)
	handle = nd.IP_db_callout_begin(bt, "db.mongodb.query", spanName)
	if len(event.Command) > 0 {
		sw := swPool.Get().(*bsonrw.SliceWriter)
		ejvw := extjPool.Get(sw, false /* non-canonical */, false /* don't escape HTML */)
		ec := bsoncodec.EncodeContext{Registry: c.bsonRegistry}
		if enc, err := bson.NewEncoderWithContext(ec, ejvw); err == nil {
			if err := enc.Encode(event.Command); err != nil {
				logger.ErrorPrint("Error : come from encode bson")
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
}

func (c *commandMonitor) failed(ctx context.Context, event *event.CommandFailedEvent) {
	c.finished(ctx, &event.CommandFinishedEvent)
}

func (c *commandMonitor) finished(ctx context.Context, event *event.CommandFinishedEvent) {

	c.mu.Lock()
	c.mu.Unlock()
	bt = ctx.Value("CavissonTx").(uint64)
	nd.IP_db_callout_end(bt, handle)
}

func collectionName(commandName string, command bson.Raw) (string, bool) {
	switch commandName {
	case
		"aggregate",
		"count",
		"distinct",
		"mapReduce",

		"geoNear",
		"geoSearch",

		"delete",
		"find",
		"findAndModify",
		"insert",
		"parallelCollectionScan",
		"update",

		"compact",
		"convertToCapped",
		"create",
		"createIndexes",
		"drop",
		"dropIndexes",
		"killCursors",
		"listIndexes",
		"reIndex",

		"collStats":

		collectionValue := command.Lookup(commandName)
		return collectionValue.StringValueOK()
	case "getMore":
		collectionValue := command.Lookup("collection")
		return collectionValue.StringValueOK()
	}
	return "", false
}

type Option func(*commandMonitor)
