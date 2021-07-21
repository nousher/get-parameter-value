package getparameter

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

func (data *DataObj) readJsonFile(fileName *string) error {

	if fileName == nil {
		return fmt.Errorf("Please provide the file name.")
	}

	byteValue, err := ioutil.ReadFile(*fileName)
	if err != nil {
		return fmt.Errorf("Unable to read data %s", err)
	}

	json.Unmarshal(byteValue, &data.DataInterface)
	if err != nil {
		return fmt.Errorf("Error in Unmarshalling Data %s", err)
	}
	return nil
}
