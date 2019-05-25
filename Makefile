# Makefile for blitzinfod
# Beware: This uses tab intend (instead of spaces)

# vars #######################################################################
VERSION=1.0
P=infoblitz
D=infoblitz
PD=package/$(VERSION)

PREFIX=/usr/local


INSTALLED_FILES=\
  $(DESTDIR)$(PREFIX)/bin/infoblitz

ARCH := $(shell uname -m)

GOBUILD := go build -v
RICE := rice embed-go


# functions ##################################################################
do-listing=\
  echo; echo '[Listing installed files:]'; ls -l --color $(INSTALLED_FILES) || :


help:
	@echo 'Makefile for infoblitz Debian Packaging                       '
	@echo '                                                              '
	@echo 'Usage:                                                        '
	@echo '   make install                   installs infoblitz          '
	@echo '   make uninstall                 removes infoblitz           '
	@echo '   make build                     build                       '
	@echo '   make clean                     clean                       '
	@echo '   make list                      list installed files        '
	@echo '   make build                     build infoblitz from go src '
	@echo '   make package                   create .deb/.dcs packagings '
	@echo '   make package-clean             delete tmp packaging files  '
	@echo '   make package-info              apt-cache showpkg           '
	@echo '   make package-purge             apt-get purge               '
	@echo '   make package-version           update the changelog        '
	@echo '                                                              '


all:
	: # do nothing


install:
	@echo "### Install"
	install -d -m 755                     $(DESTDIR)$(PREFIX)/bin
	install -m 0755 blitzinfod-$(ARCH)    $(DESTDIR)$(PREFIX)/bin/blitzinfod-$(ARCH)


uninstall:
	@echo "### Uninstall"
	-rm -f $(DESTDIR)$(PREFIX)/bin/blitzinfod-$(ARCH)


# ARCH "armv7l"
build-armv7l:
	@echo "### Building armv7l"
	$(RICE)
	GOARCH=arm GOARM=7 $(GOBUILD) -o blitzinfod-armv7l server.go rice-box.go


# ARCH "x86_64"
build-x86_64:
	@echo "### Building x86_64"
	$(RICE)
	GOARCH=amd64 $(GOBUILD) -o blitzinfod-x86_64 server.go rice-box.go


build-all:
	$(RICE)
	@echo "### Building armv7l"
	GOARCH=arm GOARM=7 $(GOBUILD) -o blitzinfod-armv7l server.go rice-box.go
	@echo "### Building x86_64"
	GOARCH=amd64 $(GOBUILD) -o blitzinfod-x86_64 server.go rice-box.go


clean:
	: # do nothing


list:
	@$(do-listing)


package:
	@echo 'packaging blitzinfod for XXX Debian unstable...'
	mkdir -p $(PD)/unstable
	rm -f $(PD)/unstable/*
	sed -i 's,\(^$(P) ('$(VERSION)') \)[a-z]*,\1unstable,' debian/changelog
	debuild -I.git -Ipackage
	@echo ""
	cd .. && for f in $(PACKAGE_FILES); do \
		[ -f $$f ] && mv $$f $(D)/$(PD)/unstable || :; done


package-info:
	apt-cache showpkg $(P)
package-purge:
	apt-get purge $(P)
package-clean:
	dh_clean
package-version:
	dch -i


# phonies
.PHONY: install uninstall build clean list package
.PHONY: package-info package-purge package-clean package-version
