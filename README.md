# Git to DailyReport

将Git日志通过大模型自动转换成每日工作报告

## 特点

* 读取指定GIT仓库的指定COMMIT，上传大模型解析成工作报告
* 支持阿里云大模型

<img src="https://s21.ax1x.com/2024/03/20/pFWHH0I.jpg" alt="喵喵喵" width="300" height="200">

## 开始使用
Windows

    $ bin\dailyreport.exe
Linux (需要自编译)

    $ ./bin/dailyreport
* 开通阿里大模型 https://dashscope.aliyun.com
* 配置文件 需要将`config.sample.yaml`重命名为`config.yaml`，并填写配置
* 确保`config.yaml`与`dailyreport.exe`在同一个目录下

## Roadmap

* 无人工干预地自动读取Git日志，根据commit时间、作者自动选取
* 支持超大规模压缩日志上传
* 更自由地配置化，支持更多的大模型
* 可以从多个Git库多个分支拉日志
* 更加傻瓜智能化，力求一键全自动处理

## License

© Zhouyu, 2024

Released under the [MIT License](https://github.com/go-gorm/gorm/blob/master/LICENSE)