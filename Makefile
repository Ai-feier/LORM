e2e:
	docker compose down 
	docker compose up -d 
	go test -tags=e2e ./... 
	docker compose down

utiltest:
	go test -v ./...