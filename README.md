# Myblogs — Microservices Blog Platform

A blog platform built with **Go microservices** (Fiber, Clean Architecture) communicating through **Kafka**, with a **React** frontend. Uses **PostgreSQL**, **Redis** caching, and is deployable with **Docker Compose** or **Kubernetes**.

![demo](image.png)

## Architecture

```
React (Vite) ──HTTP──► userservice ──Kafka event──► blogservice
                          │           (user.created)     │
                          ▼                              ▼
                      PostgreSQL  ◄── Redis (cache) ──►  PostgreSQL
```

- **userservice** — Manage user and produce event into Kafka
- **blogservice** — Manage blog and consume event from userservice

## Tech Stack

| Layer | Tech |
|---|---|
| Frontend | React 18 (Vite), TailwindCSS |
| Backend | Go 1.20 (Fiber, GORM) |
| Database | PostgreSQL 16, Redis |
| Message Queue | Kafka |
| Infra | Docker Compose, Kubernetes |

## Run locally

Require - Docker Desktop, Go, Node.js

```bash
# 1. run infra (postgres, redis, kafka, zookeeper)
docker-compose up -d

# 2. run backend 
cd userservice/app && go run main.go     
cd blogservice/app && go run main.go     

# 3. run frontend
cd my-blog-react && npm install && npm run dev   
```

## ผลลัพธ์การเชื่อมต่อ Tools และ Infras

_เร็ว ๆ นี้ — กำลังเตรียม screenshot การทำงานของแต่ละ tool_

## API (summary)

| Method | Endpoint | Service | หน้าที่ |
|---|---|---|---|
| POST | `/v1/myblogs/register` | userservice :8001 | สมัคร user |
| GET | `/v1/myblogs/getUsers` | userservice :8001 | ดู user ทั้งหมด |
| POST | `/v1/myblogs/createBlog/:userId` | blogservice :8002 | สร้างบล็อก |
| GET | `/v1/myblogs/getBlogs` | blogservice :8002 | ดึงบล็อกทั้งหมด |
