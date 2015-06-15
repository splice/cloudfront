# cloudfront signer

This signs cloudfront urls.

## Usage

```go
signer, err := cloudfront.NewSignerFromPath("./ssh/pk-APKA9ONS7QCOWEXAMPLE.pem", "APKA9ONS7QCOWEXAMPLE", "http://d111111abcdef8.cloudfront.net")
if err != nil {
	log.Fatal(err)
}

url, err := signer.SignedURL("/images/image.jpg", 5*time.Minute)
if err != nil {
	log.Fatal(err)
}
```
