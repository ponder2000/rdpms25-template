package formatter

import (
	"log/slog"
	"strconv"
	"strings"
	"time"

	"github.com/ponder2000/rdpms25-template/pkg/global"
)

func GetAvailableFormats() []string {
	availableFormats := []string{
		"UTC : converts 2023-04-12T07:24:01.105409Z to date time",
		"epoch : converts epoch since milliseconds to date time",
		"lower : converts to lower case",
		"upper : converts to upper case",
		"/1000 : divides valid int by 1000",
		"*1000 : multiplies valid int by 1000",
		"status : true false to ACTIVE CLEARED",
	}
	return availableFormats
}

func Format(f string, val string) any {
	switch f {
	case "UTC":
		t, e := time.Parse(time.RFC3339, val)
		if e != nil {
			slog.Error("unable to format", "err", e)
			return "Wrong type!"
		}
		return t.UTC().Local().Format(global.DateTimeFormat)
	case "epoc":
		t, e := strconv.ParseInt(val, 10, 64)
		if e != nil {
			return ""
		}
		return time.Unix(t/1000, 0).UTC().Local()
	case "lower":
		return strings.ToLower(val)
	case "upper":
		return strings.ToUpper(val)
	case "/1000":
		t, e := strconv.ParseInt(val, 10, 64)
		if e != nil {
			return ""
		}
		return t / 1000
	case "*1000":
		t, e := strconv.ParseInt(val, 10, 64)
		if e != nil {
			return ""
		}
		return t * 1000
	case "status":
		if val == "true" {
			return "ACTIVE"
		} else if val == "false" {
			return "CLEARED"
		} else {
			return "-"
		}
	}
	return val
}
