run:
	go build -o build/worker
	sudo ./build/worker

checkRace:
	go build -o build/worker -race

clean:
	sudo rm -rf database/database.db tmp/*
	rm build/worker

vet:
	go vet ./...