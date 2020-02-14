package common

import (
	"crypto/tls"
	"encoding/base64"
)

type TLSConfig struct {
	Cert     string `yaml:"cert"`
	Key      string `yaml:"key"`
	CertFile string `yaml:"certFile"`
	KeyFile  string `yaml:"keyFile"`
}

func LoadTLSConfig(config *TLSConfig) (*tls.Config, error) {
	if config == nil {
		return nil, nil
	}
	var certs []tls.Certificate
	if config.Key != "" && config.Cert != "" {
		cert, err := base64.StdEncoding.DecodeString(config.Cert)
		if err != nil {
			return nil, err
		}
		key, err := base64.StdEncoding.DecodeString(config.Key)
		if err != nil {
			return nil, err
		}
		c, err := tls.X509KeyPair(cert, key)
		if err != nil {
			return nil, err
		}
		certs = append(certs, c)
	}

	if config.KeyFile != "" && config.CertFile != "" {
		cert, err := tls.LoadX509KeyPair(config.CertFile, config.KeyFile)
		if err != nil {
			return nil, err
		}
		certs = append(certs, cert)
	}

	if len(certs) > 0 {
		return &tls.Config{
			Certificates: certs,
		}, nil
	}

	return nil, nil
}
