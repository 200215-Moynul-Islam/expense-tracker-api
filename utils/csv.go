package utils

import (
	"encoding/csv"
	"os"
)

func EnsureCSVFile(path string, header []string) error {

	_, err := os.Stat(path)
	if os.IsNotExist(err) {

		file, err := os.Create(path)
		if err != nil {
			return err
		}

		defer file.Close()

		writer := csv.NewWriter(file)

		err = writer.Write(header)
		if err != nil {
			return err
		}

		writer.Flush()
	}

	return nil
}
