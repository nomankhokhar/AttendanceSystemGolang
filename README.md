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

Run Backend Docker Image When Locally Build Image

```bash
sudo docker run -p 8080:8080 --env DATABASE_URL="mongodb://$(sudo docker inspect -f '{{range.NetworkSettings.Networks}}{{.IPAddress}}{{end}}' mongodb-container):27017" attendance

```

Pulled Image From Docker Hub after building

```bash
sudo docker run -p 8080:8080 --env DATABASE_URL="mongodb://$(sudo docker inspect -f '{{range.NetworkSettings.Networks}}{{.IPAddress}}{{end}}' mongodb-container):27017" nomanali1114/backend
```

Run Frontend Docker Image When Locally Build Image

```bash
docker run -p 80:80 frontend-nginx
```

Pulled Image From Docker Hub after building

```bash
docker run -p 80:80 nomanali1114/frontend
```

build by Noman Ali
