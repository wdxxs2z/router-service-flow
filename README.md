# router-service-flow
The route service is a demo to controller the flow where to go.

1.上传我们的flow-route-service应用 比例为2:3
#####1).支持算法：加权取模
```
cf push flow-route-service -c "/home/vcap/app/route-service-cf-flow -r '{\"typename\": \"modulo\",\"nodes\":[{\"index\":1,\"url\":\"http://www.demo1.com\",\"weight\":2},{\"index\":1,\"url\":\"http://www.demo2.com\",\"weight\":3}]}' " -b binary_buildpack

```
#####2).支持算法：轮询
```
cf push flow-route-service -c "/home/vcap/app/route-service-cf-flow -r '{\"typeName\": \"roundrobin\",\"nodes\":[{\"index\":1,\"url\":\"http://www.demo1.com\",\"weight\":2},{\"index\":1,\"url\":\"http://www.demo2.com\",\"weight\":3}]}' " -b binary_buildpack
```
#####3).支持算法：加权轮询
```
cf push flow-route-service -c "/home/vcap/app/route-service-cf-flow -r '{\"typeName\": \"roundrobinweight\",\"nodes\":[{\"index\":1,\"url\":\"http://www.demo1.com\",\"weight\":2},{\"index\":1,\"url\":\"http://www.demo2.com\",\"weight\":3}]}' " -b binary_buildpack
```
#####4).支持算法：源哈希
```
cf push flow-route-service -c "/home/vcap/app/route-service-cf-flow -r '{\"typeName\": \"sourcehash\",\"nodes\":[{\"index\":1,\"url\":\"http://www.demo1.com\",\"weight\":2},{\"index\":1,\"url\":\"http://www.demo2.com\",\"weight\":3}]}' " -b binary_buildpack
```
2.将我们的flow-route-service设置成一个路由服务
```
cf cups mark-route-service -r https://flow-route-service.bosh-lite.com
```
3.将这个服务绑定到受影响的路由上
```
cf bind-route-service bosh-lite.com flow-route-service --hostname fake
```
</br>
4.测试结果：
```
ubuntu@pivotal-ops-manager:~$ for i in `seq 1 1000`;do curl -q http://fake.pcf17.com/ 2>/dev/null |grep "Hello";done|sort |uniq -c
    423 >Hello World Thomas?
    577 Hello World Thomas----- new version
```
```
ubuntu@pivotal-ops-manager:~$ for i in `seq 1 10000`;do curl -q http://fake.pcf17.com/ 2>/dev/null |grep "Hello";done|sort |uniq -c
   3951 >Hello World Thomas?</h2>
   6049 Hello World Thomas----- new version
```   
可以看出比例是2/3
