package models

type ServerResponse struct {
	Error bool        `json:"error"`
	Msg   string      `json:"msg"`
	Data  interface{} `json:"data"`
}
