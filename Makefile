# The import path is where your repository can be found.
# To import subpackages, always prepend the full import path.
# If you change this, run `make clean`. Read more: https://git.io/vM7zV
IMPORT_PATH := log

V := 1 # When V is set, print commands and build progress.

export GOPATH := $(CURDIR)/.GOPATH
unexport GOBIN

all: log

log: .GOPATH/.ok
	$Q go install -tags netgo $(IMPORT_PATH)

update: .GOPATH/.ok
	$Q glide mirror set https://golang.org/x/crypto https://github.com/golang/crypto
	$Q glide mirror set https://golang.org/x/text https://github.com/golang/text
	$Q glide mirror set https://golang.org/x/sys https://github.com/golang/sys
	$Q glide mirror set https://golang.org/x/net https://github.com/golang/net
	$Q glide up -v
	
clean:
	$Q rm -rf bin pkg .GOPATH

test: .GOPATH/.ok
	$Q cd $(GOPATH)/src/$(IMPORT_PATH);go test -v -timeout=30m
	
Q := $(if $V,,@)

.GOPATH/.ok:
	$Q rm -rf $(GOPATH)
	$Q mkdir -p $(GOPATH)/src
	$Q ln -sf $(CURDIR) $(GOPATH)/src/$(IMPORT_PATH)
	$Q mkdir -p $(CURDIR)/pkg
	$Q ln -sf $(CURDIR)/pkg $(GOPATH)/pkg
	$Q mkdir -p $(CURDIR)/bin
	$Q ln -sf $(CURDIR)/bin $(GOPATH)/bin
	$Q touch $@