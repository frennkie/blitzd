# Makefile


all:
	@echo "### VFS: Embedding static assets"
	go generate web/web.go


help:
	@echo 'Makefile for blitzd Debian Packaging                          '
	@echo '                                                              '
	@echo 'Usage:                                                        '
	@echo '   make all                       complete packaging process  '

# phonies
.PHONY: all
