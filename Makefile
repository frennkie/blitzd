# Makefile for blitzinfod
# Heavily inspired by https://github.com/hamishcunningham/wiringpi
# Beware: This uses tab intend (instead of spaces)

# vars #######################################################################
VERSION=1.2
P=blitzinfod
D=blitzinfod
PD=package/$(VERSION)

PREFIX=/usr/local

# Target/Destination Host
#DEB_HOST_ARCH  := $(shell dpkg-architecture -qDEB_HOST_ARCH)

# Building Host
#DEB_BUILD_ARCH := $(shell dpkg-architecture -qDEB_BUILD_ARCH)


INSTALLED_FILES=\
  $(DESTDIR)$(PREFIX)/bin/infoblitz

PACKAGE_FILES=\
  $(P)_$(VERSION)*.deb \
  $(P)_$(VERSION)*.dsc \
  $(P)_$(VERSION)*.tar.gz \
  $(P)_$(VERSION)*.changes \
  $(P)_$(VERSION)*.build

MY_ARCH := $(shell uname -m)

GOBUILD := go build -v
RICE := /home/robbie/work/go/bin/rice embed-go

RICE_BIN=/home/robbie/work/go/bin/rice


# functions ##################################################################
do-listing=\
  echo; echo '[Listing installed files:]'; ls -l --color $(INSTALLED_FILES) || :


#########################################################
# by default `make` will build first target in Makefile #
#########################################################


all:
	: # do nothing


rice-embed-go:
	@echo "### Rice: Embedding static assets"
	cd cmd/blitzinfod/ && $(RICE_BIN) embed-go


build-all: build-amd64 build-armhf
	@echo "### Build: all - done"


# ARCH "amd64"
build-amd64: rice-embed-go
	@echo "### Build: amd64"
	GOARCH=amd64 $(GOBUILD) -o build/amd64/blitzinfod cmd/blitzinfod/main.go cmd/blitzinfod/rice-box.go
	GOARCH=amd64 $(GOBUILD) -o build/amd64/blitzinfo-cli cmd/blitzinfo-cli/main.go


# ARCH "armhf"
build-armhf: rice-embed-go
	@echo "### Build: armhf"
	GOARCH=arm GOARM=7 $(GOBUILD) -o build/armhf/blitzinfod cmd/blitzinfod/main.go cmd/blitzinfod/rice-box.go
	GOARCH=arm GOARM=7 $(GOBUILD) -o build/armhf/blitzinfo-cli cmd/blitzinfo-cli/main.go


install:
	@echo "### Install"
	install -d -m 755                     $(DESTDIR)$(PREFIX)/bin
	install -m 0755 blitzinfod-$(DEB_HOST_ARCH)    $(DESTDIR)$(PREFIX)/bin
	# install -m 0755 blitzinfod-armv7l    $(DESTDIR)$(PREFIX)/bin


uninstall:
	@echo "### Uninstall"
	-rm -f $(DESTDIR)$(PREFIX)/bin/blitzinfod-$(DEB_HOST_ARCH)


clean:
	@echo "### Clean"
	-rm -rf debian/.debhelper/
	-rm -rf debian/blitzinfod/
	-rm -f debian/debhelper-build-stamp


distclean:
	@echo "### Distclean"
	-rm -rf build/
	-rm -rf debian/.debhelper/
	-rm -rf debian/blitzinfod/
	-rm -f debian/debhelper-build-stamp


list:
	@$(do-listing)


package-armhf:
	@echo 'packaging blitzinfod armhf for XXX Debian unstable...'
	mkdir -p $(PD)/unstable
	rm -f $(PD)/unstable/*
	sed -i 's,\(^$(P) ('$(VERSION)') \)[a-z]*,\1unstable,' debian/changelog
	debuild -i -I -I.git -Ipackage -a armhf
	@echo ""
	cd .. && for f in $(PACKAGE_FILES); do \
		[ -f $$f ] && mv $$f $(D)/$(PD)/unstable || :; done


package-amd64:
	@echo 'packaging blitzinfod amd64 for XXX Debian unstable...'
	mkdir -p $(PD)/unstable
	rm -f $(PD)/unstable/*
	sed -i 's,\(^$(P) ('$(VERSION)') \)[a-z]*,\1unstable,' debian/changelog
	debuild -i -I -I.git -Ipackage -a amd64
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


# phonies
.PHONY: install uninstall build clean list package-armhf package-amd64
.PHONY: package-info package-purge package-clean package-version
