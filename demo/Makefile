.PHONY: all clean check


all: clean check
	env PATH=..:$(PATH) ANYDB_DB=local.db vhs demo.tape

clean:
	rm -vf *.db* *.gif *.json

check:
	ls -l ../anydb
