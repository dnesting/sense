package senseutil

import (
	"log"
	"net/http"
	"net/http/httputil"
	"strings"
)

func PrefixLines(prefix string, lines string) string {
	return prefix + strings.ReplaceAll(strings.TrimSpace(lines), "\n", "\n"+prefix)
}

func DumpRequest(log *log.Logger, r *http.Request) {
	if log == nil {
		return
	}
	bytes, err := httputil.DumpRequestOut(r, true)
	if err != nil {
		log.Println("error dumping request:", err)
	} else {
		log.Println("HTTP request:\n" + PrefixLines("> ", string(bytes)))
		log.Println()
	}
}

func DumpResponse(log *log.Logger, r *http.Response) {
	bytes, err := httputil.DumpResponse(r, true)
	if err != nil {
		log.Println("error dumping response:", err)
	} else {
		log.Println("HTTP response:\n" + PrefixLines("> ", string(bytes)))
		log.Println()
	}
}
