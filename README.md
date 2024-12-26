# AttandanceSystemGolang

Attendance Management System
│
├── backend
│ ├── auth
│ ├── controller
│ ├── db
│ ├── .env
│ ├── .gitignore
│ ├── go.mod
│ ├── go.sum
│ ├── main.go
│ └── Dockerfile
│
└── frontend
| ├── assets
| ├── dashboard
| ├── indexedDB.html
| └── Dockerfile
│
├── README.md
└── docker-compose.yml

Backend
sudo docker run -p 8080:8080 --env DATABASE_URL="mongodb://$(sudo docker inspect -f '{{range.NetworkSettings.Networks}}{{.IPAddress}}{{end}}' mongodb-container):27017" attendance

Frontend
docker run -p 80:80 frontend-nginx
