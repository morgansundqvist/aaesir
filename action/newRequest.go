package action

import (
	"bytes"
	"encoding/json"
	"fmt"
	internaltypes "morgansundqvist/aaesir/internalTypes"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/go-resty/resty/v2"
	"github.com/simonnilsson/ask"
)

func NewRequest(app *internaltypes.AaesirApp) {
	fmt.Println("New request")
	// Define sruvey questions for new request
	var newRequestQs = []*survey.Question{
		{
			Name: "method",
			Prompt: &survey.Select{
				Message: "What method do you want to use?",
				Options: []string{"GET", "POST", "PUT", "DELETE"},
				Default: "GET",
			},
		},
		{
			Name: "url",
			Prompt: &survey.Input{
				Message: "What URL do you want to use?",
				Default: "http://127.0.0.1:8080/api",
				Suggest: func(toComplete string) []string {
					var suggestions []string
					for _, request := range app.RecentRequests {
						if strings.Contains(request.URL, toComplete) {
							suggestions = append(suggestions, request.URL)
						}
					}
					suggestions = RemoveDuplicates(suggestions)
					return suggestions
				},
			},
		},
	}

	// Define struct for new request answers
	var newRequestAnswers struct {
		Method string
		URL    string
	}

	// Ask questions
	survey.Ask(newRequestQs, &newRequestAnswers)

	newRequest := internaltypes.Request{Method: newRequestAnswers.Method, URL: newRequestAnswers.URL, Headers: make(map[string]string), Body: make(map[string]string), SavedResponseObjects: make(map[string]string)}

	addHeaders := false

	promptAddHeaders := &survey.Confirm{
		Message: "Do you want to add headers?",
		Default: false,
	}

	survey.AskOne(promptAddHeaders, &addHeaders)

	if addHeaders {
		for {
			// Define sruvey questions for new request header
			var newRequestHeaderQs = []*survey.Question{
				{
					Name: "headerKey",
					Prompt: &survey.Input{
						Message: "What key do you want to add to the headers?",
						Suggest: func(toComplete string) []string {
							var suggestions []string
							for _, request := range app.RecentRequests {
								for headerKey := range request.Headers {
									if strings.Contains(headerKey, toComplete) {
										suggestions = append(suggestions, headerKey)
									}
								}
							}
							suggestions = RemoveDuplicates(suggestions)
							return suggestions
						},
					},
				},
				{
					Name: "headerValue",
					Prompt: &survey.Input{
						Message: "What value do you want to add to the headers?",
						Suggest: func(toComplete string) []string {
							var suggestions []string
							for _, request := range app.RecentRequests {
								for _, headerValue := range request.Headers {
									if strings.Contains(headerValue, toComplete) {
										suggestions = append(suggestions, headerValue)
									}
								}
							}
							suggestions = RemoveDuplicates(suggestions)
							return suggestions
						},
					},
				},
			}

			// Define struct for new request header answers
			var newRequestHeaderAnswers struct {
				HeaderKey   string
				HeaderValue string
			}

			// Ask questions
			survey.Ask(newRequestHeaderQs, &newRequestHeaderAnswers)

			newRequest.Headers[newRequestHeaderAnswers.HeaderKey] = newRequestHeaderAnswers.HeaderValue

			addAnotherHeader := false

			promptAddAnotherHeader := &survey.Confirm{
				Message: "Do you want to add another header?",
				Default: false,
			}

			survey.AskOne(promptAddAnotherHeader, &addAnotherHeader)

			if !addAnotherHeader {
				break
			}
		}
	}

	addBody := false

	promptAddBody := &survey.Confirm{
		Message: "Do you want to add a body?",
		Default: false,
	}
	survey.AskOne(promptAddBody, &addBody)

	if addBody {
		for {
			// Define sruvey questions for new request body
			var newRequestBodyQs = []*survey.Question{
				{
					Name: "bodyKey",
					Prompt: &survey.Input{
						Message: "What key do you want to add to the json body?",
					},
				},
				{
					Name: "bodyValue",
					Prompt: &survey.Input{
						Message: "What value do you want to add to the json body?",
					},
				},
			}

			// Define struct for new request body answers
			var newRequestBodyAnswers struct {
				BodyKey   string
				BodyValue string
			}

			// Ask questions
			survey.Ask(newRequestBodyQs, &newRequestBodyAnswers)

			newRequest.Body[newRequestBodyAnswers.BodyKey] = newRequestBodyAnswers.BodyValue

			var addMoreBodyQs = []*survey.Question{
				{
					Name: "addMoreBody",
					Prompt: &survey.Confirm{
						Message: "Do you want to add more to the body?",
						Default: false,
					},
				},
			}

			var addMoreBodyAnswers struct {
				AddMoreBody bool
			}

			survey.Ask(addMoreBodyQs, &addMoreBodyAnswers)

			if !addMoreBodyAnswers.AddMoreBody {
				break
			}

		}
	}

	saveResponseObjects := false

	promptSaveResponseObjects := &survey.Confirm{
		Message: "Do you want to save the response objects?",
		Default: false,
	}

	survey.AskOne(promptSaveResponseObjects, &saveResponseObjects)

	if saveResponseObjects {
		for {
			// Define sruvey questions for new request body
			var newRequestSaveResponseObjectsQs = []*survey.Question{
				{
					Name: "responseKey",
					Prompt: &survey.Input{
						Message: "What key do you want to save from the response?",
					},
				},
				{
					Name: "responseValue",
					Prompt: &survey.Input{
						Message: "What value do you want to save from the response?",
					},
				},
			}

			// Define struct for new request body answers
			var newRequestSaveResponseObjectsAnswers struct {
				ResponseKey   string
				ResponseValue string
			}

			// Ask questions
			survey.Ask(newRequestSaveResponseObjectsQs, &newRequestSaveResponseObjectsAnswers)

			fmt.Printf("Response key: %s\n", newRequestSaveResponseObjectsAnswers.ResponseKey)

			newRequest.SavedResponseObjects[newRequestSaveResponseObjectsAnswers.ResponseKey] = newRequestSaveResponseObjectsAnswers.ResponseValue

			var addMoreResponseObjectsyQs = []*survey.Question{
				{
					Name: "addMoreBody",
					Prompt: &survey.Confirm{
						Message: "Do you want to add more responseobjects to save?",
						Default: false,
					},
				},
			}

			var addMoreResponseObjects struct {
				AddMoreResponseObjects bool
			}

			survey.Ask(addMoreResponseObjectsyQs, &addMoreResponseObjects)

			if !addMoreResponseObjects.AddMoreResponseObjects {
				break
			}
		}
	}

	executeRequest(newRequest, app)

}

func executeRequest(newRequest internaltypes.Request, app *internaltypes.AaesirApp) {
	client := resty.New()
	request := client.R()

	for key, value := range newRequest.Headers {
		value = ReplaceStringWithContextValue(value, app)

		request = request.SetHeader(key, value)
	}

	fmt.Printf("URL: %s\n", ReplaceStringWithContextValue(newRequest.URL, app))
	fmt.Printf("Method: %s\n", newRequest.Method)

	if newRequest.Method == "GET" {
		resp, err := request.Get(ReplaceStringWithContextValue(newRequest.URL, app))
		if err != nil {
			fmt.Printf("%s\n\n", err)

		} else {
			printBody(resp)
			AddResponseDataStoreInRecentRequests(newRequest, resp, app)
		}
	} else if newRequest.Method == "POST" {
		resp, err := request.SetBody(newRequest.Body).Post(ReplaceStringWithContextValue(newRequest.URL, app))
		if err != nil {
			fmt.Printf("%s\n\n", err)

		} else {
			//pretty print the json response
			printBody(resp)

			AddResponseDataStoreInRecentRequests(newRequest, resp, app)
		}
	}
}

func printBody(resp *resty.Response) {
	var prettyJSON bytes.Buffer
	error := json.Indent(&prettyJSON, resp.Body(), "", "\t")
	if error != nil {
		fmt.Println("JSON parse error: ", error)
	} else {
		fmt.Printf("Response Body: \n%v\n\n", prettyJSON.String())
	}
}

func AddResponseDataStoreInRecentRequests(newRequest internaltypes.Request, resp *resty.Response, app *internaltypes.AaesirApp) {
	newRequest.Response = resp.String()
	var object map[string]interface{}
	json.Unmarshal([]byte(newRequest.Response), &object)
	for key, value := range newRequest.SavedResponseObjects {
		res, ok := ask.For(object, value).String("error")
		if ok {
			//Check if a specific DataKey already exists
			existed := false
			for i, v := range app.DataContext {
				if v.DataKey == key {
					app.DataContext[i].DataValue = res
					existed = true
					break
				}
			}
			if !existed {
				app.DataContext = append(app.DataContext, internaltypes.DataContextObject{DataKey: key, DataValue: res})
			}
		} else {
			fmt.Println("Error could not find key in response object")
		}
	}
	fmt.Println("")
	newRequest.Status = resp.Status()
	app.RecentRequests = append(app.RecentRequests, newRequest)
}
