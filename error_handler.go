package logger

import (
	"fmt"
	"os"
	"time"

	"github.com/TechMaster/eris"
	"github.com/goccy/go-json"
)

var ErisStringFormat = eris.StringFormat{
	Options: eris.FormatOptions{
		InvertOutput: false, // flag that inverts the error output (wrap errors shown first)
		WithTrace:    true,  // flag that enables stack trace output
		InvertTrace:  true,  // flag that inverts the stack trace output (top of call stack shown first)
		WithExternal: false,
		Skip:         11, // Bỏ qua 11 dòng lệnh cuối cùng trong Stack
		Top:          3,  //  Chỉ lấy 3 dòng lệnh đầu tiên
		//Mục tiêu để báo lỗi gọn hơn, stack trace đủ ngắn
	},
	MsgStackSep:  "\n",  // separator between error messages and stack frame data
	PreStackSep:  "\t",  // separator at the beginning of each stack frame
	StackElemSep: " | ", // separator between elements of each stack frame
	ErrorSep:     "\n",  // separator between each error in the chain
}

//Hàm chuyên xử lý Eris Error có Stack Trace. Chỉ áp dụng với cấp độ lỗi ERROR, SYSERROR, PANIC
func logErisError(err *eris.Error) {
	formattedStr := eris.ToCustomString(err, ErisStringFormat) //Định dạng lỗi Eris

	//Chỗ này log ra console
	if err.ErrType > eris.ERROR { //Với lỗi cao hơn ERROR gồm SYSERROR và PANIC, in ra mầu đỏ và ghi ra file
		colorReset := string("\033[0m")
		colorMagenta := string("\033[35m")
		fmt.Println(colorMagenta, formattedStr, colorReset)
		var textToFile string

		if dataString := marshalErisData2JSON(err); dataString != "" { //Nếu có dữ liệu đi kèm thì cũng ghi ra file
			fmt.Println(colorMagenta, dataString, colorReset)
			textToFile = time.Now().Format("2006 01 02-15:04:05 - ") + formattedStr + "\n" + dataString + "\n\n"
		} else {
			textToFile = time.Now().Format("2006 01 02-15:04:05 - ") + formattedStr + "\n\n"
		}
		//Lỗi Panic và Error nhất thiết phải ghi vào file !
		if _, err := logFile.WriteString(textToFile); err != nil {
			panic(err)
		}
	} else {
		fmt.Println(formattedStr) //Error Level
	}
}

func marshalErisData2JSON(err *eris.Error) string {
	if err.Data != nil {
		if dataStr, err := json.Marshal(err.Data); err == nil {
			return string(dataStr)
		}
	}
	return ""
}

func todayFilename() string {
	today := time.Now().Format("2006 01 02")
	return today + ".txt"
}
func newLogFile(logFolder ...string) *os.File {
	var log_folder string
	if len(logFolder) == 0 {
		log_folder = "logs/"
	} else {
		log_folder = logFolder[0]
	}
	filename := todayFilename()
	// Open the file, this will append to the today's file if server restarted.
	f, err := os.OpenFile(log_folder+filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}

	return f
}
