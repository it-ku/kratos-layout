package server

import (
	"fmt"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/transport/http"
	"google.golang.org/protobuf/proto"
	stdhttp "net/http"
	"time"
)

type HTTPError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e *HTTPError) Error() string {
	return fmt.Sprintf("HTTPError code: %d message: %s", e.Code, e.Message)
}

func errorEncoder(w stdhttp.ResponseWriter, r *stdhttp.Request, err error) {
	se := errors.FromError(err)
	codec, _ := http.CodecForRequest(r, "Accept")

	errBody := &HTTPError{
		Code:    int(se.Code),
		Message: se.Message,
	}

	body, err := codec.Marshal(errBody)
	if err != nil {
		w.WriteHeader(stdhttp.StatusInternalServerError)
		return
	}

	if se.Code > 99 && se.Code < 600 {
		w.WriteHeader(stdhttp.StatusOK)
	} else {
		w.WriteHeader(stdhttp.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/"+codec.Name())
	_, _ = w.Write(body)
}

type HTTPOk struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Ts      int64       `json:"ts"`
}

func respEncoder(w stdhttp.ResponseWriter, r *stdhttp.Request, v interface{}) error {
	codec, _ := http.CodecForRequest(r, "Accept")
	messageMap := make(map[string]interface{})
	messageStr, _ := codec.Marshal(v.(proto.Message))
	_ = codec.Unmarshal(messageStr, &messageMap)

	if len(messageMap) == 1 {
		for _, vv := range messageMap {
			v = vv
		}
	}

	reply := &HTTPOk{
		Code:    200,
		Message: "success",
		Data:    v,
		Ts:      time.Now().Unix(),
	}

	data, err := codec.Marshal(reply)
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "application/"+codec.Name())
	_, err = w.Write(data)
	if err != nil {
		return err
	}
	return nil
}
