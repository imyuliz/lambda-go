
#### 相关文档
1. 在AWS Lambda上寫Go語言搭配API Gateway https://blog.wu-boy.com/2018/01/write-golang-in-aws-lambda/ 
2. AWS Lambda已支持用Go语言编写的无服务器应用 http://www.infoq.com/cn/news/2018/02/aws-lambda-adds-golang
3. 使用 Go 和 AWS Lambda 构建无服务 API https://juejin.im/post/5af4082f518825672a02f262

创建角色
```
aws iam create-role --role-name booksuser --assume-role-policy-document file:///Users/yulibaozi/GoWorkSpace/src/github.com/yulibaozi/lambda-go/trust-policy.json
```

```
{
    "Role": {
        "AssumeRolePolicyDocument": {
            "Version": "2012-10-17",
            "Statement": [
                {
                    "Action": "sts:AssumeRole",
                    "Effect": "Allow",
                    "Principal": {
                        "Service": "lambda.amazonaws.com"
                    }
                }
            ]
        },
        "RoleId": "AROAJJCXREEZRABUJJLN4",
        "CreateDate": "2018-08-21T02:55:56Z",
        "RoleName": "booksuser",
        "Path": "/",
        "Arn": "arn:aws:iam::918101530468:role/booksuser"
  
```

绑定role
```
aws iam attach-role-policy --role-name booksuser --policy-arn arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole
```

上传可执行文件
```
aws lambda create-function --function-name booksone --runtime go1.x  --role arn:aws:iam::918101530468:role/booksuser --handler lambda-go  --zip-file fileb:///Users/yulibaozi/GoWorkSpace/src/github.com/yulibaozi/lambda-go/lambda-go.zip
```

respone:
```
{
    "TracingConfig": {
        "Mode": "PassThrough"
    },
    "CodeSha256": "TUOCU5vxWo5cEOEDSLAg6UFr/gUen8LH0/uacZmn92o=",
    "FunctionName": "booksone",
    "CodeSize": 3016412,
    "RevisionId": "361350da-f2a3-49a3-8f5d-5354f8e2bb8d",
    "MemorySize": 128,
    "FunctionArn": "arn:aws:lambda:us-east-1:918101530468:function:booksone",
    "Version": "$LATEST",
    "Role": "arn:aws:iam::918101530468:role/booksuser",
    "Timeout": 3,
    "LastModified": "2018-08-21T03:07:15.988+0000",
    "Handler": "lambda-go",
    "Runtime": "go1.x",
    "Description": ""
}
```

```
--function-name
将在 AWS 中被调用的 lambda 函数名


--runtime
lambda 函数的运行环境（在我们的例子里用 "go1.x"）


--role
你想要 lambda 函数在运行时扮演的角色的 ARN（见上面的步骤 6）


--handler
zip 文件根目录下的可执行文件的名称


--zip-file
zip 文件的路径
```

测试是否可用
```
aws lambda invoke --function-name booksone output.json
```
output是输出的结果保存的地方

```
{
    "ExecutedVersion": "$LATEST",
    "StatusCode": 200
}
```

第二部分：(DynamoDB)
aws dynamodb create-table --table-name Books \
--attribute-definitions AttributeName=ISBN,AttributeType=S \
--key-schema AttributeName=ISBN,KeyType=HASH \
--provisioned-throughput ReadCapacityUnits=5,WriteCapacityUnits=5

```
{
    "TableDescription": {
        "TableArn": "arn:aws:dynamodb:us-east-1:918101530468:table/Books",
        "AttributeDefinitions": [
            {
                "AttributeName": "ISBN",
                "AttributeType": "S"
            }
        ],
        "ProvisionedThroughput": {
            "NumberOfDecreasesToday": 0,
            "WriteCapacityUnits": 5,
            "ReadCapacityUnits": 5
        },
        "TableSizeBytes": 0,
        "TableName": "Books",
        "TableStatus": "CREATING",
        "TableId": "68bad907-e529-419c-af76-a00c0dcac46a",
        "KeySchema": [
            {
                "KeyType": "HASH",
                "AttributeName": "ISBN"
            }
        ],
        "ItemCount": 0,
        "CreationDateTime": 1534821144.495
    }
}
```
400错误：
```
{"errorMessage":"AccessDeniedException: User: arn:aws:sts::918101530468:assumed-role/booksuser/booksone is not authorized to perform: dynamodb:GetItem on resource: arn:aws:dynamodb:us-east-1:918101530468:table/Books\n\tstatus code: 400, request id: 2VIOG5K254JMUNPI64V3F5FFG3VV4KQNSO5AEMVJF66Q9ASUAAJG","errorType":"requestError"}
```
解决办法：修改policyjson文件的策略

```
aws iam put-role-policy --role-name booksuser  --policy-name dynamodb-item-crud-role --policy-document file:///Users/yulibaozi/GoWorkSpace/src/github.com/yulibaozi/lambda-go/trust-policy.json
```

第三部分：apigateway

1. 创建API网关（获取rest-api-id）
```
➜  lambda-go aws apigateway create-rest-api --name bookstore
{
    "apiKeySource": "HEADER",
    "name": "bookstore",
    "createdDate": 1534822199,
    "endpointConfiguration": {
        "types": [
            "EDGE"
        ]
    },
    "id": "ckb72xuvfd"
}
```

2. 不知道这步做什么的(获取根目录，主要是rootpath-id)
```
➜  lambda-go aws apigateway get-resources --rest-api-id ckb72xuvfd
{
    "items": [
        {
            "path": "/",
            "id": "7lgb2dvfjf"
        }
    ]
}
```
3. 创建访问目录（/books）获取resourceId
```
➜  lambda-go aws apigateway create-resource --rest-api-id ckb72xuvfd --parent-id 7lgb2dvfjf --path-part books
{
    "path": "/books",
    "pathPart": "books",
    "id": "00h2c6",
    "parentId": "7lgb2dvfjf"
}
```
4. 修改http的请求方式
```
➜  lambda-go aws apigateway put-method --rest-api-id ckb72xuvfd --resource-id 00h2c6 --http-method ANY --authorization-type NONE
{
    "apiKeyRequired": false,
    "httpMethod": "ANY",
    "authorizationType": "NONE"
}
```
5. 把上面创建的HTTP 网关整合进lambda
```
➜  lambda-go aws apigateway put-integration --rest-api-id ckb72xuvfd --resource-id 00h2c6 --http-method ANY --type AWS_PROXY --integration-http-method POST \
--uri arn:aws:apigateway:us-east-1:lambda:path/2015-03-31/functions/arn:aws:lambda:us-east-1:918101530468:function:booksone/invocations

值得注意的是：918101530468:function:booksone 是:serviceid:function:function-name
{
    "passthroughBehavior": "WHEN_NO_MATCH",
    "timeoutInMillis": 29000,
    "uri": "arn:aws:apigateway:us-east-1:lambda:path/2015-03-31/functions/arn:aws:lambda:us-east-1:918101530468:function:booksone/invocations",
    "httpMethod": "POST",
    "cacheNamespace": "00h2c6",
    "type": "AWS_PROXY",
    "cacheKeyParameters": []
}
```

6. 测试API

```
➜  lambda-go aws apigateway test-invoke-method --rest-api-id ckb72xuvfd --resource-id 00h2c6 --http-method "GET"
{
    "status": 500,
    "body": "{\"message\": \"Internal server error\"}",
    "log": "Execution log for request 1cda20c4-a4f5-11e8-a408-8d0b78241a21\nTue Aug 21 03:48:49 UTC 2018 : Starting execution for request: 1cda20c4-a4f5-11e8-a408-8d0b78241a21\nTue Aug 21 03:48:49 UTC 2018 : HTTP Method: GET, Resource Path: /books\nTue Aug 21 03:48:49 UTC 2018 : Method request path: {}\nTue Aug 21 03:48:49 UTC 2018 : Method request query string: {}\nTue Aug 21 03:48:49 UTC 2018 : Method request headers: {}\nTue Aug 21 03:48:49 UTC 2018 : Method request body before transformations: \nTue Aug 21 03:48:49 UTC 2018 : Endpoint request URI: https://lambda.us-east-1.amazonaws.com/2015-03-31/functions/arn:aws:lambda:us-east-1:918101530468:function:booksone/invocations\nTue Aug 21 03:48:49 UTC 2018 : Endpoint request headers: {x-amzn-lambda-integration-tag=1cda20c4-a4f5-11e8-a408-8d0b78241a21, Authorization=************************************************************************************************************************************************************************************************************************************************************************************************************************0629eb, X-Amz-Date=20180821T034849Z, x-amzn-apigateway-api-id=ckb72xuvfd, X-Amz-Source-Arn=arn:aws:execute-api:us-east-1:918101530468:ckb72xuvfd/test-invoke-stage/GET/books, Accept=application/json, User-Agent=AmazonAPIGateway_ckb72xuvfd, X-Amz-Security-Token=FQoGZXIvYXdzEPT//////////wEaDE5kQLYCjn/fIbbUwiK3A1MMtK+JHJkaYiDbGHaBChG3Ym/EzLzKJr0HEc/akNODnkJ5vUWYYbbDkwT+0ktp/XYrXZ29jCkhf713/V7PsTxWLAtlfWw+T62AsXg4sQLydIh+cBkAtVhZAALa0guEeHQjB1e8V3eLb+6JIrg8yQq+iaj60AJEfODS55C8Lg4XWb/OHNp7qzxQq7Yi1bN4nTNEJlQBjsx0hQAsZxIdoupq8uq/Opiz41zjkrLaRWg37jlmVHZuzTeIu8x+o1eRuR9KRwTY1ccNcp19T/aev/9MTOU732H2AOwQLAJ [TRUNCATED]\nTue Aug 21 03:48:49 UTC 2018 : Endpoint request body after transformations: {\"resource\":\"/books\",\"path\":\"/books\",\"httpMethod\":\"GET\",\"headers\":null,\"queryStringParameters\":null,\"pathParameters\":null,\"stageVariables\":null,\"requestContext\":{\"path\":\"/books\",\"accountId\":\"918101530468\",\"resourceId\":\"00h2c6\",\"stage\":\"test-invoke-stage\",\"requestId\":\"1cda20c4-a4f5-11e8-a408-8d0b78241a21\",\"identity\":{\"cognitoIdentityPoolId\":null,\"cognitoIdentityId\":null,\"apiKey\":\"test-invoke-api-key\",\"cognitoAuthenticationType\":null,\"userArn\":\"arn:aws:iam::918101530468:root\",\"apiKeyId\":\"test-invoke-api-key-id\",\"userAgent\":\"aws-cli/1.15.81 Python/2.7.13 Darwin/17.7.0 botocore/1.10.80\",\"accountId\":\"918101530468\",\"caller\":\"918101530468\",\"sourceIp\":\"test-invoke-source-ip\",\"accessKey\":\"AKIAJPUW7XFZACKNO4IA\",\"cognitoAuthenticationProvider\":null,\"user\":\"918101530468\"},\"resourcePath\":\"/books\",\"httpMethod\":\"GET\",\"extendedRequestId\":\"L9LBMFuUIAMFwjA=\",\"apiId\":\"ckb72xuvfd\"},\"body\":null,\"isBase64Encoded\":false}\nTue Aug 21 03:48:49 UTC 2018 : Sending request to https://lambda.us-east-1.amazonaws.com/2015-03-31/functions/arn:aws:lambda:us-east-1:918101530468:function:booksone/invocations\nTue Aug 21 03:48:49 UTC 2018 : Execution failed due to configuration error: Invalid permissions on Lambda function\nTue Aug 21 03:48:49 UTC 2018 : Method completed with status: 500\n",
    "latency": 15,
    "headers": {}
}
```
7. 修复权限导致的不可访问
```
Execution failed due to configuration error: Invalid permissions on Lambda function
```
```
aws lambda add-permission --function-name booksone --statement-id a-GUID --action lambda:InvokeFunction --principal apigateway.amazonaws.com
```
8.   Malformed Lambda proxy response\nTue Aug 21 04:49:54 UTC 2018 : Method completed with status: 50
```
这步是修改代码为events，应该格式不能解析
```
9. 命令行访问
```
aws apigateway test-invoke-method --rest-api-id ckb72xuvfd \
--resource-id 00h2c6 --http-method "GET" \
--path-with-query-string "/books?isbn=978-1420931693"
```
respone:
```
{
    "status": 200,
    "body": "{\"isbn\":\"978-1420931693\",\"title\":\"The Republic\",\"author\":\"Plato\"}",
    "log": "Execution log for request 2be4ca29-a4ff-11e8-9a45-afa7c876e658\nTue Aug 21 05:00:49 UTC 2018 : Starting execution for request: 2be4ca29-a4ff-11e8-9a45-afa7c876e658\nTue Aug 21 05:00:49 UTC 2018 : HTTP Method: GET, Resource Path: /book
s\nTue Aug 21 05:00:49 UTC 2018 : Method request path: {}\nTue Aug 21 05:00:49 UTC 2018 : Method request query string: {isbn=978-1420931693}\nTue Aug 21 05:00:49 UTC 2018 : Method request headers: {}\nTue Aug 21 05:00:49 UTC 2018 : Method request body before transformations: \nTue Aug 21 05:00:49 UTC 2018 : Endpoint request URI: https://lambda.us-east-1.amazonaws.com/2015-03-31/functions/arn:aws:lambda:us-east-1:918101530468:function:booksone/invocations\nTue Aug 21 05:00:49 UTC 2018 : Endpoint request headers: {x-amzn-lambda-integration-tag=2be4ca29-a4ff-11e8-9a45-afa7c876e658, Authorization=************************************************************************************************************************************************************************************************************************************************************************************************************************54975b, X-Amz-Date=20180821T050049Z, x-amzn-apigateway-api-id=ckb72xuvfd, X-Amz-Source-Arn=arn:aws:execute-api:us-east-1:918101530468:ckb72xuvfd/test-invoke-stage/GET/books, Accept=application/json, User-Agent=AmazonAPIGateway_ckb72xuvfd, X-Amz-Security-Token=FQoGZXIvYXdzEPX//////////wEaDFO7bZQWkbwKO+DjwiK3AwOpGZgwBHzpDLkPPgEA/XlWozpdVIGPm7EiHqZKqWXEwJ5eALYaHE2VnuNaCIoUJBm+40Csy4kbebDPJdLc9ER/Dj3WHJ7n5bbNwjpjDUABeBhCY1pzNYKis9nbEPx7nGf8+W0fq6uUVOMAvC3YFGNXAQVo8i7AX4/DjgP78iS3t3wAKLkOGyDH5Py9zwMBXNtgKvF3+gvDVlVVooJs/JUeKvuIGpzqf7HR4BSvrBvchgrRGuCvwJ0f4RAb/twbTakhkAdc3ekCDZ6onic8XA4PLB3ggWkf7kodOli [TRUNCATED]\nTue Aug 21 05:00:49 UTC 2018 : Endpoint request body after transformations: {\"resource\":\"/books\",\"path\":\"/books\",\"httpMethod\":\"GET\",\"headers\":null,\"queryStringParameters\":{\"isbn\":\"978-1420931693\"},\"pathParameters\":null,\"stageVariables\":null,\"requestContext\":{\"path\":\"/books\",\"accountId\":\"918101530468\",\"resourceId\":\"00h2c6\",\"stage\":\"test-invoke-stage\",\"requestId\":\"2be4ca29-a4ff-11e8-9a45-afa7c876e658\",\"identity\":{\"cognitoIdentityPoolId\":null,\"cognitoIdentityId\":null,\"apiKey\":\"test-invoke-api-key\",\"cognitoAuthenticationType\":null,\"userArn\":\"arn:aws:iam::918101530468:root\",\"apiKeyId\":\"test-invoke-api-key-id\",\"userAgent\":\"aws-cli/1.15.81 Python/2.7.13 Darwin/17.7.0 botocore/1.10.80\",\"accountId\":\"918101530468\",\"caller\":\"918101530468\",\"sourceIp\":\"test-invoke-source-ip\",\"accessKey\":\"AKIAJPUW7XFZACKNO4IA\",\"cognitoAuthenticationProvider\":null,\"user\":\"918101530468\"},\"resourcePath\":\"/books\",\"httpMethod\":\"GET\",\"extendedRequestId\":\"L9VkOGvIIAMFm1A=\",\"apiId\":\"ckb72xuvfd\"},\"body\":null,\"isBase64Encoded\":false}\nTue Aug 21 05:00:49 UTC 2018 : Sending request to https://lambda.us-east-1.amazonaws.com/2015-03-31/functions/arn:aws:lambda:us-east-1:918101530468:function:booksone/invocations\nTue Aug 21 05:00:50 UTC 2018 : Received response. Integration latency: 1082 ms\nTue Aug 21 05:00:50 UTC 2018 : Endpoint response body before transformations: {\"statusCode\":200,\"headers\":null,\"body\":\"{\\\"isbn\\\":\\\"978-1420931693\\\",\\\"title\\\":\\\"The Republic\\\",\\\"author\\\":\\\"Plato\\\"}\"}\nTue Aug 21 05:00:50 UTC 2018 : Endpoint response headers: {X-Amz-Executed-Version=$LATEST, x-amzn-Remapped-Content-Length=0, Connection=keep-alive, x-amzn-RequestId=2be7899a-a4ff-11e8-84f3-addfd7877e53, Content-Length=120, Date=Tue, 21 Aug 2018 05:00:50 GMT, X-Amzn-Trace-Id=root=1-5b7b9c81-2dc977373387bc48d6098f3b;sampled=0, Content-Type=application/json}\nTue Aug 21 05:00:50 UTC 2018 : Method response body after transformations: {\"isbn\":\"978-1420931693\",\"title\":\"The Republic\",\"author\":\"Plato\"}\nTue Aug 21 05:00:50 UTC 2018 : Method response headers: {X-Amzn-Trace-Id=Root=1-5b7b9c81-2dc977373387bc48d6098f3b;Sampled=0}\nTue Aug 21 05:00:50 UTC 2018 : Successfully completed execution\nTue Aug 21 05:00:50 UTC 2018 : Method completed with status: 200\n",
    "latency": 1099,
    "headers": {
        "X-Amzn-Trace-Id": "Root=1-5b7b9c81-2dc977373387bc48d6098f3b;Sampled=0"
    }
}
```

9. 开始部署 deployment(api命名为staging)
```
aws apigateway create-deployment --rest-api-id ckb72xuvfd \
--stage-name staging

respone:

{
    "id": "e2mudv",
    "createdDate": 1534827820
}
```

10. API访问：
```
https://rest-api-id.execute-api.us-east-1.amazonaws.com/staging

```