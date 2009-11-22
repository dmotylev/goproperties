# Copyright 2009 The Go Authors.  All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.
 
include $(GOROOT)/src/Make.$(GOARCH)
 
TARG=properties
#CGOFILES=$(TARG).go
CGOFILES=
GOFILES=properties.go
CGO_LDFLAGS=
 
CLEANFILES+=$(TARG)_test
 
include $(GOROOT)/src/Make.pkg
 
%: install %.go
	$(GC) $*.go
	$(LD) -o $@ $*.$O