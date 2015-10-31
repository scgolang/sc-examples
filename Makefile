SRCS  = $(wildcard *.go)
PROGS = $(subst .go,,$(SRCS))

all:
	for f in $(SRCS); do \
		go build $$f & \
	done; wait

clean:
	rm -f $(PROGS)

.PHONY: all clean
