.PHONY: demo clean check clean-demo

VHS  = vhs



clean-demo:
	rm -f local.db*

%.gif: %.tape
	@echo "vhs $<"
	env PATH=..:$(PATH) ANYDB_DB=local.db vhs $<

clean:
	rm -vf *.db* *.json

check:
	ls -l ../anydb

demo: check clean-demo intro.gif advanced.gif

