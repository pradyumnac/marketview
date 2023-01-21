.PHONY: all
INSTALL_PREFIX := $${HOME}/.local/bin
PWD := $(shell pwd)

all: clean build install

clean:
	@rm -f marketview $(INSTALL_PREFIX)/marketview
	@# @rm -fR ~/.local/share/marketview
	@rm -f ~/.local/share/marketview/symbols/data.db
	@rm -f ~/.local/share/marketview/data.db
	@echo "Clean Done"

build:
	@go build -ldflags "-s -w" .
	# @go build  .
	@echo "Build Done"

install: 
	@mv marketview $(INSTALL_PREFIX)/marketview
	@ln -s -f $(PWD)/bse_cd.sh $(INSTALL_PREFIX)/bse_cd
	@ln -s -f $(PWD)/nse_cd.sh $(INSTALL_PREFIX)/nse_cd
	@echo "Install Done"
