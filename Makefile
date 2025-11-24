run:
	go build -o build/worker
	./build/worker

clean:
	sudo rm -rf database/database.db tmp/*
	rm build/worker