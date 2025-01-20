.PHONY: all clean


all:
	rm -vf local.db*
	env PATH=..:$(PATH) ANYDB_DB=local.db vhs demo.tape

clean:
	rm -vf *.db* *.gif *.json
