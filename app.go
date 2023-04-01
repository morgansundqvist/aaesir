package main

import (
	"encoding/json"
	"io/ioutil"
	"morgansundqvist/aaesir/action"
	internaltypes "morgansundqvist/aaesir/internalTypes"
	"os"

	"github.com/AlecAivazis/survey/v2"
)

const (
	// AppFilePath is the path to the app file
	AppFilePath = "./maesir.json"
)

func SaveApp(app *internaltypes.AaesirApp) {
	fileContent, _ := json.MarshalIndent(app, "", " ")
	_ = ioutil.WriteFile(AppFilePath, fileContent, 0644)
}

func main() {

	filePath := AppFilePath
	app := internaltypes.AaesirApp{RecentRequests: []internaltypes.Request{}, SavedRequests: []internaltypes.SavedRequest{}, DataContext: []internaltypes.DataContextObject{}, ExecutionChains: []internaltypes.ExecutionChain{}}

	// Check if the filePath exists

	// Try to open the file
	file, err := os.Open(filePath)
	if err != nil {
		// If the file doesn't exist, create it
		if os.IsNotExist(err) {
			SaveApp(&app)
		} else {
			panic(err)
		}
	} else {
		fileContents, err := ioutil.ReadAll(file)
		if err != nil {
			panic(err)
		}

		// Unmarshal the JSON into a struct

		err = json.Unmarshal(fileContents, &app)
		if err != nil {
			panic(err)
		}

	}
	file.Close()

	var mainQs = []*survey.Question{
		{
			Name: "mainActionToExecute",
			Prompt: &survey.Select{
				Message: "What do you want to do?",
				Options: []string{"New request", "Saved requests", "Recent requests", "Data context", "Execution chains", "Exit"},
				Default: "New request",
			},
		},
	}

	var mainAnswers struct {
		MainActionToExecute string
	}

	for {
		survey.Ask(mainQs, &mainAnswers)

		if mainAnswers.MainActionToExecute == "Exit" {
			break
		} else if mainAnswers.MainActionToExecute == "New request" {
			action.NewRequest(&app)
		} else if mainAnswers.MainActionToExecute == "Saved requests" {
			action.SavedRequests(&app)
		} else if mainAnswers.MainActionToExecute == "Recent requests" {
			action.RecentRequests(&app)
		} else if mainAnswers.MainActionToExecute == "Data context" {
			action.DataContext(&app)
		} else if mainAnswers.MainActionToExecute == "Execution chains" {
			action.ExecutionChains(&app)
		}
		SaveApp(&app)
	}
}
