package cloudfront

import "fmt"

// Request is a URL that will expire.
type Request struct {
	URL     string
	Expires int64
}

// CannedPolicy is the default Cloudfront policy for expiring URLs.
func (req *Request) CannedPolicy() string {
	if req == nil {
		return ""
	}
	s := `{"Statement":[{"Resource":"%s","Condition":{"DateLessThan":{"AWS:EpochTime":%d}}}]}`
	return fmt.Sprintf(s, req.URL, req.Expires)
}
