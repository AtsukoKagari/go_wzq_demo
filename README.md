# go_wzq_demo
## Demo介绍
一个简单的go websocket实例，五子棋对战的Demo，可以支持多组对战
## Demo用法
1.修改rule.js中的`websocket = new WebSocket("ws://localhost:8011/ws");`将localhost改成自己的ip

2.go run server.go

3.开启两个页面即可进行对战，也可以把文件挂载在nginx或者tomcat上
## Demo引用：
>前端代码以及后端核心逻辑：https://blog.csdn.net/xuhangsong/article/details/81356657
