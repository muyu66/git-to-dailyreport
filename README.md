# Git to DailyReport

将Git日志通过大模型自动转换成每日工作报告

## 特点

* **所码即所报**：无人工干预地自动读取Git日志，根据commit时间、作者自动选取，上传大模型解析成工作报告
* **多仓库报告聚合**：支持多个仓库的工作量自动聚合，自动加权均衡上报
* **多模型适配**：暂时只支持通义千问大模型、通义千问开源大模型
* **多仓库**：自动读取指定多个目录下的所有仓库
* 自动/手动指定Git Author，防止工作日志错领
* 支持切换日报详细程度模式 (简单/详细)，简单模式可以节省AI TOKEN消费

![](https://z1.wzznft.com/i/2024/03/21/iu9yyz.gif "to boss")

## 开始使用

Windows

    执行 dailyreport.exe

* 开通阿里大模型 https://dashscope.aliyun.com
* 配置文件 需要将`config.sample.yaml`重命名为`config.yaml`，并填写配置
* 确保`config.yaml`与`dailyreport.exe`在同一个目录下

## 效果展示

![](https://free.wzznft.com/i/2024/03/21/p1idmz.png)

## 编译 Win64

    .\build.bat # 编译脚本
    .\bin\dailyreport.exe # 开始使用

## Roadmap

* 将工作区、Stage区纳入上报的工作范围
* 支持超大规模压缩日志上传
* 可视化配置
* 更多大模型支持 (百度系、Meta、Google等)
* 更加傻瓜智能化，力求一键全自动处理
* 多种结果输出方式 (文本/终端/Webhook)
* 兼容Linux/MACOS
* 将工作范围简要报告给使用者
* 周报
* 应对加班等可能跨天的日报
* 日报的天数可自定义

## Contributors

[Thank you](https://github.com/muyu66/git-to-dailyreport/graphs/contributors) for contributing to the Anti-Boss!

## License

© Zhouyu, 2024

Released under the [MIT License](https://github.com/muyu66/git-to-dailyreport/blob/master/LICENSE)