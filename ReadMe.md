# Thư viện log dành riêng cho iris và techmaster/eris

Thư viện này viết để log lỗi và thông báo. 

Nếu client gửi AJAX request hoặc request.ContentType='application/json' thì trả về JSON error. Nếu client gọi GET/POST từng web page thì trả về web page.

Ứng dụng có hai phương thức chính:
```go
logger.Log(ctx, eris.NewFrom(err, "Failed to query post"))  //Log lỗi
logger.Info(ctx,"Kích hoạt thành công")  //Thông báo
```

## Cài đặt
```
go get -u github.com/TechMaster/logger
```

Ví dụ một ứng dụng web đơn giản log lỗi khi gọi vào database. Chú ý chỉ những lỗi eris cập độ SYSERR và PANIC mới ghi vào log l
Thư mục mặc định lưu log file là `logs`

Cấu hình mặc định
```go
package main

import (
	"github.com/Techmaster/logger"
  "github.com/Techmaster/eris"
	"github.com/kataras/iris/v12"	
)

func main() {
	app := iris.New() 

	logFile := logger.Init()
  if logFile != nil {
		defer logFile.Close()
	}

  app.Get("/", homepage)
	app.Listen(":8080")
}
```


Cấu hình chi tiết hơn
```go
package main

import (
	"github.com/Techmaster/logger"
  "github.com/Techmaster/eris"
	"github.com/kataras/iris/v12"	
)

func main() {
	app := iris.New() 

	logFile := logger.Init(logger.LogConfig{
		LogFolder:     "logs/", //Nếu để rỗng thì không ghi log
		ErrorTemplate: "error", //Cần phải có file error.html ở thư mục views để render error page
    InfoTemplate:  "info", //Cần phải có file error.html ở thư mục views để render error page
		Top:            3,       //Lấy 3 hàm đầu tiên trên đỉnh stack trace
		Skip:           11,      //hoặc loại đi 11 hàm đáy của stack trace
	})
	if logFile != nil {
		defer logFile.Close()
	}

  app.Get("/", homepage)
	app.Listen(":8080")
}

func homepage(ctx iris.Context) {
  if posts, err := db.QueryPost(); err != nil {
    logger.Log(ctx, eris.NewFrom(err, "Failed to query post"))  //Log lối
    return
  } else {
    _, _ = ctx.JSON(posts)
  }
}
```


## Ví dụ bổ xung
```go
func Handler(ctx iris.Context) {
	if err := PhuTro("val1", 10); err != nil {
		logger.Log(err)  //Log error ở đây sau đó return luôn
		return
	}
	
}

func PhuTro(para1 string, para2 int) error {
	if err := Db.Connnect(connectionstr); err != nil {
    //Luôn bọc lỗi thông thường bằng eris để có stack trace
		return eris.NewFromMsg(err, "Lỗi kết nối CSDL").BadRequest.SetType(eris.SysError) 
	}
}
```

### Thông báo chứ không phải báo lỗi
Tham số truyền như sau:
  1. Thông điệp
  2. Đường link "/dashboard"
  3. Tên đường link "Dashboard"

```go
logger.Info(ctx, "Tài khoản đã được kích hoạt, ấn link để đi tiếp", "/dashboard", "Dashboard")
```

## Trả về JSON Error hay HTML Error Page tuỳ thuộc vào

Request gọi lên là AJAX Request hoặc có Content Type là "application/json"

```go
func Log(ctx iris.Context, err error) {
	//Trả về JSON error khi client gọi lên bằng AJAX hoặc request.ContentType dạng application/json
	shouldReturnJSON := ctx.IsAjax() || ctx.GetContentTypeRequested() == "application/json"
  ...
}
```

## Xử lý JSON Error trả về
Cấu trúc JSON trả về gồm 2 trường:
1. error dạng string
2. data dạng struct bất kỳ

Hãy truy cập `err.responseJSON` để lấy dữ liệu lỗi trả về

```javascript
function sendEmail(type) {
  $.post("/email/send?type=" + type, 
  { 
    to: $("#to").val(),
    subject: $("#subject").val(),
    body: $("#body").val(),
  })
    .done(data => {  //status code 200
      $("#result").html(data).css('color', 'black');
    })
    .fail(data => {  //400, 401, 404, 500
      console.log(err);
      $("#result").html(err.responseJSON.error).css('color', 'red');
    })
}
```


## Publish new module
```bash
$ git add .
$ git commit -m "v0.1.0"
$ git tag v0.1.0
$ git push origin v0.1.0
```