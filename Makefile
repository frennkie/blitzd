# Makefile for blitzd
# Heavily inspired by https://github.com/hamishcunningham/wiringpi
# Beware: This uses tab indent (instead of spaces)

# vars #######################################################################
VERSION := $(shell grep -e "^blitzd " debian/changelog | head -1 | cut -d "(" -f2 | cut -d ")" -f1)
DATE := $(shell date --iso-8601=seconds)
GIT_VERSION := $(shell git describe --abbrev=40 --dirty)
P=blitzd
D=blitzd
PD=package/$(VERSION)

PKG=github.com/frennkie/blitzd
BUILDFLAGS="-X $(PKG)/pkg/cmd/blitzd.BuildTime=$(DATE) \
  -X $(PKG)/pkg/cmd/blitzd.BuildVersion=$(VERSION) \
  -X $(PKG)/pkg/cmd/blitzd.BuildGitVersion=$(GIT_VERSION)"

PREFIX=/usr/local

# Target/Destination Host
DEB_HOST_ARCH  := $(shell dpkg-architecture -qDEB_HOST_ARCH)
# Building Host
DEB_BUILD_ARCH := $(shell dpkg-architecture -qDEB_BUILD_ARCH)


INSTALLED_FILES=\
  $(DESTDIR)$(PREFIX)/bin/blitzd \
  $(DESTDIR)$(PREFIX)/bin/blitz-cli \
  $(DESTDIR)/etc/blitzd.toml \
  $(DESTDIR)/lib/systemd/system/blitzd.service


PACKAGE_FILES=\
  $(P)_$(VERSION)*.deb \
  $(P)_$(VERSION)*.dsc \
  $(P)_$(VERSION)*.tar.xz \
  $(P)_$(VERSION)*.build \
  $(P)_$(VERSION)*.buildinfo \
  $(P)_$(VERSION)*.changes


GOBUILD := go build -v


# functions ##################################################################
do-listing=\
  echo; echo 'Listing installed files:'; ls -l --color $(INSTALLED_FILES) || :


#########################################################
# by default `make` will build first target in Makefile #
#########################################################


nothing:
	: # do nothing


all: package dist fix-sign clean


vfsgen:
	@echo "### VFS: Embedding static assets"
	go generate web/web.go


build: build-amd64 build-armhf
	@echo "### Build: all - done"


# ARCH "amd64"
build-amd64:
	@echo "### Build: amd64"
	@echo "GitVersion: $(GIT_VERSION)"
	GOARCH=amd64 $(GOBUILD) -ldflags $(BUILDFLAGS) -o build/amd64/blitzd cmd/blitzd/main.go
	GOARCH=amd64 $(GOBUILD) -ldflags $(BUILDFLAGS) -o build/amd64/blitz-cli cmd/cli/main.go


# ARCH "armhf"
build-armhf:
	@echo "### Build: armhf"
	GOARCH=arm GOARM=7 $(GOBUILD) -ldflags $(BUILDFLAGS) -o build/armhf/blitzd cmd/blitzd/main.go
	GOARCH=arm GOARM=7 $(GOBUILD) -ldflags $(BUILDFLAGS) -o build/armhf/blitz-cli cmd/cli/main.go


dist:
	@echo "### Dist"
	cd build/amd64/ && tar czf ../../package/$(VERSION)/unstable/blitzd_$(VERSION)_amd64.tgz blitzd blitz-cli
	cd build/armhf/ && tar czf ../../package/$(VERSION)/unstable/blitzd_$(VERSION)_armhf.tgz blitzd blitz-cli


fix-sign:
	@echo "### Fix Debian Signature for armhf"
	cd package/$(VERSION)/unstable && debsign blitzd_$(VERSION)_armhf.changes -a armhf


install:
	@echo "### Install: Arch=$(DEB_HOST_ARCH)"
	install -d -m 755 $(DESTDIR)$(PREFIX)/bin
	install -d -m 755 $(DESTDIR)/etc
	install -m 0755 build/$(DEB_HOST_ARCH)/blitzd 	$(DESTDIR)$(PREFIX)/bin
	install -m 0755 build/$(DEB_HOST_ARCH)/blitz-cli 	$(DESTDIR)$(PREFIX)/bin
	install -m 0755 configs/blitzd.toml 	$(DESTDIR)/etc


uninstall:
	@echo "### Uninstall"
	-rm -f $(DESTDIR)$(PREFIX)/bin/blitzd
	-rm -f $(DESTDIR)$(PREFIX)/bin/blitz-cli
	-rm -f $(DESTDIR)/etc/blitzd.toml


clean: distclean
	@echo "### Clean"
	-rm -rf build/


distclean:
	@echo "### Distclean (this is called by debuild)"
	-rm -rf debian/.debhelper/
	-rm -f debian/blitzd.substvars
	-rm -rf debian/blitzd/
	-rm -f debian/debhelper-build-stamp
	-rm -f debian/files


list:
	@$(do-listing)


package: package-amd64 package-armhf package-move
	@echo 'Packaging all - done'


package-amd64: build-amd64
	@echo 'Package: amd64'
	mkdir -p $(PD)/unstable
	rm -f $(PD)/unstable/*
	sed -i 's,\(^$(P) ('$(VERSION)') \)[a-z]*,\1unstable,' debian/changelog
	debuild -i -I -I.git -Ipackage -a amd64
	@echo ""


package-armhf: build-armhf
	@echo 'Package: armhf'
	mkdir -p $(PD)/unstable
	rm -f $(PD)/unstable/*
	sed -i 's,\(^$(P) ('$(VERSION)') \)[a-z]*,\1unstable,' debian/changelog
	debuild -uc -us -i -I -I.git -Ipackage -a armhf
	@echo ""


package-move:
	@echo 'Move: Packages from parent to sub directory'
	cd .. && for f in $(PACKAGE_FILES); do \
		[ -f $$f ] && mv $$f $(D)/$(PD)/unstable || :; done


package-info:
	apt-cache showpkg $(P)
package-purge:
	apt-get purge $(P)
package-clean:
	dh_clean
package-version:
	debchange -i -D unstable


goget:
	@echo 'Install: Missing Go libs'
	go get ./...

help:
	@echo 'Makefile for infoblitz Debian Packaging                       '
	@echo '                                                              '
	@echo 'Usage:                                                        '
	@echo '   make all                       complete packaging process  '
	@echo '   make install                   installs infoblitz          '
	@echo '   make uninstall                 removes infoblitz           '
	@echo '   make build                     build                       '
	@echo '   make clean                     clean                       '
	@echo '   make list                      list installed files        '
	@echo '   make build                     build infoblitz from go src '
	@echo '   make build-amd64               build infoblitz from go src '
	@echo '   make build-armhf               build infoblitz from go src '
	@echo '   make package                   create .deb/.dcs packagings '
	@echo '   make package-amd64             create .deb/.dcs packagings '
	@echo '   make package-armhf             create .deb/.dcs packagings '
	@echo '   make package-clean             delete tmp packaging files  '
	@echo '   make package-info              apt-cache showpkg           '
	@echo '   make package-purge             apt-get purge               '
	@echo '   make package-version           update the changelog        '
	@echo '   make goget                     install missing go libs     '
	@echo '                                                              '


# phonies
.PHONY: all install uninstall build clean list package-armhf package-amd64
.PHONY: package-move package-info package-purge package-clean package-version
