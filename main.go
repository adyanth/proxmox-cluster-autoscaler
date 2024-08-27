/*
Copyright 2022 The Kubernetes Authors.

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

package main

import (
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/adyanth/proxmox-cluster-autoscaler/proxmox"
	"github.com/adyanth/proxmox-cluster-autoscaler/wrapper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"k8s.io/autoscaler/cluster-autoscaler/cloudprovider/externalgrpc/protos"
	"k8s.io/klog/v2"
)

// MultiStringFlag is a flag for passing multiple parameters using same flag
type MultiStringFlag []string

// String returns string representation of the node groups.
func (flag *MultiStringFlag) String() string {
	return "[" + strings.Join(*flag, " ") + "]"
}

// Set adds a new configuration.
func (flag *MultiStringFlag) Set(value string) error {
	*flag = append(*flag, value)
	return nil
}

var (
	// flags needed by the external grpc provider service
	address     = flag.String("address", ":8086", "The address to expose the grpc service.")
	keyCert     = flag.String("key-cert", "", "The path to the certificate key file. Empty string for insecure communication.")
	cert        = flag.String("cert", "", "The path to the certificate file. Empty string for insecure communication.")
	cacert      = flag.String("ca-cert", "", "The path to the ca certificate file. Empty string for insecure communication.")
	cloudConfig = flag.String("config", "", "The path to the proxmox autoscaler config.")
)

func main() {
	klog.InitFlags(nil)
	flag.Parse()
	var s *grpc.Server

	// tls config
	var serverOpt grpc.ServerOption
	if *keyCert == "" || *cert == "" || *cacert == "" {
		fmt.Println("no cert specified, using insecure")
		s = grpc.NewServer()
	} else {

		certificate, err := tls.LoadX509KeyPair(*cert, *keyCert)
		if err != nil {
			panic(fmt.Errorf("failed to read certificate files: %s", err))
		}
		certPool := x509.NewCertPool()
		bs, err := os.ReadFile(*cacert)
		if err != nil {
			panic(fmt.Errorf("failed to read client ca cert: %s", err))
		}
		ok := certPool.AppendCertsFromPEM(bs)
		if !ok {
			panic(fmt.Errorf("failed to append client certs"))
		}
		transportCreds := credentials.NewTLS(&tls.Config{
			ClientAuth:   tls.RequireAndVerifyClientCert,
			Certificates: []tls.Certificate{certificate},
			ClientCAs:    certPool,
		})
		serverOpt = grpc.Creds(transportCreds)
		s = grpc.NewServer(serverOpt)
	}
	cloudProvider := proxmox.BuildProxmoxEngine(*cloudConfig)
	srv := wrapper.NewCloudProviderGrpcWrapper(cloudProvider)

	// listen
	lis, err := net.Listen("tcp", *address)
	if err != nil {
		panic(fmt.Errorf("failed to listen: %s", err))
	}

	// serve
	protos.RegisterCloudProviderServer(s, srv)
	fmt.Printf("Server ready at: %s\n", *address)
	if err := s.Serve(lis); err != nil {
		panic(fmt.Errorf("failed to serve: %v", err))
	}

}
