# H2G

contracted but not simple

h2g is a plugins to test grpc service when the testing tool does not support grpc in iris

## How to use

Download in your project
> go get github.com/culionbear/h2g@latest

```golang
//Use H2G in iris
func main() {
    m := iris.New()
    //New H2G Manager
    g := h2g.New(&h2g.Config{
        Handler: map[string]h2g.Func{
            //server name : GrpcClient
            "live.Debug": func() interface{} {
                //New Grpc Connect
                conn, err := grpc.Dial(model.GrpcAddr, grpc.WithInsecure())
                if err != nil {
                    log.Fatalf("did not connect: %v", err)
                }
                //New Grpc Client with protobuf
                return pb.NewClient(conn)
            },
        },
        //get service name from url
        Service: "service",
        //get method name from url
        Method: "method",
    })
    //m.Post("/grpc/{service}/{method}", g.Service)
    m.Post(g.Path("/grpc"), g.Service)
    m.Run(iris.Addr("127.0.0.1:80"))
}
```
