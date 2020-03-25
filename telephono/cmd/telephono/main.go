package main

import (
	"telephono/ui"
)
func main() {
	defer func() {
		//if recovered := recover(); recovered != nil {
		//	if panicFile, panicFileErr :=
		//		os.OpenFile("panic.txt", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0755); panicFileErr == nil {
		//		defer func() { panicFile.Close() }()
		//
		//		fmt.Println("Panicing!")
		//		if _, err := panicFile.WriteString(fmt.Sprintln("paniced! [", time.Now(), "]!!!\n")); err != nil {
		//			fmt.Println(err.Error())
		//		}
		//		panicFile.WriteString("Yup\n")
		//		panicFile.WriteString(fmt.Sprintln(recovered))
		//		panicFile.WriteString(fmt.Sprintln(string(debug.Stack())))
		//		//panicFile.WriteString(fmt.Sprintln(panicFileErr.Error()))
		//
		//		if err := panicFile.Sync(); err != nil {
		//			fmt.Println("Couldn't sync")
		//		}
		//	}
		//}
	}()
	ui.Main()

}
