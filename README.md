# Log

该模块在Logrus基础上进行二次封装，支持打印行号，支持格式化日志，支持本地生成日志文件和以JSON格式远程上报日志。

## 环境配置

- 本项目编译之前，需用户自行安装`Golang`和`glide`工具；
- 在编译时，本项目会先从`github`上拉去以来的项目，编译时间视具体网络环境而定；
- 要求`linux`或`macos`系统即可。

## 编译

- 第一次下载程序后，先`make update`，下载第三方依赖库文件；
- 最后再`make`， 编译程序（如果第三方库安装vendor目录存在，则可直接make，无需操作install和update）;


## 引用

- import "github.com/dongjialong2006/log"

## 示例
{
	package main
	
	import "github.com/dongjialong2006/log"
	
	func main() {
		logger := log.New("main")
		logger.Debug("test")
	}
}

time="2018-07-25 15:37:05.2244" level=debug msg=test file="/Users/dongchaofeng/go/src/log/log_test.go" line=19 model=test