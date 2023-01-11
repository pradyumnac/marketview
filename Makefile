INSTALL_PREFIX := ~/.local/bin

.PHONY all
all: clean build install

clean:
	@rm -f marketview ~/$(INSTALL_PREFIX)/marketview
	@echo "Clean Done"

build:
	go build .
	@echo "Build Done"

install: 
	mv marketview ~/$(INSTALL_PREFIX)/marketview
	@echo "Install Done"
