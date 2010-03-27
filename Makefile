include $(GOROOT)/src/Make.$(GOARCH)
 
TARG=properties
GOFILES=properties.go
 
CLEANFILES+=$(TARG)_test
 
include $(GOROOT)/src/Make.pkg
