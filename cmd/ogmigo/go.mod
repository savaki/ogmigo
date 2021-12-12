module blah

go 1.17

require (
	github.com/savaki/ogmigo v0.7.1
	github.com/urfave/cli/v2 v2.3.0
)

require (
	github.com/aws/aws-sdk-go v1.42.22 // indirect
	github.com/buger/jsonparser v1.1.1 // indirect
	github.com/cpuguy83/go-md2man/v2 v2.0.0-20190314233015-f79a8a8ca69d // indirect
	github.com/fxamacker/cbor v1.5.1 // indirect
	github.com/gorilla/websocket v1.4.2 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/russross/blackfriday/v2 v2.0.1 // indirect
	github.com/shurcooL/sanitized_anchor_name v1.0.0 // indirect
	github.com/x448/float16 v0.8.4 // indirect
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c // indirect
)

replace github.com/savaki/ogmigo/store/badgerstore => ../../store/badgerstore

replace github.com/savaki/ogmigo => ../..
