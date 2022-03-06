package reflection

import (
	"fmt"
	"reflect"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

const (
	protoFieldMessage = "message"
	structFieldName   = "Message"
)

var (
	errNotAStruct           = fmt.Errorf("request object is not a struct")
	errImmutableField       = fmt.Errorf("request object field is immutable")
	errFieldNotFound        = fmt.Errorf("field not found in request object")
	errFieldNotOfTypeString = fmt.Errorf("field is not of type string")
	errNotAProto            = fmt.Errorf("request object is not a proto message")
)

// extractByReflection extract a field value from provided request using so reflection
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

// setByReflection set the provided string value to a string field in request using go reflection
func setByReflection(req interface{}, valueToSet string) error {
	reqType := reflect.TypeOf(req)
	// check if request is of kind pointer pointing to an element type struct
	if reqType.Kind() == reflect.Ptr && reqType.Elem().Kind() != reflect.Struct {
		return errNotAStruct
	}
	// check if struct type request has field protoFieldMessage
	field, exist := reqType.Elem().FieldByName(structFieldName)
	if !exist {
		return errFieldNotFound
	}
	// check if field is of type string
	if field.Type.Kind() != reflect.String {
		return errFieldNotOfTypeString
	}
	// extract value object of field
	reqTypeVal := reflect.ValueOf(req)
	// check if value object is from an exported field and can be modified
	if !reqTypeVal.Elem().FieldByName(structFieldName).CanSet() {
		return errImmutableField
	}
	// set the value
	reqTypeVal.Elem().FieldByName(structFieldName).SetString(valueToSet)
	return nil
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

// setByProtoReflection set provided string value to proto message request string field using proto reflect.
func setByProtoReflection(req interface{}, valueToSet string) error {
	// check if request is of type proto message
	reqMessage, ok := req.(proto.Message)
	if !ok {
		return errNotAProto
	}
	// check if message has field protoFieldMessage and is of kind string
	messageStringField := reqMessage.ProtoReflect().Descriptor().Fields().ByName(protoreflect.Name(protoFieldMessage))
	if messageStringField == nil {
		return errFieldNotFound
	}
	if messageStringField.Kind() != protoreflect.StringKind {
		return errFieldNotOfTypeString
	}
	reqMessage.ProtoReflect().Set(messageStringField, protoreflect.ValueOfString(valueToSet))
	return nil
}
