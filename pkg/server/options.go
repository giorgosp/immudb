/*
Copyright 2019 vChain, Inc.

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

type Options struct {
	Dir     string
	Network string
	Address string
	Port    int16
}

func DefaultOptions() Options {
	return Options{
		Dir:     ".",
		Network: "tcp",
		Address: "127.0.0.1",
		Port:    8080,
	}
}

func (o Options) WithDir(dir string) Options {
	o.Dir = dir
	return o
}

func (o Options) WithNetwork(network string) Options {
	o.Network = network
	return o
}

func (o Options) WithAddress(address string) Options {
	o.Address = address
	return o
}

func (o Options) WithPort(port int16) Options {
	o.Port = port
	return o
}