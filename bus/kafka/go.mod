module github.com/neutrinocorp/quark/bus/kafka

go 1.15

require (
	github.com/Shopify/sarama v1.28.0
	github.com/google/uuid v1.2.0 // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-multierror v1.1.1
	github.com/klauspost/compress v1.12.1 // indirect
	github.com/neutrinocorp/quark v1.0.0
	github.com/stretchr/testify v1.7.0
	golang.org/x/crypto v0.0.0-20210322153248-0c34fe9e7dc2 // indirect
	golang.org/x/net v0.0.0-20210414194228-064579744ee0 // indirect
)

replace github.com/neutrinocorp/quark => ../..
