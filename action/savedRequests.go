package action

import (
	"fmt"
	internaltypes "morgansundqvist/aaesir/internalTypes"

	"github.com/AlecAivazis/survey/v2"
)

func SavedRequests(app *internaltypes.AaesirApp) {
	fmt.Println("")
	for _, savedRequest := range app.SavedRequests {
		fmt.Printf("%s\t%s\t\t%s\n", FixedLengthString(savedRequest.Name, 30), savedRequest.Request.Method, savedRequest.Request.URL)
	}

	fmt.Println("")

	//Ask user what to do with saved requests
	var savedRequestsQs = []*survey.Question{
		{
			Name: "savedRequestActionToExecute",
			Prompt: &survey.Select{
				Message: "What do you want to do?",
				Options: []string{"Execute request", "Delete request", "Exit"},
				Default: "Execute request",
			},
		},
	}

	var savedRequestsAnswers struct {
		SavedRequestActionToExecute string
	}

	survey.Ask(savedRequestsQs, &savedRequestsAnswers)

	switch savedRequestsAnswers.SavedRequestActionToExecute {
	case "Execute request":
		ExecuteSavedRequest(app)
	case "Delete request":
		DeleteSavedRequest(app)
	case "Exit":
		return

	}
}

func ExecuteSavedRequest(app *internaltypes.AaesirApp) {
	//Ask user which request to execute
	var savedRequestsQs = []*survey.Question{
		{
			Name: "savedRequestToExecute",
			Prompt: &survey.Select{
				Message: "Which request do you want to execute?",
				Options: func() []string {
					var options []string
					for _, savedRequest := range app.SavedRequests {
						options = append(options, savedRequest.Name)
					}
					return options
				}(),
			},
		},
	}

	var savedRequestsAnswers struct {
		SavedRequestToExecute string
	}

	survey.Ask(savedRequestsQs, &savedRequestsAnswers)

	//Find the request to execute
	for _, savedRequest := range app.SavedRequests {
		if savedRequest.Name == savedRequestsAnswers.SavedRequestToExecute {
			//Execute the request
			executeRequest(savedRequest.Request, app)
		}
	}
}

func DeleteSavedRequest(app *internaltypes.AaesirApp) {
	//Ask user which request to delete
	var savedRequestsQs = []*survey.Question{
		{
			Name: "savedRequestToDelete",
			Prompt: &survey.Select{
				Message: "Which request do you want to delete?",
				Options: func() []string {
					var options []string
					for _, savedRequest := range app.SavedRequests {
						options = append(options, savedRequest.Name)
					}
					return options
				}(),
			},
		},
	}

	var savedRequestsAnswers struct {
		SavedRequestToDelete string
	}

	survey.Ask(savedRequestsQs, &savedRequestsAnswers)

	//Find the request to delete
	for i, savedRequest := range app.SavedRequests {
		if savedRequest.Name == savedRequestsAnswers.SavedRequestToDelete {
			//Delete the request
			app.SavedRequests = append(app.SavedRequests[:i], app.SavedRequests[i+1:]...)
		}
	}
}
