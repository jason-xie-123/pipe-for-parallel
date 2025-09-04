package main

type Package struct {
	UUID    string `json:"uuid"`
	Action  string `json:"action"`
	Message string `json:"message"`
}

type ResponsePackage struct {
	UUID string `json:"uuid"`
}
