
build:
	wsk -i --apihost http://localhost:3233 action update web-crawler-main main.go --kind go:1.15-actor --main Main
	wsk -i --apihost http://localhost:3233 action update web-crawler-c0 main.go --kind go:1.15-actor --main MainChild1
	wsk -i --apihost http://localhost:3233 action update web-crawler-c1 main.go --kind go:1.15-actor --main MainChild2

buildc:
	wsk -i --apihost http://localhost:3233 action update web-crawler-main exec.zip --kind go:1.15-actor --main Main
	wsk -i --apihost http://localhost:3233 action update web-crawler-c0 exec.zip --kind go:1.15-actor --main MainChild1
	wsk -i --apihost http://localhost:3233 action update web-crawler-c1 exec.zip --kind go:1.15-actor --main MainChild2

# https://www.ibm.com/cloud/blog/running-existing-app-ibm-cloud-functions
# broken
mybuild:
	env GOOS=linux GOARCH=amd64 go build -o exec
	zip exec.zip exec
	wsk -i --apihost http://localhost:3233 action update web-crawler-main exec.zip --kind go:1.15-actor --main Main

run:
	wsk -i --apihost http://localhost:3233 action invoke web-crawler-main -r

runv:
	wsk -i --apihost http://localhost:3233 action invoke web-crawler-main -r -v -d

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

deploy: package.done $(BINZIP)
	$(WSK) action update $(PACKAGE)/$(NAME) $(BINZIP) --main $(MAIN) --docker $(OW_RUNTIME)

devel: package.done $(SRCZIP)
	$(WSK) action update $(PACKAGE)/$(NAME) $(SRCZIP) --main $(MAIN) --docker $(OW_COMPILER)

$(BINZIP): $(SRCS) $(VENDORS) $(SRCZIP)
	docker run -i $(OW_COMPILER) -compile $(MAIN) -debug <$(SRCZIP) >$(BINZIP)

$(SRCZIP): $(SRCS)
	zip $@ -qr $^

clean:
	-$(WSK) action delete $(PACKAGE)/$(NAME)
	-rm $(BINZIP) $(SRCZIP) package.done test.json
	-rm test.out

package.done:
	$(WSK) package update $(PACKAGE)
	touch package.done

.PHONY: deploy devel test clean