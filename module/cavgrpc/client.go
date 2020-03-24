
package cavgrpc

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	_"google.golang.org/grpc/metadata"
        nd "goAgent"
)

func NewUnaryClientInterceptor() grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req, resp interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		
                unique_id :="1"
                name := method
                bt_old := nd.BT_begin(name, "")
                new_ctx := nd.Updated_context(ctx, bt_old)
                
                bt := nd.Current_Transaction(new_ctx)
                defer nd.BT_end(bt)
                nd.BT_store(bt, unique_id)

		return invoker(ctx, method, req, resp, cc, opts...)
	}
}
