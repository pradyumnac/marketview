.PHONY: all
INSTALL_PREFIX := $${HOME}/.local/bin

all: clean build install

clean:
	@rm -f marketview $(INSTALL_PREFIX)/marketview
	@echo "Clean Done"

build:
	go build .
	@echo "Build Done"

install: 
	@mv marketview $(INSTALL_PREFIX)/marketview
	@echo "Install Done"
