package cloudfront

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"strings"
	"time"
)

// Signer holds on to the trusted signer's private key, and can sign urls.
// This doesn't bother with query params, because we don't need them.
type Signer struct {
	key     *rsa.PrivateKey
	ID      string
	baseURL string
	now     time.Time
}

// NewSigner loads the private key from the pem file at the provided path.
// In order to sign URLs, it also needs the id of the key pair, as well
// as the base url of the distribution.
func NewSigner(path, kpID, url string) (*Signer, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(b)
	pk, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	if err := pk.Validate(); err != nil {
		return nil, err
	}

	return &Signer{key: pk, ID: kpID, baseURL: url}, nil
}

// URL is the unsigned url for the given path.
func (s *Signer) URL(path string) string {
	return fmt.Sprintf("%s%s", s.baseURL, path)
}

// SignedURL provides a signed url valid until the provided expiration time.
// The path should be relative to the configured origin. E.g. if the origin is
// defined as s3-yourbucket/images, then you would omit "/images" from the path
// that you want to sign.
func (s *Signer) SignedURL(path string, expiration time.Duration) (string, error) {
	req := &Request{URL: s.URL(path), Expires: s.time(expiration).Unix()}
	sig, err := s.sig(req.CannedPolicy())
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s?Expires=%d&Signature=%s&Key-Pair-Id=%s", req.URL, req.Expires, sig, s.ID), nil
}

func (s *Signer) sig(policy string) (string, error) {
	hash := sha1.New()
	_, err := hash.Write([]byte(policy))
	if err != nil {
		return "", err
	}
	signed, err := rsa.SignPKCS1v15(nil, s.key, crypto.SHA1, hash.Sum(nil))
	if err != nil {
		return "", err
	}

	var r = strings.NewReplacer("=", "_", "+", "-", "/", "~")
	return r.Replace(base64.StdEncoding.EncodeToString(signed)), nil
}

// SetTime lets you fix "now" to a predetermined time.
// This is only useful for debugging.
func (s *Signer) SetTime(t time.Time) {
	if s == nil {
		return
	}
	s.now = t
}

func (s *Signer) time(d time.Duration) time.Time {
	if s != nil && !s.now.IsZero() {
		return s.now.Add(d)
	}
	return time.Now().Add(d)
}
