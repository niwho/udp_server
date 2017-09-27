all:
	mkdir -p output/bin output/conf
	cp -r conf/* output/conf/
	cp script/bootstrap.sh script/settings.py output
	go build -o output/bin/pusic_monitor_agent

clean:
	rm -rf output

