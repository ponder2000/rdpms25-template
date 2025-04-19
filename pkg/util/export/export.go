package export

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"

	"github.com/Jeffail/gabs/v2"
)

type Conf struct {
	Format       string
	FileBasePath string
	Cols         []*ColConf

	Pdf *PdfConf
}

func (r *Conf) getAliases() []string {
	res := make([]string, 0, len(r.Cols))
	for _, col := range r.Cols {
		res = append(res, col.Alias)
	}
	return res
}

type PdfConf struct {
	Title       string
	Orientation string // L or P
	PageSize    string // A4, A3 etc
}

type ColConf struct {
	Key         string
	Alias       string
	Transformer func(any) string
}

func Table[T any](conf *Conf, data []T) (string, error) {
	if conf.Format == "json" {
		rawData, _ := json.Marshal(data)
		return saveJsonFile(conf, rawData)
	}

	tableContent := make([][]string, 0)
	for _, row := range data {
		rawData, e := json.Marshal(row)
		if e != nil {
			slog.Warn("Unable to marshal object", "row", row)
			continue
		}

		jsonParsed, e := gabs.ParseJSON(rawData)
		if e != nil {
			slog.Warn("Unable to gab parse", "err", e.Error())
			continue
		}

		values := make([]string, 0)
		for _, cc := range conf.Cols {
			values = append(values, getValue(cc, jsonParsed))
		}
		tableContent = append(tableContent, values)
	}

	switch conf.Format {
	case "pdf":
		if conf.Pdf == nil {
			return "", errors.New("pdf configuration is missing")
		}
		return savePdfFile(conf, tableContent)
	case "xlsx":
		return saveExcelFile(conf, tableContent)
	case "csv":
		return saveCsvFile(conf, tableContent)
	default:
		return "", fmt.Errorf("unsupported format: %s", conf.Format)
	}
}

func getValue(conf *ColConf, jsonParsed *gabs.Container) string {
	var data any
	if conf.Key == "" {
		data = jsonParsed.Data()
	} else {
		data = jsonParsed.Path(conf.Key).Data()
	}

	if conf.Transformer != nil {
		return conf.Transformer(data)
	} else if data == nil {
		return ""
	} else {
		return fmt.Sprintf("%v", data)
	}
}
