.PHONY: all clean check

VHS  = vhs
TAPE = demo.tape

all: clean check
	env PATH=..:$(PATH) ANYDB_DB=local.db $(VHS) $(TAPE)

clean:
	rm -vf *.db* *.gif *.json

check:
	ls -l ../anydb
