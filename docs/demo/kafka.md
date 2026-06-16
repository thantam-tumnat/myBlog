# Kafka (message queue) — หลักฐานว่าใช้งานจริง

Kafka ทำให้สอง service คุยกันแบบ async โดยไม่เรียกกันตรง ๆ

- `userservice` สมัคร user เสร็จ → **produce** event `myblogs.user.created`
- `blogservice` **consume** event นั้น → เก็บ user ลงตาราง `blog_users`

โค้ดฝั่งรับอยู่ที่ [`blogservice/modules/consumer/usecases/blogConsume.go`](../../blogservice/modules/consumer/usecases/blogConsume.go)
นิยาม event อยู่ที่ [`blogservice/modules/entities/event.go`](../../blogservice/modules/entities/event.go)

## วิธีพิสูจน์ว่า produce/consume ทำงาน

ต้องรัน `docker-compose up -d` + userservice + blogservice พร้อมกัน

### 1. สมัคร user แล้วดู event ไหลข้าม service

```powershell
curl -X POST http://localhost:8001/v1/myblogs/register `
  -H "Content-Type: application/json" `
  -d '{\"username\":\"tum\",\"password\":\"1234\",\"name\":\"Tum\",\"description\":\"Writer\",\"userImage\":\"https://i.pravatar.cc/150\"}'
```

จากนั้นดู log ของ **blogservice** จะขึ้น

```
Received Topic: myblogs.user.created
Received message: {"userId":1,"name":"Tum",...}
```

> หลักฐานว่า event ที่ userservice ส่ง วิ่งผ่าน Kafka มาถึง blogservice จริง

📸 _screenshot:_ `images/kafka-consume-log.png` (วาง log สอง service คู่กันให้เห็น produce → consume)

### 2. ยืนยันว่า user ถูก sync ข้ามมาที่ DB ของ blog

```powershell
docker exec -it myblogs-postgres psql -U postgres -d myblogs -c "SELECT id, name FROM blog_users;"
```

จะเห็นแถว user ที่เพิ่งสมัคร โผล่ในตาราง `blog_users` ของ blogservice
(ทั้งที่เราสมัครผ่าน userservice) — นี่คือผลของ event-driven sync

📸 _screenshot:_ `images/kafka-blog-users-row.png`

### 3. (เสริม) ส่อง message ใน topic ด้วย kafka console consumer

```powershell
docker exec -it myblogs-kafka kafka-console-consumer `
  --bootstrap-server localhost:9092 `
  --topic myblogs.user.created --from-beginning
```

จะเห็น raw message ที่ถูก produce เข้า topic โดยตรง

📸 _screenshot:_ `images/kafka-console-consumer.png`

## สิ่งที่หลักฐานนี้บอก

- เข้าใจ **producer / consumer / topic** และ event-driven communication
- เข้าใจว่าทำไมต้องใช้ async (loose coupling — userservice ไม่ต้องรอ blogservice)
- รู้วิธี debug Kafka ผ่าน console consumer
