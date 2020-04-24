/*
Copyright 2019-2020 vChain, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package gw

import (
	"context"
	"encoding/json"
	"github.com/codenotary/immudb/pkg/api/schema"
	"github.com/codenotary/immudb/pkg/client"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/grpc-ecosystem/grpc-gateway/utilities"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
	"net/http"
)

type SafesetHandler interface {
	Safeset(w http.ResponseWriter, req *http.Request, pathParams map[string]string)
}

type safesetHandler struct {
	mux    *runtime.ServeMux
	client client.ImmuClient
	rs     client.RootService
}

func NewSafesetHandler(mux *runtime.ServeMux, client client.ImmuClient, rs client.RootService) SafesetHandler {
	return &safesetHandler{
		mux:    mux,
		client: client,
		rs:     rs,
	}
}

func (h *safesetHandler) Safeset(w http.ResponseWriter, req *http.Request, pathParams map[string]string) {
	ctx, cancel := context.WithCancel(req.Context())
	defer cancel()
	inboundMarshaler, outboundMarshaler := runtime.MarshalerForRequest(h.mux, req)

	rctx, err := runtime.AnnotateContext(ctx, h.mux, req)
	if err != nil {
		runtime.HTTPError(ctx, h.mux, outboundMarshaler, w, req, err)
		return
	}

	var protoReq schema.SafeSetOptions
	var metadata runtime.ServerMetadata

	newReader, berr := utilities.IOReaderFactory(req.Body)
	if berr != nil {
		runtime.HTTPError(ctx, h.mux, outboundMarshaler, w, req, status.Errorf(codes.InvalidArgument, "%v", berr))
		return
	}
	if err := inboundMarshaler.NewDecoder(newReader()).Decode(&protoReq); err != nil && err != io.EOF {
		runtime.HTTPError(ctx, h.mux, outboundMarshaler, w, req, status.Errorf(codes.InvalidArgument, "%v", err))
		return
	}

	msg, err := h.client.SafeSet(rctx, protoReq.Kv.Key, protoReq.Kv.Value)
	if err != nil {
		runtime.HTTPError(ctx, h.mux, outboundMarshaler, w, req, err)
		return
	}

	ctx = runtime.NewServerMetadataContext(ctx, metadata)
	w.Header().Set("Content-Type", "application/json")

	newData, err := json.Marshal(msg)
	if err != nil {
		runtime.HTTPError(ctx, h.mux, outboundMarshaler, w, req, err)
		return
	}

	if _, err := w.Write(newData); err != nil {
		runtime.HTTPError(ctx, h.mux, outboundMarshaler, w, req, err)
		return
	}
	return
}
