package utils

import "log"

func Showerr(msg string, rollback bool) {
	tag := "Error : "

	if rollback {
		tag += "Rollback : "
	}

	log.Println(tag + msg)
}
