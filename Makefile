bin/cobra:
	cd vendor/github.com/spf13/cobra/cobra \
	  && go build -o $(PWD)/$@ ./
