package action

import (
	"fmt"
	internaltypes "morgansundqvist/aaesir/internalTypes"

	"github.com/AlecAivazis/survey/v2"
)

func ExecutionChains(app *internaltypes.AaesirApp) {
	fmt.Println("Execution Chains")
	for _, executionChain := range app.ExecutionChains {
		fmt.Println(executionChain)
	}
	fmt.Println("")

	var executionChainsQs = []*survey.Question{
		{
			Name: "executionChainActionToExecute",
			Prompt: &survey.Select{
				Message: "What do you want to do?",
				Options: []string{"Create chain", "Execute chain", "Delete chain", "Exit"},
				Default: "Execute chain",
			},
		},
	}

	var executionChainsAnswers struct {
		ExecutionChainActionToExecute string
	}

	survey.Ask(executionChainsQs, &executionChainsAnswers)

	switch executionChainsAnswers.ExecutionChainActionToExecute {
	case "Create chain":
		CreateExecutionChain(app)
	case "Execute chain":
		ExecuteExecutionChain(app)
	case "Delete chain":
		DeleteExecutionChain(app)
	case "Exit":
		return
	}
}

func ExecuteExecutionChain(app *internaltypes.AaesirApp) {
	//Ask user which request to execute
	var executionChainsQs = []*survey.Question{
		{
			Name: "executionChainToExecute",
			Prompt: &survey.Select{
				Message: "Which execution chain do you want to execute?",
				Options: func() []string {
					var options []string
					for _, executionChain := range app.ExecutionChains {
						options = append(options, executionChain.Name)
					}
					return options
				}(),
			},
		},
	}

	var executionChainsAnswers struct {
		ExecutionChainToExecute string
	}

	survey.Ask(executionChainsQs, &executionChainsAnswers)

	for _, executionChain := range app.ExecutionChains {
		if executionChain.Name == executionChainsAnswers.ExecutionChainToExecute {
			for _, executionChainItem := range executionChain.Items {
				for _, savedRequest := range app.SavedRequests {
					if savedRequest.Name == executionChainItem.SavedRequestName {
						executeRequest(savedRequest.Request, app)
					}
				}
			}
		}
	}
}

func DeleteExecutionChain(app *internaltypes.AaesirApp) {
	//Ask user which request to execute
	var executionChainsQs = []*survey.Question{
		{
			Name: "executionChainToDelete",
			Prompt: &survey.Select{
				Message: "Which execution chain do you want to delete?",
				Options: func() []string {
					var options []string
					for _, executionChain := range app.ExecutionChains {
						options = append(options, executionChain.Name)
					}
					return options
				}(),
			},
		},
	}

	var executionChainsAnswers struct {
		ExecutionChainToDelete string
	}

	survey.Ask(executionChainsQs, &executionChainsAnswers)

	for i, executionChain := range app.ExecutionChains {
		if executionChain.Name == executionChainsAnswers.ExecutionChainToDelete {
			app.ExecutionChains = append(app.ExecutionChains[:i], app.ExecutionChains[i+1:]...)
		}
	}
}

func CreateExecutionChain(app *internaltypes.AaesirApp) {
	// Define sruvey questions for new request
	var executionChainsQs = []*survey.Question{
		{
			Name: "executionChainName",
			Prompt: &survey.Input{
				Message: "What do you want to name the execution chain?",
			},
		},
	}

	var executionChainsAnswers struct {
		ExecutionChainName string
	}

	survey.Ask(executionChainsQs, &executionChainsAnswers)

	executionChain := internaltypes.ExecutionChain{Name: executionChainsAnswers.ExecutionChainName}

	for {
		var savedRequestQs = []*survey.Question{
			{
				Name: "savedRequest",
				Prompt: &survey.Select{
					Message: "Which saved request do you want to add to the execution chain?",
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

		var savedRequestAnswers struct {
			SavedRequest string
		}

		survey.Ask(savedRequestQs, &savedRequestAnswers)

		executionChain.Items = append(executionChain.Items, internaltypes.ExecutionChainItem{SavedRequestName: savedRequestAnswers.SavedRequest})

		addMore := false
		prompt := &survey.Confirm{
			Message: "Do you want to add more saved requests to the execution chain?",
		}

		survey.AskOne(prompt, &addMore)

		if !addMore {
			break
		}

	}
	app.ExecutionChains = append(app.ExecutionChains, executionChain)
}
