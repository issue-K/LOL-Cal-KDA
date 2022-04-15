### 优点

可以基于此代码进行二次开发。

尽量把很多代码都化简了,而且没有写日志之类的模块,都是单纯的逻辑,最后发送到lol聊天框内的是json数据。

所以有很多地方可以优化, 也可以做GUI界面, 很方便学习的一个项目.

这里有lcu api, 就是lol客户端接收的所有请求: [link](http://www.mingweisamuel.com/lcu-schema/tool/#/Plugin%20lol-champ-select)


### 使用方法

①.在后台运行程序

②.打开英雄联盟客户端

③.随意开一把游戏

④.可以在聊天框内收到队友近期的KDA战绩


核心是利用wmi获取lol client进程的端口号和token值

进入英雄选择界面后, 计算己方队友KDA，并发送到聊天界面.
