SUBDIRS := telephono telephono-ui launchpad

.PHONY: all $(SUBDIRS)

all: $(SUBDIRS)

clean:
	@for d in $(SUBDIRS); do $(MAKE) -C $$d clean; done

telephono:
	$(MAKE) -C telephono all

telephono-ui:
	$(MAKE) -C telephono-ui all

launchpad:
	$(MAKE) -C launchpad all
