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

package server

import (
	"context"
	"github.com/rs/xid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"io/ioutil"
	"os"
	"path"
)

const IDENTIFIER_FNAME  = "immudb.identifier"
const SERVER_UUID_HEADER  = "immudb-uuid"

type uuidContext struct {
	Uuid xid.ID
}

type UuidContext interface {
	UuidStreamContextSetter(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error
	UuidContextSetter(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error)
}

func NewUuidContext(id xid.ID) uuidContext{
	return uuidContext{id}
}

func getOrSetUuid(dir string) (xid.ID,error) {
	fname := path.Join(dir, IDENTIFIER_FNAME)
	if fileExists(fname){
		b, err := ioutil.ReadFile(fname)
		if err != nil {
			return xid.ID{}, err
		}
		return xid.FromBytes(b)
	}
	guid := xid.New()
	err := ioutil.WriteFile(fname, guid.Bytes(), os.ModePerm)
	return guid, err
	}

// fileExists checks if a file exists and is not a directory before we
// try using it to prevent further errors.
func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

type WrappedServerStream struct {
	grpc.ServerStream
}

func (w *WrappedServerStream) RecvMsg(m interface{}) error {
	return w.ServerStream.RecvMsg(m)
}

func (w *WrappedServerStream) SendMsg(m interface{}) error {
	return w.ServerStream.SendMsg(m)
}

func  (u *uuidContext) UuidStreamContextSetter(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	ss.Context()
	header := metadata.Pairs(SERVER_UUID_HEADER, u.Uuid.String())
	ss.SendHeader(header)
	return handler(srv, &WrappedServerStream{ss})
}

func  (u *uuidContext) UuidContextSetter(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	header := metadata.Pairs(SERVER_UUID_HEADER, u.Uuid.String())
	grpc.SendHeader(ctx, header)
	m, err := handler(ctx, req)
	return m, err
}