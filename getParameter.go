package getparameter

import (
	"fmt"

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

func GetParameterValue(tempFileName string) (map[string]string, error) {
	var parameterData *map[string]string

	if len(tempFileName) == 0 {
		return nil, fmt.Errorf("Please provide the config file name.")
	}

	fmt.Println("Getting Parameters from the Parameter Store")

	var interfaceObj DataObj

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	// if err != nil {
	// 	return nil, fmt.Errorf("Error in creating new session. Error %s", err)
	// }

	interfaceObj.readJsonFile(&tempFileName)
	if interfaceObj.DataInterface.(map[string]interface{}) == nil {
		return nil, fmt.Errorf("Sorry there is no data to process")
	} else {
		parameterData = setData(&interfaceObj)
	}

	if parameterData == nil {
		return nil, fmt.Errorf("Sorry there is no parameterData")
	}

	resultMap, err := getValue(sess, parameterData)
	if err != nil {
		return nil, err
	}

	return resultMap, nil

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

func getValue(sess *session.Session, parameterValue *map[string]string) (map[string]string, error) {

	var configData ConfigValue

	configData.ConfigMap = make(map[string]string)

	for key, parameter := range *parameterValue {

		fmt.Printf("key %v | value %v \n\n", key, parameter)

		svc := ssm.New(sess)

		outputValue, err := svc.GetParameter(
			&ssm.GetParameterInput{
				Name:           &parameter,
				WithDecryption: aws.Bool(true),
			},
		)
		if err != nil {
			return nil, fmt.Errorf("The error from the Get Paramter is %s", err)
		}
		tempVal1 := aws.StringValue(outputValue.Parameter.Value)
		configData.ConfigMap[key] = tempVal1

	}
	return configData.ConfigMap, nil
}
