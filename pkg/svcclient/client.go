/*
Copyright 2022 Nokia.

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

package svcclient

import (
	"context"
	"crypto/tls"
	"fmt"
	"time"

	"github.com/henderiw-k8s-lcnc/fn-svc-sdk/pkg/api/fnservicepb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	defaultTimeout = 30 * time.Second
	maxMsgSize     = 512 * 1024 * 1024
)

type ServiceClient interface {
	Create() (fnservicepb.ServiceFunctionClient, error)
	Get() fnservicepb.ServiceFunctionClient
	Close()
}

func New(cfg *Config) (ServiceClient, error) {
	if cfg == nil {
		return nil, fmt.Errorf("must provide non-nil Configw")
	}
	c := &client{
		cfg: cfg,
	}

	return c, nil
}

type client struct {
	cfg    *Config
	conn   *grpc.ClientConn
	client fnservicepb.ServiceFunctionClient
}

func (r *client) Close() {
	if r.conn != nil {
		r.conn.Close()
	}
}

func (r *client) Get() fnservicepb.ServiceFunctionClient {
	return r.client
}

func (r *client) Create() (fnservicepb.ServiceFunctionClient, error) {
	if r.cfg == nil {
		return nil, fmt.Errorf("must provide non-nil Configw")
	}
	var opts []grpc.DialOption
	fmt.Printf("grpc client config: %v\n", r.cfg)
	if r.cfg.Insecure {
		//opts = append(opts, grpc.WithInsecure())
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	} else {
		tlsConfig, err := r.newTLS()
		if err != nil {
			return nil, err
		}
		opts = append(opts, grpc.WithTransportCredentials(credentials.NewTLS(tlsConfig)))
	}
	timeoutCtx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	var err error
	r.conn, err = grpc.DialContext(timeoutCtx, r.cfg.Address, opts...)
	if err != nil {
		return nil, err
	}
	//defer conn.Close()
	return fnservicepb.NewServiceFunctionClient(r.conn), nil
}

func (r *client) newTLS() (*tls.Config, error) {
	tlsConfig := &tls.Config{
		Renegotiation:      tls.RenegotiateNever,
		InsecureSkipVerify: r.cfg.SkipVerify,
	}
	//err := loadCerts(tlsConfig)
	//if err != nil {
	//	return nil, err
	//}
	return tlsConfig, nil
}

/*
func loadCerts(tlscfg *tls.Config) error {
	if c.TLSCert != "" && c.TLSKey != "" {
		certificate, err := tls.LoadX509KeyPair(*c.TLSCert, *c.TLSKey)
		if err != nil {
			return err
		}
		tlscfg.Certificates = []tls.Certificate{certificate}
		tlscfg.BuildNameToCertificate()
	}
	if c.TLSCA != nil && *c.TLSCA != "" {
		certPool := x509.NewCertPool()
		caFile, err := ioutil.ReadFile(*c.TLSCA)
		if err != nil {
			return err
		}
		if ok := certPool.AppendCertsFromPEM(caFile); !ok {
			return errors.New("failed to append certificate")
		}
		tlscfg.RootCAs = certPool
	}
	return nil
}
*/
