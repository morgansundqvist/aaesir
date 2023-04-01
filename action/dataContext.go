package action

//Function that prints all Data Context objects in the app

import (
	"fmt"
	internaltypes "morgansundqvist/aaesir/internalTypes"
	"regexp"
	"strings"

	"github.com/AlecAivazis/survey/v2"
)

func DataContext(app *internaltypes.AaesirApp) {
	fmt.Println("Data Context")
	for _, dataContextObject := range app.DataContext {
		fmt.Printf("%s : %s\n", dataContextObject.DataKey, TruncateString(dataContextObject.DataValue, 30))
	}
	fmt.Println("")

	var dataContextQs = []*survey.Question{
		{
			Name: "dataContextActionToExecute",
			Prompt: &survey.Select{
				Message: "What do you want to do?",
				Options: []string{"Add data context object", "Delete data context object", "Exit"},
				Default: "Add data context object",
			},
		},
	}

	var dataContextAnswers struct {
		DataContextActionToExecute string
	}

	survey.Ask(dataContextQs, &dataContextAnswers)

	switch dataContextAnswers.DataContextActionToExecute {
	case "Add data context object":
		AddDataContextObject(app)
	case "Delete data context object":
		DeleteDataContextObject(app)
	case "Exit":
		return
	}

}

func AddDataContextObject(app *internaltypes.AaesirApp) {
	var dataContextObjectQs = []*survey.Question{
		{
			Name: "dataKey",
			Prompt: &survey.Input{
				Message: "What is the data key?",
			},
		},
		{
			Name: "dataValue",
			Prompt: &survey.Input{
				Message: "What is the data value?",
			},
		},
	}

	var dataContextObjectAnswers struct {
		DataKey   string
		DataValue string
	}

	survey.Ask(dataContextObjectQs, &dataContextObjectAnswers)

	app.DataContext = append(app.DataContext, internaltypes.DataContextObject{
		DataKey:   dataContextObjectAnswers.DataKey,
		DataValue: dataContextObjectAnswers.DataValue,
	})

}

func DeleteDataContextObject(app *internaltypes.AaesirApp) {
	var dataContextObjectQs = []*survey.Question{
		{
			Name: "dataKey",
			Prompt: &survey.Select{
				Message: "Which data context object do you want to delete?",
				Options: func() []string {
					var options []string
					for _, dataContextObject := range app.DataContext {
						options = append(options, dataContextObject.DataKey)
					}
					return options
				}(),
			},
		},
	}

	var dataContextObjectAnswers struct {
		DataKey string
	}

	survey.Ask(dataContextObjectQs, &dataContextObjectAnswers)

	for i, dataContextObject := range app.DataContext {
		if dataContextObject.DataKey == dataContextObjectAnswers.DataKey {
			app.DataContext = append(app.DataContext[:i], app.DataContext[i+1:]...)
		}
	}
}

// Function which iterates over app.DataContext and finds a DataContextObject with a specific DataKey
func findDataContextObject(dataKey string, app *internaltypes.AaesirApp) internaltypes.DataContextObject {
	for _, dataContextObject := range app.DataContext {
		if dataContextObject.DataKey == dataKey {
			return dataContextObject
		}
	}

	return internaltypes.DataContextObject{}
}

func ReplaceStringWithContextValue(value string, app *internaltypes.AaesirApp) string {
	for _, match := range regexp.MustCompile(`\{([^}]+)\}`).FindAllString(value, -1) {
		variableName := strings.TrimPrefix(match, "{")
		variableName = strings.TrimSuffix(variableName, "}")
		dataContextObject := findDataContextObject(variableName, app)
		if dataContextObject.DataKey != "" {
			value = strings.Replace(value, match, dataContextObject.DataValue, -1)
		}
	}

	return value
}
