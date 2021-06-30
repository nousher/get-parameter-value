package getparameter

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

func (data *DataObj) readJsonFile(fileName *string) {

	if fileName == nil {
		log.Printf("Please provide the file name.")
	}

	byteValue, err := ioutil.ReadFile(*fileName)
	if err != nil {
		log.Printf("Unable to read data %s", err)
	}

	json.Unmarshal(byteValue, &data.DataInterface)
	if err != nil {
		log.Printf("Error in Unmarshalling Data %s", err)
	}

}
