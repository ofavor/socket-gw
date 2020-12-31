# echo example

## TCP Sampler.jmx

Jmeter test plan file

1. start server

```
go run main.go
```

2. start jmeter

## pprof

```
http://127.0.0.1:9990/debug/pprof/
```

参考：

* https://www.freecodecamp.org/news/how-i-investigated-memory-leaks-in-go-using-pprof-on-a-large-codebase-4bec4325e192/
* https://segmentfault.com/a/1190000019222661

获取不同时间点的heap数据

```
curl http://127.0.0.1:9990/debug/pprof/heap > heap.out
... jmeter test run
curl http://127.0.0.1:9990/debug/pprof/heap > heap1.out
```

查看使用量top
```
go tool pprof -base heap_2.out
或
go tool pprof -base heap_1.out heap_2.out

(pprof) top
Showing nodes accounting for 885.89kB, 40.50% of 2187.14kB total
Showing top 10 nodes out of 23
      flat  flat%   sum%        cum   cum%
 1024.48kB 46.84% 46.84%  1024.48kB 46.84%  github.com/ofavor/socket-gw/session.NewSession
 -650.62kB 29.75% 17.09%  -650.62kB 29.75%  compress/flate.(*compressor).init
  512.03kB 23.41% 40.50%   512.03kB 23.41%  syscall.anyToSockaddr
         0     0% 40.50%  -650.62kB 29.75%  compress/flate.NewWriter
         0     0% 40.50%  -650.62kB 29.75%  compress/gzip.(*Writer).Write
         0     0% 40.50%  1536.52kB 70.25%  github.com/ofavor/socket-gw/server.(*tcpServer).Run.func1
         0     0% 40.50%   512.03kB 23.41%  internal/poll.(*FD).Accept
         0     0% 40.50%   512.03kB 23.41%  internal/poll.accept
         0     0% 40.50%   512.03kB 23.41%  net.(*TCPListener).Accept
         0     0% 40.50%   512.03kB 23.41%  net.(*TCPListener).accept

(pprof) list session.NewSession
Total: 2.14MB
ROUTINE ======================== github.com/ofavor/socket-gw/session.NewSession in /Users/henry.sha/Workspaces/agent/socket-gw/session/session.go
       1MB        1MB (flat, cum) 46.84% of Total
         .          .     67:		conn:      conn,
         .          .     68:		transport: t,
         .          .     69:		handler:   handler,
         .          .     70:		needAuth:  auth,
         .          .     71:		sendCh:    make(chan *transport.Packet, SendBufferSize),
       1MB        1MB     72:		recvCh:    make(chan *transport.Packet, RecvBufferSize),
         .          .     73:		stopCh:    make(chan interface{}),
         .          .     74:		status:    StatusInitial,
         .          .     75:		meta:      make(map[string]string),
         .          .     76:	}
         .          .     77:}

```

查看图形展示

```
go tool pprof -web heap_2.out
或
go tool pprof -web -base heap_1.out heap_2.out
```



```
go tool pprof -web -base heap_1.out heap_2.out
```