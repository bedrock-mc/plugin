.PHONY: proto proto-go proto-php proto-ts lint breaking clean

proto: proto-go proto-php proto-ts proto-py proto-cpp proto-rs

proto-go:
	buf generate && cd .. && go mod tidy

proto-php:
	buf generate --template buf.gen.php.yaml

proto-ts:
	buf generate --template buf.gen.ts.yaml && cd generated/ts && npm install && npm run build

proto-py:
	buf generate --template buf.gen.py.yaml

proto-cpp:
	buf generate --template buf.gen.cpp.yaml

proto-rs:
	buf generate --template buf.gen.rs.yaml

lint:
	buf lint

breaking:
	buf breaking --against '.git#branch=main'

clean:
	rm -rf generated/*.pb.go generated/php/ generated/ts/