package export

import (
	"encoding/csv"
	"os"
	"path"

	"github.com/oklog/ulid/v2"
)

func saveCsvFile(conf *Conf, content [][]string) (string, error) {
	filePath := path.Join(conf.FileBasePath, ulid.Make().String()+".csv")

	file, e := os.Create(filePath)
	if e != nil {
		return "", e
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	if e = writer.Write(conf.getAliases()); e != nil {
		return "", e
	}
	if e = writer.WriteAll(content); e != nil {
		return "", e
	}
	writer.Flush()
	return filePath, nil
}
