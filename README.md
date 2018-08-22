一些见闻和一些思考
---

#### 典型的serverless(lambda)的场景

1. WEB应用：某WEB网站在用户注册成功后，会发一封欢迎邮件，通过函数计算把邮件内容定制成模板，每次触发，每次执行都是幂等无状态。
2. 图片处理：基于对象存储的事件触发，当用户上传的图片转入到某Bucket中后，自动触发函数岁图片进行可定制化处理
3. 定时任务
4. 音频转换文字处理：当用户通过语音来发出某些指令的时候，可以通过函数计算来触发阿里云的ET公开API获取到音频转换成文字的方式。



#### 目前框架

1. Fission
2. openFaas
3. kubeless
4. project riff

#### 相关链接：

11. Serverless的本质是什么？https://mp.weixin.qq.com/s/fuIhHI9VraOv64uII8NzlA
12. 迁移到Serverless之前应该想清楚的几件事 https://mp.weixin.qq.com/s/P1fKZwPaIsOQKvmQXqj2yg
1. 阿里云 https://serverless.aliyun.com/?spm=5176.137990.709885.ww4.3efa224eaiXzhJ
2. faas https://github.com/openfaas/faas
3. kubeless https://github.com/kubeless/kubeless
4. 一套平台覆盖全部主流云无服务器：Knative介绍 http://dockone.io/article/7746
5. Google发布基于Kubernetes的Serverless管理平台Knative https://www.kubernetes.org.cn/4388.html
6. Fission：基于 Kubernetes 的 Serverless 函数框架 https://www.kubernetes.org.cn/2523.html
7. 部署成功的第一个serverless API https://ckb72xuvfd.execute-api.us-east-1.amazonaws.com/staging/books?name=paas
8. 腾讯云函数式服务 scf 几个场景的Demo https://github.com/Masonlu/SCF-Demo
9. docker-lambda https://github.com/lambci/docker-lambda
10. aws-lambda文档 https://docs.aws.amazon.com/zh_cn/lambda/latest/dg/welcome.html
13. 在AWS Lambda上寫Go語言搭配API Gateway https://blog.wu-boy.com/2018/01/write-golang-in-aws-lambda/ 
14. AWS Lambda已支持用Go语言编写的无服务器应用 http://www.infoq.com/cn/news/2018/02/aws-lambda-adds-golang
15. 使用 Go 和 AWS Lambda 构建无服务 API https://juejin.im/post/5af4082f518825672a02f262
16. Serverless 实战:打造个人阅读追踪系统 https://blog.jimmylv.info/2017-06-30-serverless-in-action-build-personal-reading-statistics-system/
17. 带您玩转Lambda，轻松构建Serverless后台！https://amazonaws-china.com/cn/blogs/china/lambda-serverless/
18. knative文档 https://github.com/knative/docs
19. Serverless 应用开发指南 https://serverless.ink/

#### 使用场景：

1. 对用户不太敏感的
2. 单一的功能
3. 触发式行为的。

#### 需要考虑的要点：
1. 如何利用函数构建并部署容器
2. 如何实现扩展伸缩以响应调用函数的方式
3. 如何实现基于事件的函数调用