package httpapi

import (
        "fmt"
	"io/ioutil"
	"net/http"

	"github.com/tony-huhan/go-workwx/internal/lowlevel/envelope"
)

type EnvelopeHandler interface {
	OnIncomingEnvelope(rx envelope.Envelope) (body []byte, err error)
}

func (h *LowlevelHandler) eventHandler(
	rw http.ResponseWriter,
	r *http.Request,
) {
	// request bodies are assumed small
	// we can't do streaming parse/decrypt/verification anyway
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

        fmt.Printf("HandleIncomingMsg body msg is %v", body)
	// signature verification is inside EnvelopeProcessor
	ev, err := h.ep.HandleIncomingMsg(r.URL, body)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	msg, err := h.eh.OnIncomingEnvelope(ev)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
        fmt.Printf("OnIncomingEnvelope result msg is %v\n", string(msg))

	// if passive msg is not nil ,call inside MakeOutgoingEnvelope
	// then write to http response
	if msg != nil {
		resp, err := h.ep.MakeOutgoingEnvelope(msg)
                // fmt.Printf("MakeOutgoingEnvelope resp is %v\n", resp)
                fmt.Printf("MakeOutgoingEnvelope resp str is %s\n", string(resp))
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}
		rw.Write(resp)
	}

	// currently we always return empty 200 responses
	// any reply is to be sent asynchronously
	// this might change in the future (maybe save a couple of RTT or so)
	rw.WriteHeader(http.StatusOK)
}
