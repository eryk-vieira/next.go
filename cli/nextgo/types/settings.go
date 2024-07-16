package types

type Settings struct {
	Package    string `json:"package"`
	RootFolder string `json:"root_folder"`
	Server     struct {
		Port string `json:"port"`
	} `json:"server"`
	HTTP struct {
		HandlerName string `json:"handler_name"`
	} `json:"http"`
}
