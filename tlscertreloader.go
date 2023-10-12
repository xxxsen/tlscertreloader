package tlscertreloader

import (
	"crypto/tls"
	"log"
	"sync/atomic"
	"time"
)

type CertReloader struct {
	keyFile  string
	certFile string

	c    *config
	cert atomic.Value
}

func NewCertReloader(cert string, key string, opts ...Option) (*CertReloader, error) {
	c := &config{
		period: 24 * time.Hour,
	}
	for _, opt := range opts {
		opt(c)
	}

	cr := &CertReloader{
		keyFile:  key,
		certFile: cert,
		c:        c,
	}
	if err := cr.init(); err != nil {
		return nil, err
	}
	return cr, nil
}

func MustNewCertReloader(cert string, key string, opts ...Option) *CertReloader {
	cr, err := NewCertReloader(cert, key, opts...)
	if err != nil {
		panic(err)
	}
	return cr
}

func (cr *CertReloader) GetCertificate(*tls.ClientHelloInfo) (*tls.Certificate, error) {
	v := cr.cert.Load()
	return v.(*tls.Certificate), nil
}

func (cr *CertReloader) init() error {
	if err := cr.reload(); err != nil {
		return err
	}
	go cr.periodicReload()
	return nil
}

func (cr *CertReloader) reload() error {
	cert, err := tls.LoadX509KeyPair(cr.certFile, cr.keyFile)
	if err != nil {
		return err
	}

	cr.cert.Store(&cert)
	return nil
}

func (cr *CertReloader) periodicReload() {
	for {
		time.Sleep(cr.c.period)
		if err := cr.reload(); err != nil {
			log.Printf("reload cert fail, err:%v", err)
		}
	}
}
