include $(GOROOT)/src/Make.inc
 
TARG=properties
GOFILES=properties.go
 
CLEANFILES+=$(TARG)_test
 
include $(GOROOT)/src/Make.pkg
