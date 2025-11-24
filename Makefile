.PHONY: run proto

run:
	cd cmd && go run .
proto:
	cd proto && buf generate
	./scripts/post_generation.sh

tag-php-sdk:
	@if [ -z "$(VERSION)" ]; then echo "Usage: make tag-php-sdk VERSION=X.Y.Z"; exit 1; fi
	git tag php-sdk-v$(VERSION)
	git push origin php-sdk-v$(VERSION)
