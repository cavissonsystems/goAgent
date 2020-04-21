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
        "time"
)
type server struct{}
func m3(bt uint64) {
        nd.Method_entry(bt, "a.b.m3")
        time.Sleep(2*time.Millisecond)
        logger.TracePrint("a.b.m1 called")    
        nd.Method_exit(bt, "a.b.m3")
}

func m4(bt uint64) {
        nd.Method_entry(bt, "a.b.m4")
        logger.TracePrint("a.b.m4 called")    
        time.Sleep(2*time.Millisecond)
        nd.Method_exit(bt, "a.b.m4")
}
func m5(bt uint64) {
        nd.Method_entry(bt, "a.b.m5")
        time.Sleep(2*time.Millisecond)
        logger.TracePrint("a.b.m5 called")    
        nd.Method_exit(bt, "a.b.m5")
}
func m6(bt uint64) {
        nd.Method_entry(bt, "a.b.m6")
        time.Sleep(2*time.Millisecond)
        logger.TracePrint("a.b.m6 called")    
        nd.Method_exit(bt, "a.b.m6")
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
        m5(bt)
        return &pb.Response{Result:result},nil
}

func (s *server) Multiply(ctx context.Context, request *pb.Request) (*pb.Response, error) {
        a,b := request.GetA(),request.GetB()
        result := a*b
        bt := ctx.Value("CavissonTx").(uint64)
        m4(bt)
        m6(bt)
        return &pb.Response{Result:result},nil
}

