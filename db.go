package main

type PublicProperty struct {
	RdsServers    string `json:"RdsServers"`
	RdsUser       string `json:"RdsUser"`
	RdsPasswd     string `json:"RdsPasswd"`
	SrServers     string `json:"SrServers"`
	SrUser        string `json:"SrUser"`
	SrPasswd      string `json:"SrPasswd"`
	SourceELB     string `json:"SourceELB"`
	AwsAK         string `json:"AwsAK"`
	AwsSK         string `json:"AwsSK"`
	MQServer      string `json:"MQServer"`
	MQTopics      string `json:"MQTopics"`
	MQBroker      string `json:"MQBroker"`
	RedisServer   string `json:"RedisServer"`
	RedisDBNumber int32 `json:"RedisDBNumber"`
	ETCDServer    string `json:"ETCDServer"`
}