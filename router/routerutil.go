package router

import "encoding/json"

type messageUtil struct {
}

type MessageFormat struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func NewMessageUtil() *messageUtil {
	return &messageUtil{}
}

func (msgUtil *messageUtil) BuildMessage(obj interface{}) interface{} {
	return &MessageFormat{
		Code:    0,
		Message: "complete",
		Data:    obj,
	}
}

func (msgUtil *messageUtil) BuildErrorMessage(code int, message string, obj interface{}) interface{} {
	return &MessageFormat{
		Code:    code,
		Message: message,
		Data:    obj,
	}
}

func (msgUtil *messageUtil) BuildMessageBytes(obj interface{}) []byte {
	data, _ := json.Marshal(msgUtil.BuildMessage(obj))
	return data
}

func (msgUtil *messageUtil) BuildErrorMessageBytes(code int, message string, obj interface{}) []byte {
	data, _ := json.Marshal(msgUtil.BuildErrorMessage(code, message, obj))
	return data
}

func (msgUtil *messageUtil) ParseMessageFromBytes(msg []byte) (code int, message string, obj interface{}) {
	messageData := &MessageFormat{}
	if err := json.Unmarshal(msg, &messageData); err != nil {
		return -1, err.Error(), nil
	} else {
		return messageData.Code, messageData.Message, messageData.Data
	}
}
