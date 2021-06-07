N=100
WSK?=wsk
MAIN=main
MAIN_CHILD=Mainchild
SRCS=main_$N.go go.mod go.sum
NAME=module-main
SRCZIP=$(MAIN)-src.zip
ARGS_FILE=seeds.json
TIMEOUT=240000 # 4 minutes

HOST=http://172.17.0.1:3233

build: $(SRCZIP)
	wsk -i --apihost $(HOST) action update web-crawler-main $(SRCZIP) --timeout $(TIMEOUT) --kind go:1.15-actor --main Main
	i=0; while [ "$$i" -le $$(($N - 1)) ]; do \
		wsk -i --apihost $(HOST) action update web-crawler-$$i $(SRCZIP) --timeout $(TIMEOUT) --kind go:1.15-actor --main ${MAIN_CHILD}$$i; \
		i=$$((i + 1));\
	done

run:
	wsk -i --apihost $(HOST) action invoke web-crawler-main --param-file $(ARGS_FILE) -b

runb:
	wsk -i --apihost $(HOST) action invoke web-crawler-main --param-file $(ARGS_FILE)

runv:
	wsk -i --apihost $(HOST) action invoke web-crawler-main -b -v -d

runc:
	wsk -i --apihost $(HOST) action invoke web-crawler-c0 -b

build_file:
	./build_main.sh $N

$(SRCZIP): build_file $(SRCS)
	zip $@ -qr $^

clean:
	-rm main_*.go
	-rm *.zip

.PHONY: deploy devel test clean
