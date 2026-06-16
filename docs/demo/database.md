# PostgreSQL — หลักฐานว่า schema และ query ทำงาน

ฐานข้อมูลในโปรเจกต์นี้รันใน `docker-compose` (container `myblogs-postgres`)
ตาราง `blogs` และ `blog_users` ถูกสร้างอัตโนมัติด้วย GORM `AutoMigrate`
(ดู [`blogservice/modules/blogs/repositories/blog_repositories.go`](../../blogservice/modules/blogs/repositories/blog_repositories.go))

> demo ทั้งหมดนี้แคปจาก terminal ของ PostgreSQL ใน Docker โดยตรง

## วิธีพิสูจน์

ต้องรัน `docker-compose up -d` + blogservice (เพื่อให้ AutoMigrate สร้างตาราง) ก่อน

เข้า psql ใน container

```powershell
docker exec -it myblogs-postgres psql -U postgres -d myblogs
```

### 1. ดูตารางที่ถูกสร้าง

```sql
\dt
```

ต้องเห็นตาราง `blogs` และ `blog_users`

📸 _screenshot:_ `images/db-tables.png`

### 2. ดูโครงสร้างคอลัมน์ของตาราง blogs

```sql
\d blogs
```

เห็นคอลัมน์ที่ map มาจาก struct `Blog` ใน Go (title, blog_desc, content, cover_image, user_id ...)

📸 _screenshot:_ `images/db-describe-blogs.png`

### 3. ผล query ข้อมูลบล็อก

```sql
SELECT id, title, user_name FROM blogs;
```

📸 _screenshot:_ `images/db-select-blogs.png`

### 4. ยืนยัน user ถูก sync ข้ามมาจาก userservice (ผ่าน Kafka)

```sql
SELECT id, name FROM blog_users;
```

แถวที่เห็นมาจากการสมัครผ่าน userservice แล้ว event วิ่งผ่าน Kafka มาลง
(เชื่อมกับ [demo Kafka](kafka.md))

📸 _screenshot:_ `images/db-blog-users.png`

## สิ่งที่หลักฐานนี้บอก

- เข้าใจ schema migration ด้วย ORM (GORM AutoMigrate) — struct ใน Go → ตารางจริง
- อ่าน/ตรวจข้อมูลด้วย SQL และ meta-command ของ psql เป็น (มีพื้นจากการเป็น TA วิชา Database Systems)
- เห็นภาพรวมว่าข้อมูลจากคนละ service มาอยู่ใน DB ได้อย่างไร

---

> เมื่อ deploy ขึ้น **Supabase** ภายหลัง หลักฐานชุดเดียวกันนี้ดูได้จาก Table Editor / SQL Editor
> (เปลี่ยน `DB_HOST`/`DB_PASSWORD` ชี้ไป Supabase และตั้ง `DB_SSLMOD=require`) — ค่อยเพิ่ม screenshot ตอนนั้น
