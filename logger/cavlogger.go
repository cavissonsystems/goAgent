package cavlogger

import (
	"io"
	"log"
	"os"
)

var (
	Trace *log.Logger
	Error *log.Logger
)

func Init(traceHandle io.Writer, errorHandle io.Writer) {
	Trace = log.New(traceHandle, "TRACE: ", log.Ldate|log.Ltime|log.Lshortfile)
	Error = log.New(errorHandle, "Error: ", log.Ldate|log.Ltime|log.Lshortfile)
}

func OpenFiles() {
	ndhome := os.Getenv("NDHOME")
	trace_file, trace_err := os.OpenFile(ndhome+"log/trace_logs.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if trace_err != nil {
		log.Fatalln("Failed to open trace log file", trace_err)
	}

	error_file, error_err := os.OpenFile(ndhome+"/log/error_logs.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if error_err != nil {
		log.Fatalln("Failed to open error log file", error_err)
	}

	Init(trace_file, error_file)
}

func TracePrint(s string) {
	Trace.Println(s)
}

func ErrorPrint(s string) {
	Error.Println(s)
}
