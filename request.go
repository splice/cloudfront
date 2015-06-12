package cloudfront

import "fmt"

type Request struct {
	URL     string
	Expires int64
}

func (req *Request) CannedPolicy() string {
	if req == nil {
		return ""
	}
	s := `{"Statement":[{"Resource":"%s","Condition":{"DateLessThan":{"AWS:EpochTime":%d}}}]}`
	return fmt.Sprintf(s, req.URL, req.Expires)
}
