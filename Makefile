up-app:
	docker compose up --build -d 
down-app:
	docker compose down -v
check-container:
	docker ps -a