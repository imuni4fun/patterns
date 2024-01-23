
build:
	go build .

run:
	go run .

test:
	go test -v -count 1 .

trackB:
	go test -v -count 1 -run Duplicates . | grep ":B" > /tmp/tmpOutput
	@echo ""
	@echo "all output"
	@cat /tmp/tmpOutput
	@echo ""
	@echo "just main"
	cat /tmp/tmpOutput | grep "main:"
	@echo ""
	@echo "just start pull"
	cat /tmp/tmpOutput | grep "start pull:"
	@echo ""
	@echo "just puller"
	cat /tmp/tmpOutput | grep "puller:"
	@echo ""
	@echo "just wait for pull"
	cat /tmp/tmpOutput | grep "wait for pull:"
	