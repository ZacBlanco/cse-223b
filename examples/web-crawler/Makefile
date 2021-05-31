
WSK?=wsk
MAIN=main
SRCS=main.go go.mod go.sum
NAME=module-main
SRCZIP=$(MAIN)-src.zip
ARGS_FILE=seeds.json

HOST=http://172.17.0.1:3233

build: $(SRCZIP)
	wsk -i --apihost $(HOST) action update web-crawler-main $(SRCZIP) --kind go:1.15-actor --main Main

run:
	wsk -i --apihost $(HOST) action invoke web-crawler-main --param-file $(ARGS_FILE) -b

runv:
	wsk -i --apihost $(HOST) action invoke web-crawler-main -b -v -d

runc:
	wsk -i --apihost $(HOST) action invoke web-crawler-c0 -b

$(SRCZIP): $(SRCS)
	zip $@ -qr $^

clean:
	-rm test.out
	-rm *.zip

.PHONY: deploy devel test clean
