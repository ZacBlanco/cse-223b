
OW_USER?=whisk
OW_RUNTIME?=$(OW_USER)/action-golang-v1.15-actor:latest
OW_COMPILER?=$(OW_USER)/action-golang-v1.15-actor:latest
WSK?=wsk
MAIN=main
PACKAGE=crawler
SRCS=main.go go.mod go.sum
NAME=module-main
BINZIP=$(MAIN)-bin.zip
SRCZIP=$(MAIN)-src.zip
ARGS_FILE=seeds.json

HOST=http://localhost:3233

build: $(SRCZIP)
	wsk -i --apihost $(HOST) action update web-crawler-main $(SRCZIP) --kind go:1.15-actor --main Main

run:
	wsk -i --apihost $(HOST) action invoke web-crawler-main --param-file $(ARGS_FILE) -b

runv:
	wsk -i --apihost $(HOST) action invoke web-crawler-main -b -v -d

runc:
	wsk -i --apihost $(HOST) action invoke web-crawler-c0 -b

$(BINZIP): $(SRCS) $(VENDORS) $(SRCZIP)
	docker run -i $(OW_COMPILER) -compile $(MAIN) -debug <$(SRCZIP) >$(BINZIP)

$(SRCZIP): $(SRCS)
	zip $@ -qr $^

clean:
	-$(WSK) action delete $(PACKAGE)/$(NAME)
	-rm $(BINZIP) $(SRCZIP) package.done test.json
	-rm test.out
	-rm *.zip

.PHONY: deploy devel test clean
