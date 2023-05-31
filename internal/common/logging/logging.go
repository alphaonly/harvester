package logging

import "log"

func LogFatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
func LogPrintln(err error) {
	if err != nil {
		log.Println(err)
	}
}