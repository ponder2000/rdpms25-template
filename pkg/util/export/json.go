package export

import (
	"os"
	"path"

	"github.com/oklog/ulid/v2"
)

func saveJsonFile(conf *Conf, bytes []byte) (string, error) {
	filePath := path.Join(conf.FileBasePath, ulid.Make().String()+".json")

	file, e := os.Create(filePath)
	if e != nil {
		return "", e
	}
	defer file.Close()

	if _, e = file.Write(bytes); e != nil {
		return "", e
	}
	return filePath, nil
}
