用于存放cwgo生成的rpc客户端代码
cwgo client -I ../../idl/ --type RPC --module xxx/xxx/xxx/xxx  --service demo_proto --idl ../../idl/echo.proto
- server 表示生成服务代码 client 表示客户端代码
- -I ../../idl/ 表示idl文件搜索路径
- --type RPC 表示用的是kitex框架生成RPC代码 --type HTTP 表示用的是hertz框架生成HTTP代码
- --service 表示服务名
- --idl ../../idl/echo.proto  idl文件位置
  --module xxx/xxx/xxx/xxx 如果项目文件不在gowork下需要指定 如果在则忽略


示例：
    1.先生成客户端代码到rpc_gen（注意先进入到rpc_gen文件夹内）
	cwgo client --type RPC --service cart --I ../idl --idl ../idl/xxx.proto
    2.再生成服务端代码（注意切换到app下你的目录）
	cwgo server --type RPC --service cart --I ../idl --idl ../idl/xxx.proto --pass "-use ../../rpc_gen/kitex_gen"
	*注意： --pass "-use ../../rpc_gen/kitex_gen" --pass表示传递参数给底层工具 -use表示使用../../rpc_gen/kitex_gen路径下的客户端代码不再另外生成 
