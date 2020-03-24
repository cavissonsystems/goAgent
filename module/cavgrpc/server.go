package cavgrpc

import (
	_"strings"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	_"google.golang.org/grpc/codes"
	_"google.golang.org/grpc/metadata"
	_"google.golang.org/grpc/status"
        nd "goAgent"
)


func NewUnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp interface{}, err error) {

                unique_id :="1"
                bt_old := nd.BT_begin("server_name", "")
          
                new_ctx := nd.Updated_context(ctx, bt_old)
                
                bt := nd.Current_Transaction(new_ctx)
         
                defer nd.BT_end(bt)
                nd.BT_store(bt, unique_id)

		resp, err = handler(ctx, req)
		return resp, err
	}
}
