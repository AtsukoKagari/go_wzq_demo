package model

type Chess struct{
    xy		string		`json:"xy"`
    message	string		`json:"message"`
    bout	bool		`json:"bout"`
    color	string		`json:"color"`
}