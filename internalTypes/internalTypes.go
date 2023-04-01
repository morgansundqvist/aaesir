package internaltypes

// Struct which represents a request
type Request struct {
	ID int `json:"id"`
	//Method of the request
	Method string `json:"method"`
	//URL of the request
	URL string `json:"url"`
	//Headers of the request
	Headers map[string]string `json:"headers"`
	//Body of the request
	Body map[string]string `json:"body"`
	//Response of the request
	Response string `json:"response"`
	//Status of the request
	Status string `json:"status"`

	SavedResponseObjects map[string]string `json:"savedResponseObjects"`
}

type SavedRequest struct {
	ID      int     `json:"id"`
	Name    string  `json:"name"`
	Request Request `json:"request"`
}

type DataContextObject struct {
	DataKey   string `json:"dataKey"`
	DataValue string `json:"dataValue"`
}

type ExecutionChain struct {
	Name  string               `json:"name"`
	Items []ExecutionChainItem `json:"items"`
}

type ExecutionChainItem struct {
	SavedRequestName string `json:"savedRequestName"`
}

type AaesirApp struct {
	RecentRequests  []Request           `json:"recentRequests"`
	SavedRequests   []SavedRequest      `json:"savedRequests"`
	DataContext     []DataContextObject `json:"dataContext"`
	ExecutionChains []ExecutionChain    `json:"executionChains"`
}
