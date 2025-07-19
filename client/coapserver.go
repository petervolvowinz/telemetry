package main

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/plgd-dev/go-coap/v3"
	"github.com/plgd-dev/go-coap/v3/message"
	"github.com/plgd-dev/go-coap/v3/message/codes"
	"github.com/plgd-dev/go-coap/v3/mux"
	"log"
	dp "telemetry/datapoint"
)

func runCoapServer1(cp *dp.CurrentDataPoint, ctx context.Context) {
	router := mux.NewRouter()

	loggingMiddleWare := func(next mux.Handler) mux.Handler {
		return mux.HandlerFunc(func(w mux.ResponseWriter, r *mux.Message) {
			log.Printf("ClientAddress %v, %v\n", w.Conn().RemoteAddr(), r.String())
			next.ServeCOAP(w, r)
		})
	}

	hfunc := func(w mux.ResponseWriter, r *mux.Message) {
		dp := cp.GetCurrentDataPoint()

		payload, err := json.Marshal(dp)
		if err != nil {
			w.SetResponse(codes.InternalServerError, message.TextPlain, bytes.NewReader([]byte("marshal error")))
			return
		}

		w.SetResponse(codes.Content, message.AppJSON, bytes.NewReader(payload))
	}
	router.Use(loggingMiddleWare)
	for {
		select {
		case <-ctx.Done():
			return
		default:
			router.Handle("/bike-pw/telemetry", mux.HandlerFunc(hfunc))
			log.Println("serving coap on 5688")
			// listen to both Ipv4 and 6, resorts to a dual stack socket
			log.Fatal(coap.ListenAndServe("udp", ":5688", router))
		}
	}

}
