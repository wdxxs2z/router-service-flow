# router-service-flow
The route service is a demo to controller the flow where to go.

## Push 2 app, one is V1, another is V2
cf push fakea -b java_buildpack_offine -p hello_word.war -n fake
cf push fakeb -b java_buildpack_offine -p hello_work_2.war -n fakeb
## push router_service_flow app to cf  the -f is "mark":"url"
cf push mark-route-service -c "/home/vcap/app/route-service-cf -f '{\"fakea\":\"http://fake.pcf17.com\",\"fakeb\":\"http://fakeb.pcf17.com\"}' " -b binary_buildpack
## Create CUPS
cf cups mark-route-service -r https://mark-route-service.pcf17.com
## Bind the route fake
cf bind-route-service bosh-lite.com mark-route-service --hostname fake

**Now you can see:**
"fake.pcf17.com":[{"address":"172.30.51.167:60298","ttl":0,"route_service_url":"https://mark-route-service.pcf17.com"}]
</br>

## Result are：
Visit fake.pcf17.com ADD a header ：X-CF-Mark:fakea and the result is ：<h2>Hello World Thomas?</h2>
Visit fake.pcf17.com ADD a header ：X-CF-Mark:fakeb and the result is ：<h2>Hello World Thomas----- new version </h2>
