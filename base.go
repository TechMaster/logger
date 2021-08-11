package logger

import (
	"fmt"
	"os"

	"github.com/TechMaster/eris"
	"github.com/goccy/go-json"
	"github.com/kataras/iris/v12"
)

var logFile *os.File

func InitErrLog(logFolder ...string) *os.File {
	logFile = newLogFile(logFolder...)
	return logFile
}

// Chuyên xử lý các err mà controller trả về
func Log(ctx iris.Context, err error) {
	switch e := err.(type) {
	case *eris.Error:
		if e.ErrType > eris.WARNING { //Chỉ log ra console hoặc file
			logErisError(e)
		}

		if ctx.IsAjax() { //Có trả về báo lỗi dạng JSON cho REST API request không?
			if e.Data == nil { //không có dữ liệu đi kèm thì chỉ cần in thông báo lỗi
				ctx.StatusCode(e.Code)
				_, _ = ctx.JSON(e.Error())
			} else { // Có dữ liệu bổ xung
				errorBody := map[string]interface{}{
					"error": e.Error(),
					"data":  e.Data,
				}
				ctx.StatusCode(e.Code)
				_, _ = ctx.JSON(errorBody) //Trả về cho client gọi REST API
			}
			return //Xuất ra JSON rồi thì không hiển thị Error Page nữa
		}

		if e.Data == nil {
			_ = ctx.View("error", iris.Map{
				"ErrorMsg": e.Error(),
			})
		} else {
			if bytes, err := json.Marshal(e.Data); err == nil {
				_ = ctx.View("error", iris.Map{
					"ErrorMsg": e.Error(),
					"Data":     string(bytes),
				})
			} else {
				_ = ctx.View("error", iris.Map{
					"ErrorMsg": e.Error(),
				})
			}

		}

	default: //Lỗi thông thường
		fmt.Println(err.Error()) //In ra console
		if ctx.IsAjax() {        //Trả về JSON
			ctx.StatusCode(iris.StatusInternalServerError)
			_, _ = ctx.JSON(err.Error())
		} else {
			_ = ctx.View("error", iris.Map{
				"ErrorMsg": err.Error(),
			})
		}
	}
}

//Kiểm tra xem có lỗi thì báo lỗi
func CheckErr(ctx iris.Context, err error) {
	if err != nil {
		Log(ctx, eris.WrapFrom(err, 4))
	}
}
