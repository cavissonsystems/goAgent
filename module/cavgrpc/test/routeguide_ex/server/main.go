package main

import(
        "context"
        "net"
grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
        "google.golang.org/grpc"
        "google.golang.org/grpc/reflection"
     pb "google.golang.org/grpc/examples/routeguide_ex/routeguide_ex"
        "goAgent/module/cavgrpc"
     nd "goAgent"
logger "goAgent/logger"
)
type server struct{}
func m3(bt uint64) {
        nd.Method_entry(bt, "m3")
        logger.TracePrint("m1 called")    
        nd.Method_exit(bt, "m3")
}

func m4(bt uint64) {
        nd.Method_entry(bt, "m4")
        logger.TracePrint("m4 called")    
        nd.Method_exit(bt, "m4")
}



func main(){
        nd.Sdk_init()
        interceptors := []grpc.UnaryServerInterceptor{grpc_recovery.UnaryServerInterceptor()}
        serverOpts := []grpc.ServerOption{}

        interceptors = append(interceptors, cavgrpc.NewUnaryServerInterceptor())
        serverOpts = append(serverOpts, grpc_middleware.WithUnaryServerChain(interceptors...))

        srv :=grpc.NewServer(serverOpts...)

        pb.RegisterAddServiceServer(srv,&server{})
        listener ,err := net.Listen("tcp",":4040")
        if err != nil{
       panic(err)
       }

       reflection.Register(srv)

       if e := srv.Serve(listener); e != nil{

       panic(err)
      }

      nd.Sdk_free()

}

func (s *server) Add(ctx context.Context, request *pb.Request) (*pb.Response, error) {
        a,b := request.GetA(),request.GetB()
        result := a+b
        bt := ctx.Value("CavissonTx").(uint64)
        m3(bt)
        return &pb.Response{Result:result},nil
}

func (s *server) Multiply(ctx context.Context, request *pb.Request) (*pb.Response, error) {
        a,b := request.GetA(),request.GetB()
        result := a*b
        bt := ctx.Value("CavissonTx").(uint64)
        m4(bt)
        return &pb.Response{Result:result},nil
}

