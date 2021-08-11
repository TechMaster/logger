# Thư viện log dành riêng cho iris và techmaster/eris

Thư viện này viết để log lỗi.

## Cài đặt
```
go get -u github.com/Techmaster/logger
```

Ví dụ một ứng dụng web đơn giản log lỗi khi gọi vào database. Chú ý chỉ những lỗi eris cập độ SYSERR và PANIC mới ghi vào log l
Thư mục mặc định lưu log file là `logs`

```go
package main

import (
	"github.com/Techmaster/logger"
  "github.com/Techmaster/eris"
	"github.com/kataras/iris/v12"	
)

func main() {
	app := iris.New() 

	logFile := logger.InitErrLog("logs/")  //chọn thư mục để logs ra file
	defer logFile.Close()
  app.Get("/", homepage)
	app.Listen(":8080")
}

func homepage(ctx iris.Context) {
  if posts, err := db.QueryPost(); err != nil {
    logger.Log(eris.NewFrom(err, "Failed to query post"))  //Log lối
    return
  } else {
    _, _ = ctx.JSON(posts)
  }
}
```

## Publish new module
```bash
$ git add .
$ git commit -m "v0.1.0"
$ git tag v0.1.0
$ git push origin v0.1.0
```