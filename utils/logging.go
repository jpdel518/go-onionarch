package utils

import (
	"io"
	"log"
	"os"
)

func LoggingSettings(logFile string) {
	logfile, err := os.OpenFile(logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666) // なかったら作成、Read, Write, Appendも許可
	if err != nil {
		log.Fatalln(err)
	}
	multiLogFile := io.MultiWriter(os.Stdout, logfile)   // logの出力先を標準出力とログファイルに指定
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile) // format指定
	log.SetOutput(multiLogFile)
}
