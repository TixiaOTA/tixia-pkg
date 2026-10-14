package panorama

type FlightGetProviderListResultsEventData struct {
	Jsonrpc string      `json:"jsonrpc"`
	ID      interface{} `json:"id"`
	Result  struct {
		ErrorCode int    `json:"error_code"`
		ErrorMsg  string `json:"error_msg"`
		Response  struct {
			Providers []struct {
				Provider string   `json:"provider"`
				Carriers []string `json:"carriers"`
			} `json:"providers"`
		} `json:"response"`
	} `json:"result"`
}
