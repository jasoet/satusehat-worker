package file

import (
	"encoding/json"
	"os"
)

func WritePrettyJson(fileName string, payload []byte, perm os.FileMode) error {
	var jsonData any
	err := json.Unmarshal(payload, &jsonData)
	if err != nil {
		return err
	}

	prettyPayload, err := json.MarshalIndent(jsonData, "", "    ")
	if err != nil {
		return err
	}

	return os.WriteFile(fileName, prettyPayload, perm)
}
