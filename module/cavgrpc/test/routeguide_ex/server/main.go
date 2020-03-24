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
)
type server struct{}


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

        return &pb.Response{Result:result},nil
}

func (s *server) Multiply(ctx context.Context, request *pb.Request) (*pb.Response, error) {
        a,b := request.GetA(),request.GetB()
        result := a*b

        return &pb.Response{Result:result},nil
}

