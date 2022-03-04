package reflection

import (
	"reflect"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

const (
	protoFieldMessage = "message"
	structFieldName   = "Message"
)

func extractByReflection(req interface{}) (string, bool) {
	reqType := reflect.TypeOf(req)
	// check if request is of kind pointer pointing to an element type struct
	if reqType.Kind() == reflect.Ptr && reqType.Elem().Kind() != reflect.Struct {
		return "", false
	}
	// check if struct type request has field protoFieldMessage of kind string
	// if field value is not empty, return the value, true
	if field, exist := reqType.Elem().FieldByName(structFieldName); exist && field.Type.Kind() == reflect.String {
		reqTypeVal := reflect.ValueOf(req)
		messageString := reqTypeVal.Elem().FieldByName(structFieldName).String()
		if messageString != "" {
			return messageString, true
		}
	}
	return "", false
}

func extractByProtoReflection(req interface{}) (string, bool) {
	// check if request is of type proto message
	reqMessage, ok := req.(proto.Message)
	if !ok {
		return "", false
	}
	// check if message has field protoFieldMessage and is of kind string
	messageStringField := reqMessage.ProtoReflect().Descriptor().Fields().ByName(protoreflect.Name(protoFieldMessage))
	if messageStringField != nil && messageStringField.Kind() == protoreflect.StringKind {
		// get the field value and cast to string
		fieldValue := reqMessage.ProtoReflect().Get(messageStringField).String()
		if fieldValue != "" {
			return fieldValue, true
		}
	}
	return "", false
}
