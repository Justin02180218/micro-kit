# micro-kit

这是关于微服务架构实战的项目，基于Go语言，结合Gin web框架和Go-kit微服务框架。内容包括限流、服务注册与发现、熔断降级、监控和分布式链路追踪、构建docker镜像，部署发布到k8s中。 

一，[项目介绍与架构](https://mp.weixin.qq.com/s?__biz=Mzg5MjA1ODYzNg==&mid=2247484113&idx=1&sn=f8ffdee22f429c18b696468d5b65ee98&chksm=cfc2afb2f8b526a426c5fe8dd125d3edcf0cefddd4b088cb300d8383cb6e28187d2036206224&token=249966631&lang=zh_CN#rd)  
二，[基础框架搭建](https://mp.weixin.qq.com/s?__biz=Mzg5MjA1ODYzNg==&mid=2247484123&idx=1&sn=156ea3f05a44ed2cb044ff767683e6a4&chksm=cfc2afb8f8b526ae4cc750c1505cafce281f5f45a426547a2b3a7bbe34577f3bfbfbc499e43d&token=404066195&lang=zh_CN#rd)  
三，[微服务library-user-service](https://mp.weixin.qq.com/s?__biz=Mzg5MjA1ODYzNg==&mid=2247484133&idx=1&sn=2635ca4e79c7a9561bd4bb1b0eb09ce3&chksm=cfc2af86f8b526903610adb87e5c48f948c83d54d77397d1c022d0776a53484363eb54f89eb7&token=404066195&lang=zh_CN#rd)  
四，[微服务library-book-service](https://mp.weixin.qq.com/s?__biz=Mzg5MjA1ODYzNg==&mid=2247484143&idx=1&sn=7d855e6039a4b4a8470736da1c2e32d9&chksm=cfc2af8cf8b5269ad85fcadae5c38d76a19dc24d6861dd5485b2f9f5b90a0c28e6013fc48411&token=404066195&lang=zh_CN#rd)  
五，[微服务library-book-grpc-service](https://mp.weixin.qq.com/s?__biz=Mzg5MjA1ODYzNg==&mid=2247484156&idx=1&sn=24ef4f146be831a291a929ca679d56d4&chksm=cfc2af9ff8b52689a0c28cc069b95819adda25407473b553b6d19788c31561ad626e5ed19019&token=404066195&lang=zh_CN#rd)  
六，[限流](https://mp.weixin.qq.com/s?__biz=Mzg5MjA1ODYzNg==&mid=2247484164&idx=1&sn=859b976d91f5e2e8578821f982af74b7&chksm=cfc2ae67f8b5277160626173db194c5d929e6faf3dc49532220079631134bcb8ecfda79ce198&token=404066195&lang=zh_CN#rd)    
七，[服务注册与发现](https://mp.weixin.qq.com/s?__biz=Mzg5MjA1ODYzNg==&mid=2247484176&idx=1&sn=8eaed74976d9ea5a0a09931e359beeaa&chksm=cfc2ae73f8b527657f95052cce2a4384874e67037539e04e9b6fd02f96316ead848bc1582787&token=404066195&lang=zh_CN#rd)  
八，[熔断、降级](https://mp.weixin.qq.com/s?__biz=Mzg5MjA1ODYzNg==&mid=2247484187&idx=1&sn=ddd393b1576df432bba95edca6f4940e&chksm=cfc2ae78f8b5276e31781197a3613d94b651c91c00c28ba20cc1a547bfaae14e672fdd739f80&token=404066195&lang=zh_CN#rd)  
九，[网关library-apigateway](https://mp.weixin.qq.com/s?__biz=Mzg5MjA1ODYzNg==&mid=2247484198&idx=1&sn=079ac356daf1370b60a74e07be56004f&chksm=cfc2ae45f8b527538f5a05bf88d0e9c1cba536e8c7cb5439787ae012b213cc45eb9f788c6ea5&token=404066195&lang=zh_CN#rd)  
十，[分布式链路追踪](https://mp.weixin.qq.com/s?__biz=Mzg5MjA1ODYzNg==&mid=2247484209&idx=1&sn=a0a597ed6c35ef7c5a7e31148f75a5d4&chksm=cfc2ae52f8b52744001755a2712acd29cee53a8c33bab3b12ad1b4b11a6d51fc3d1d6bd44213&token=404066195&lang=zh_CN#rd)  
十一，[微服务监控](https://mp.weixin.qq.com/s?__biz=Mzg5MjA1ODYzNg==&mid=2247484219&idx=1&sn=8a121493d692e6d3c3853cda7a3fa2b2&chksm=cfc2ae58f8b5274e9370ce34f309da23dec49566f4fad7c539ce6630341477827d8a433c418d&token=404066195&lang=zh_CN#rd)  
十二，[构建Docker镜像](https://mp.weixin.qq.com/s?__biz=Mzg5MjA1ODYzNg==&mid=2247484231&idx=1&sn=867a07fe549057e62d2014db8b9df384&chksm=cfc2ae24f8b527321525ff46b16966171084ac93ac2839ea08d8383982332e113a6aa5693fd2&token=404066195&lang=zh_CN#rd)  
十三，[部署到K8S](https://mp.weixin.qq.com/s?__biz=Mzg5MjA1ODYzNg==&mid=2247484249&idx=1&sn=5257c8924f136d0d3329f756c46fc893&chksm=cfc2ae3af8b5272c19a84a2e21e6fa9a8ef96c3e70a789bf6c33d01f158e9837a29d3edec984&token=404066195&lang=zh_CN#rd)  
后续未完 。。。。。。  

### 代码详解，请关注微信公众号：coding到灯火阑珊

![Image](https://github.com/Justin02180218/distribute-election-bully/blob/master/qrcode_for_gh_8a5b7b90c100_258.jpg)
