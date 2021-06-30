package getparameter

import (
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
)

type DataObj struct {
	DataInterface interface{}
}

type ConfigValue struct {
	ConfigMap map[string]string
}

func GetParameterValue(tempFileName string) *map[string]string {
	var parameterData *map[string]string

	if len(tempFileName) == 0 {
		log.Printf("Please provide the config file name.")
		os.Exit(1)
	}

	fmt.Println("Program has been started!!!!")

	var interfaceObj DataObj

	sess, err := session.NewSession()
	if err != nil {
		log.Printf("Error in creating new session. Error %s", err)
	}

	interfaceObj.readJsonFile(&tempFileName)

	if interfaceObj.DataInterface.(map[string]interface{}) == nil {
		log.Printf("Sorry there is no data to process")
	} else {
		parameterData = setData(&interfaceObj)
	}

	if parameterData == nil {
		log.Printf("Sorry there is no parameterData")
		os.Exit(1)
	}

	resultMap := getValue(sess, parameterData)

	return resultMap

}

func setData(interfaceObj *DataObj) *map[string]string {

	data := interfaceObj.DataInterface.(map[string]interface{})
	parameterPath := make(map[string]string)

	for key, valueStage1 := range data {

		for _, ParameterValue := range valueStage1.([]interface{}) {
			stringValue := fmt.Sprintf("%s", ParameterValue)
			tempValue := "/" + "dev" + "/" + key + "/" + stringValue

			parameterPath[key+"/"+stringValue] = tempValue

		}
	}
	return &parameterPath
}

func getValue(sess *session.Session, parameterValue *map[string]string) *map[string]string {

	var configData ConfigValue

	configData.ConfigMap = make(map[string]string)

	for key, parameter := range *parameterValue {

		svc := ssm.New(sess)

		outputValue, err := svc.GetParameter(
			&ssm.GetParameterInput{
				Name:           &parameter,
				WithDecryption: aws.Bool(true),
			},
		)
		if err != nil {
			log.Printf("The error from the Get Paramter is %s, err:", err)
		}
		tempVal1 := aws.StringValue(outputValue.Parameter.Value)
		configData.ConfigMap[key] = tempVal1

	}
	return &configData.ConfigMap
}
