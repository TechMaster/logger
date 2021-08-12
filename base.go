package logger

import (
	"fmt"
	"os"

	"github.com/TechMaster/eris"
	"github.com/goccy/go-json"
	"github.com/kataras/iris/v12"
)

type LogConfig struct {
	Log_folder     string // thư mục chứa log file. Nếu rỗng có nghĩa là không ghi log ra file
	Error_template string // tên view template sẽ render error page
	Skip           int    // số dòng cuối cùng trong stack trace sẽ bị bỏ qua
	Top            int    // số dòng đỉnh stack trace sẽ được in ra
}

var logConfig LogConfig
var logFile *os.File

func Init(_logConfig LogConfig) *os.File {
	logConfig = _logConfig
	if logConfig.Log_folder != "" {
		logFile = newLogFile(logConfig.Log_folder)
		return logFile
	} else {
		return nil
	}
}

// Chuyên xử lý các err mà controller trả về
func Log(ctx iris.Context, err error) {
	switch e := err.(type) {
	case *eris.Error:
		if e.ErrType > eris.WARNING { //Chỉ log ra console hoặc file
			logErisError(e)
		}

		if ctx.IsAjax() { //Có trả về báo lỗi dạng JSON cho REST API request không?
			errorBody := iris.Map{
				"error": e.Error(),
			}
			if e.Data != nil { //không có dữ liệu đi kèm thì chỉ cần in thông báo lỗi
				errorBody["data"] = e.Data
			}
			ctx.StatusCode(e.Code)
			_, _ = ctx.JSON(errorBody) //Trả về cho client gọi REST API
			return                     //Xuất ra JSON rồi thì không hiển thị Error Page nữa
		}

		// Nếu request không phải là REST request (AJAX request) thì render error page
		ctx.ViewData("ErrorMsg", e.Error())
		if e.Data != nil {
			if bytes, err := json.Marshal(e.Data); err == nil {
				ctx.ViewData("Data", string(bytes))
			}
		}
		_ = ctx.View(logConfig.Error_template)
		return

	default: //Lỗi thông thường
		fmt.Println(err.Error()) //In ra console
		if ctx.IsAjax() {        //Trả về JSON
			ctx.StatusCode(iris.StatusInternalServerError)
			_, _ = ctx.JSON(err.Error())
		} else {
			_ = ctx.View(logConfig.Error_template, iris.Map{
				"ErrorMsg": err.Error(),
			})
		}
		return
	}
}

//Kiểm tra xem có lỗi thì báo lỗi
func CheckErr(ctx iris.Context, err error) {
	if err != nil {
		Log(ctx, eris.WrapFrom(err, 4))
	}
}
