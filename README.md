## 简介

** gg **，没想好名字，随便起的，主要用途是基于 MySQL 数据库增量日志解析，提供增量数据订阅和消费

当前的 canal 支持源端 MySQL 版本包括： 
- 5.6.x
- 5.7.x
- 8.0.x
- 阿里云RDS

我们基于「gg」的业务包括
- 数据库实时备份
- 业务 cache 处理
- 业务统计处理
- 推送消息

## 客户端
- nats: [https://nats.io/](https://nats.io/)
- nats_stream [https://nats.io/](https://nats.io/)
- rabbitmq [https://www.rabbitmq.com/](https://www.rabbitmq.com/)

## 重要版本更新说明
- 支持docker 详情参考 Dockerfile

## 问题反馈
- 报告 issue: [github issues](https://github.com/lvxin0315/gg/issues)