package marshaler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	nhk_service "github.com/nhk-news-web-easy/nhk-easy-service-proto"
	"google.golang.org/protobuf/proto"
	"io"
	"reflect"
	"strconv"
)

type protoEnum interface {
	fmt.Stringer
	EnumDescriptor() ([]byte, []int)
}

var typeProtoEnum = reflect.TypeOf((*protoEnum)(nil)).Elem()

var (
	// protoMessageType is stored to prevent constant lookup of the same type at runtime.
	protoMessageType = reflect.TypeOf((*proto.Message)(nil)).Elem()
)

type NewsMarshaler struct {
	runtime.JSONPb
}

// Serialize the NewsReply.News field only
func (newsMarshaler *NewsMarshaler) Marshal(v interface{}) ([]byte, error) {
	if newsReply, ok := v.(*nhk_service.NewsReply); ok {
		var buf bytes.Buffer
		if err := newsMarshaler.marshalTo(&buf, newsReply.News); err != nil {
			return nil, err
		}
		return buf.Bytes(), nil
	} else {
		return nil, nil
	}
}

// borrowed from marshal_jsonpb.go
// marshalNonProto marshals a non-message field of a protobuf message.
// This function does not correctly marshal arbitrary data structures into JSON,
// it is only capable of marshaling non-message field values of protobuf,
// i.e. primitive types, enums; pointers to primitives or enums; maps from
// integer/string types to primitives/enums/pointers to messages.
func (newsMarshaler *NewsMarshaler) marshalNonProtoField(v interface{}) ([]byte, error) {
	if v == nil {
		return []byte("null"), nil
	}
	rv := reflect.ValueOf(v)
	for rv.Kind() == reflect.Ptr {
		if rv.IsNil() {
			return []byte("null"), nil
		}
		rv = rv.Elem()
	}

	if rv.Kind() == reflect.Slice {
		if rv.IsNil() {
			if newsMarshaler.EmitUnpopulated {
				return []byte("[]"), nil
			}
			return []byte("null"), nil
		}

		if rv.Type().Elem().Implements(protoMessageType) {
			var buf bytes.Buffer
			err := buf.WriteByte('[')
			if err != nil {
				return nil, err
			}
			for i := 0; i < rv.Len(); i++ {
				if i != 0 {
					err = buf.WriteByte(',')
					if err != nil {
						return nil, err
					}
				}
				if err = newsMarshaler.marshalTo(&buf, rv.Index(i).Interface().(proto.Message)); err != nil {
					return nil, err
				}
			}
			err = buf.WriteByte(']')
			if err != nil {
				return nil, err
			}

			return buf.Bytes(), nil
		}

		if rv.Type().Elem().Implements(typeProtoEnum) {
			var buf bytes.Buffer
			err := buf.WriteByte('[')
			if err != nil {
				return nil, err
			}
			for i := 0; i < rv.Len(); i++ {
				if i != 0 {
					err = buf.WriteByte(',')
					if err != nil {
						return nil, err
					}
				}
				if newsMarshaler.UseEnumNumbers {
					_, err = buf.WriteString(strconv.FormatInt(rv.Index(i).Int(), 10))
				} else {
					_, err = buf.WriteString("\"" + rv.Index(i).Interface().(protoEnum).String() + "\"")
				}
				if err != nil {
					return nil, err
				}
			}
			err = buf.WriteByte(']')
			if err != nil {
				return nil, err
			}

			return buf.Bytes(), nil
		}
	}

	if rv.Kind() == reflect.Map {
		m := make(map[string]*json.RawMessage)
		for _, k := range rv.MapKeys() {
			buf, err := newsMarshaler.Marshal(rv.MapIndex(k).Interface())
			if err != nil {
				return nil, err
			}
			m[fmt.Sprintf("%v", k.Interface())] = (*json.RawMessage)(&buf)
		}
		if newsMarshaler.Indent != "" {
			return json.MarshalIndent(m, "", newsMarshaler.Indent)
		}
		return json.Marshal(m)
	}
	if enum, ok := rv.Interface().(protoEnum); ok && !newsMarshaler.UseEnumNumbers {
		return json.Marshal(enum.String())
	}
	return json.Marshal(rv.Interface())
}

// borrowed from marshal_jsonpb.go
func (newsMarshaler *NewsMarshaler) marshalTo(w io.Writer, v interface{}) error {
	p, ok := v.(proto.Message)
	if !ok {
		buf, err := newsMarshaler.marshalNonProtoField(v)
		if err != nil {
			return err
		}
		_, err = w.Write(buf)
		return err
	}
	b, err := newsMarshaler.MarshalOptions.Marshal(p)
	if err != nil {
		return err
	}

	_, err = w.Write(b)
	return err
}
