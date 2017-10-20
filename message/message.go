package message

type ApplicationMessage struct {
	Title   string   `json:"title"`
	MType   string   `json:"mtype"`   //dingidng, kafka, influx
	Content string   `json:"content"` // 具体数据信息，json格式
	At      []string `json:"at"`
	Token   string   `json:"token,omitempty"`
}
