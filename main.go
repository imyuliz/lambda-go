package main

import (
	"fmt"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// type book struct {
// 	ISBN   string `json:"isbn"`
// 	Title  string `json:"title"`
// 	Author string `json:"author"`
// }

// func show() (*book, error) {
// 	// 从 DynamoDB 数据库获取特定的 book 记录。在下一章中，
// 	// 我们可以让这个行为更加动态。
// 	bk, err := getItem("978-0486298238")
// 	if err != nil {
// 		return nil, err
// 	}

// 	return bk, nil
// }

// func main() {
// 	lambda.Start(show)
// }

// // 声明一个新的 DynamoDB 实例。注意它在并发调用时是
// // 安全的。
// var db = dynamodb.New(session.New(), aws.NewConfig().WithRegion("us-east-1"))

// func getItem(isbn string) (*book, error) {
// 	// 准备查询的输入
// 	input := &dynamodb.GetItemInput{
// 		TableName: aws.String("Books"),
// 		Key: map[string]*dynamodb.AttributeValue{
// 			"ISBN": {
// 				S: aws.String(isbn),
// 			},
// 		},
// 	}

// 	// 从 DynamoDB 检索数据。如果没有符合的数据
// 	// 返回 nil。
// 	result, err := db.GetItem(input)
// 	if err != nil {
// 		return nil, err
// 	}
// 	if result.Item == nil {
// 		return nil, nil
// 	}

// 	// 返回的 result.Item 对象具有隐含的
// 	// map[string]*AttributeValue 类型。我们可以使用 UnmarshalMap helper
// 	// 解析成对应的数据结构。注意：
// 	// 当你需要处理多条数据时，可以使用
// 	// UnmarshalListOfMaps。
// 	bk := new(book)
// 	err = dynamodbattribute.UnmarshalMap(result.Item, bk)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return bk, nil
// }

// // 第二版本

// var isbnRegexp = regexp.MustCompile(`[0-9]{3}\-[0-9]{10}`)
// var errorLogger = log.New(os.Stderr, "ERROR ", log.Llongfile)

// type book struct {
// 	ISBN   string `json:"isbn"`
// 	Title  string `json:"title"`
// 	Author string `json:"author"`
// }

// func show(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
// 	// 从请求中获取查询 `isbn` 的字符串参数
// 	// 并校验。
// 	isbn := req.QueryStringParameters["isbn"]
// 	if !isbnRegexp.MatchString(isbn) {
// 		return clientError(http.StatusBadRequest)
// 	}

// 	// 根据 isbn 值从数据库中取出 book 记录
// 	bk, err := getItem(isbn)
// 	if err != nil {
// 		return serverError(err)
// 	}
// 	if bk == nil {
// 		return clientError(http.StatusNotFound)
// 	}

// 	// APIGatewayProxyResponse.Body 域是个字符串，所以
// 	// 我们将 book 记录解析成 JSON。
// 	js, err := json.Marshal(bk)
// 	if err != nil {
// 		return serverError(err)
// 	}

// 	// 返回一个响应，带有代表成功的 200 状态码和 JSON 格式的 book 记录
// 	// 响应体。
// 	return events.APIGatewayProxyResponse{
// 		StatusCode: http.StatusOK,
// 		Body:       string(js),
// 	}, nil
// }

// // 添加一个用来处理错误的帮助函数。它会打印错误日志到 os.Stderr
// // 并返回一个 AWS API 网关能够理解的 500 服务器内部错误
// // 的响应。
// func serverError(err error) (events.APIGatewayProxyResponse, error) {
// 	errorLogger.Println(err.Error())

// 	return events.APIGatewayProxyResponse{
// 		StatusCode: http.StatusInternalServerError,
// 		Body:       http.StatusText(http.StatusInternalServerError),
// 	}, nil
// }

// // 加一个简单的帮助函数，用来发送和客户端错误相关的响应。
// func clientError(status int) (events.APIGatewayProxyResponse, error) {
// 	return events.APIGatewayProxyResponse{
// 		StatusCode: status,
// 		Body:       http.StatusText(status),
// 	}, nil
// }

func main() {
	lambda.Start(Hi)

}

type Do func(name string) string

func (do Do) Hello(name string) string {
	return do(name)

}

type Sayer interface {
	Hello(name string) string
}

type Object struct {
	say Sayer
}

func Hi(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string(object.say.Hello(req.QueryStringParameters["name"])),
	}, nil
}

var (
	sweetWordsMap = map[string]string{
		"roc":       "congratulate!!! %s is a xxx boy..., maybe.... and.... ",
		"yulibaozi": "%s, WINNER WINNER ,CHICKEN DINNER!",
		"yeepay":    "Meet You, 15 years old %s",
		"paas":      "welcome to yce.yeepay.com.1111 %s",
	}
	prefix = "来自serverless(lambda)的问候: "
	object = &Object{
		say: Do(
			func(name string) string {
				if data, ok := sweetWordsMap[name]; ok {
					return fmt.Sprintf(prefix+data, name)
				}
				return fmt.Sprintf(prefix+"%s 你好啊~", name)
			},
		),
	}
)
