# kit
一个用于解析http协议的post/get请求数据的封装，其中还带有一些针对类型的判断功能，非常简单实用。

实例:
```go

import "github.com/ahwek/kit/args"

// 设置必填项 required:"请输入id"
// 设置检查手机号 phone:"true" phoneErr:"手机号不合法"
// 设置检查邮箱 email:"true" email:"邮箱格式不合法" *** 迭代中，后期实现 ***
// 设置必填项abc和def enum:"abc,def"
// 设置检查最小长度 minLen:"10" minLenErr:"姓名不得小于10个字"
// 设置检查最大长度 maxLen:"12" maxLenErr:"姓名长度不得大于12个字"
// 设置只允许的整型值 int:"1,3" intErr:"分类id只能输入1或3的数字"
type PostArg struct {
	Id      int64  `json:"id" required:"必传参数 id"`
	Cid     int32  `json:"cid" int:"1,3" intErr:"分类id只能输入1或3的数字"`
	Phone   string `json:"phone" required:"请输入手机号" phone:"true" phoneErr:"手机号不合法"`
	Title   string `json:"title" required:"请输入标题" minLen:"2" minLenErr:"长度不得小于2个字"`
	Content string `"json:content" maxLen:"10" maxLenErr:"长度不得大于10个字"`
	Tag     string `json:"tag" enum:"abc,123,456"`
}

type GetArg struct {
	Ts int64
}

// -------------------------------

// 服务器实例 http.ServeHTTP 方法
func (o *OpenApi) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    r = r.WithContext(r.Context())
	ctx := args.NewHTTPContext(w, r)

    // 设置接参
    var arg PostArg
    if err := ctx.PostArg(&arg);err != nil {
        log.Printf("err reason:%s", err.Error())
        return
    }

    fmt.Println("POST 参数数据打印", arg)

    var arg2 GetArg
    if err := ctx.GetArg(&arg2); err != nil {
        log.Printf("err reason:%s", err.Error())
        return
    }
    fmt.Println("GET 参数数据打印", arg2)
}

```

