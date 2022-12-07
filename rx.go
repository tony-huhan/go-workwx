package workwx

import (
        "fmt"
	"net/http"

	"github.com/tony-huhan/go-workwx/internal/lowlevel/envelope"
	"github.com/tony-huhan/go-workwx/internal/lowlevel/httpapi"
)

// RxMessageHandler 用来接收消息的接口。
type RxMessageHandler interface {
	// OnIncomingMessage 一条消息到来时的回调。
	OnIncomingMessage(msg *RxMessage) ([]byte, error)
}

type lowlevelEnvelopeHandler struct {
	highlevelHandler RxMessageHandler
}

var _ httpapi.EnvelopeHandler = (*lowlevelEnvelopeHandler)(nil)

func (h *lowlevelEnvelopeHandler) OnIncomingEnvelope(rx envelope.Envelope) (body []byte, err error) {
	msg, err := fromEnvelope(rx.Msg)
	if err != nil {
		return nil, err
	}
        body, err = h.highlevelHandler.OnIncomingMessage(msg)
        fmt.Printf("rxMessageHandler err is %v, body string is \n", err, string(body))
	return body,err
}

type HTTPHandler struct {
	inner *httpapi.LowlevelHandler
}

var _ http.Handler = (*HTTPHandler)(nil)

func NewHTTPHandler(
	token string,
	encodingAESKey string,
	rxMessageHandler RxMessageHandler,
) (*HTTPHandler, error) {
	lleh := &lowlevelEnvelopeHandler{
		highlevelHandler: rxMessageHandler,
	}

	llHandler, err := httpapi.NewLowlevelHandler(token, encodingAESKey, lleh)
	if err != nil {
		return nil, err
	}

	obj := HTTPHandler{
		inner: llHandler,
	}

	return &obj, nil
}

func (h *HTTPHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	h.inner.ServeHTTP(rw, r)
}
