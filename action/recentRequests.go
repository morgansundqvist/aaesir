package action

import (
	"fmt"
	internaltypes "morgansundqvist/aaesir/internalTypes"

	"github.com/AlecAivazis/survey/v2"
)

func RecentRequests(app *internaltypes.AaesirApp) {

	if len(app.RecentRequests) > 5 {
		app.RecentRequests = app.RecentRequests[1:6]
	}

	for i, request := range app.RecentRequests {
		request.ID = i + 1
		app.RecentRequests[i] = request
	}

	for _, request := range app.RecentRequests {
		fmt.Printf("%d\t%s\t%s\t\t%s\n", request.ID, request.Method, FixedLengthString(request.Status, 3), request.URL)
	}
	fmt.Println("")

	var recentRequestsQs = []*survey.Question{
		{
			Name: "recentRequestActionToExecute",
			Prompt: &survey.Select{
				Message: "What do you want to do?",
				Options: []string{"Execute request", "Save request", "Delete request", "Exit"},
				Default: "Execute request",
			},
		},
	}

	var recentRequestsAnswers struct {
		RecentRequestActionToExecute string
	}

	survey.Ask(recentRequestsQs, &recentRequestsAnswers)

	switch recentRequestsAnswers.RecentRequestActionToExecute {
	case "Execute request":
		ExecuteRecentRequest(app)
	case "Save request":
		SaveRecentRequest(app)
	case "Delete request":
		DeleteRecentRequest(app)
	case "Exit":
		return
	}

}

func SaveRecentRequest(app *internaltypes.AaesirApp) {
	// Define sruvey questions for new request
	var recentRequestsQs = []*survey.Question{
		{
			Name: "recentRequestToSave",
			Prompt: &survey.Input{
				Message: "What request do you want to save?",
			},
		},
		{
			Name: "nameOfSavedRequest",
			Prompt: &survey.Input{
				Message: "What do you want to name the request?",
			},
		},
		{
			Name: "overwriteSavedRequest",
			Prompt: &survey.Confirm{
				Message: "Do you want to overwrite the request if it already exists?",
				Default: false,
			},
		},
	}

	// Define struct for new request answers
	var recentRequestsAnswers struct {
		RecentRequestToSave   int
		NameOfSavedRequest    string
		OverwriteSavedRequest bool
	}

	// Ask questions
	survey.Ask(recentRequestsQs, &recentRequestsAnswers)

	if recentRequestsAnswers.RecentRequestToSave > len(app.RecentRequests) {
		fmt.Println("Invalid request")
		return
	}

	requestToSave := app.RecentRequests[recentRequestsAnswers.RecentRequestToSave-1]

	for i, savedRequest := range app.SavedRequests {
		if savedRequest.Name == recentRequestsAnswers.NameOfSavedRequest {
			if recentRequestsAnswers.OverwriteSavedRequest {
				app.SavedRequests[i].Request = requestToSave
				return
			} else {
				fmt.Println("Request already exists")
				return
			}
		}
	}

	app.SavedRequests = append(app.SavedRequests, internaltypes.SavedRequest{
		Name:    recentRequestsAnswers.NameOfSavedRequest,
		Request: requestToSave,
	})
}

func DeleteRecentRequest(app *internaltypes.AaesirApp) {
	// Define sruvey questions for new request
	var recentRequestsQs = []*survey.Question{
		{
			Name: "recentRequestToDelete",
			Prompt: &survey.Input{
				Message: "What request do you want to delete?",
			},
		},
	}

	// Define struct for new request answers
	var recentRequestsAnswers struct {
		RecentRequestToDelete int
	}

	// Ask questions
	survey.Ask(recentRequestsQs, &recentRequestsAnswers)

	if recentRequestsAnswers.RecentRequestToDelete > len(app.RecentRequests) {
		fmt.Println("Invalid request")
		return
	}

	app.RecentRequests = append(app.RecentRequests[:recentRequestsAnswers.RecentRequestToDelete-1], app.RecentRequests[recentRequestsAnswers.RecentRequestToDelete:]...)

}

func ExecuteRecentRequest(app *internaltypes.AaesirApp) {
	// Define sruvey questions for new request
	var recentRequestsQs = []*survey.Question{
		{
			Name: "recentRequestToExecute",
			Prompt: &survey.Input{
				Message: "What request do you want to execute?",
			},
		},
	}

	// Define struct for new request answers
	var recentRequestsAnswers struct {
		RecentRequestToExecute int
	}

	// Ask questions
	survey.Ask(recentRequestsQs, &recentRequestsAnswers)

	if recentRequestsAnswers.RecentRequestToExecute > len(app.RecentRequests) {
		fmt.Println("Invalid request")
		return
	}

	executeRequest(app.RecentRequests[recentRequestsAnswers.RecentRequestToExecute-1], app)

}
