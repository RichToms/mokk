package logging

import (
	"fmt"
	"time"
)

type Logger struct{}

func NewLogger() Logger {
	return Logger{}
}

func (l Logger) PrintLogo() {
	fmt.Println(" _______         __     __    ")
	fmt.Println("|   |   |.-----.|  |--.|  |--.")
	fmt.Println("|       ||  _  ||    < |    < ")
	fmt.Println("|__|_|__||_____||__|__||__|__|")
}

func (l Logger) TimestampedRow(str string) {
	fmt.Printf("%s | %s\n", time.Now().Format(time.TimeOnly), str)
}

func (l Logger) NewLine() {
	fmt.Println()
}
