# Redis (caching) — หลักฐานว่าใช้งานจริง

Redis ในโปรเจกต์นี้ทำหน้าที่ **แคชผลลัพธ์ของ `GetBlogs()`** เพื่อลดภาระ PostgreSQL
โค้ดอยู่ที่ [`blogservice/modules/blogs/repositories/blog_repositories.go`](../../blogservice/modules/blogs/repositories/blog_repositories.go)

แนวคิด: query ครั้งแรกดึงจาก Postgres แล้วเก็บลง Redis (TTL 10 วินาที) — ครั้งถัด ๆ ไปภายใน 10 วิ จะอ่านจาก Redis แทน

## วิธีพิสูจน์ว่า cache ทำงาน

ต้องรันระบบด้วย `docker-compose up -d` + blogservice ก่อน

### 1. ดู log ตอน cache hit

เรียก endpoint เดิมสองครั้งติดกัน

```powershell
curl http://localhost:8002/v1/myblogs/getBlogs
curl http://localhost:8002/v1/myblogs/getBlogs
```

ครั้งที่สอง log ของ blogservice จะขึ้นบรรทัด

```
blogs Retrieved From Redis
```

> นี่คือหลักฐานว่ารอบสองไม่ได้แตะ Postgres แต่อ่านจาก Redis

📸 _screenshot:_ `images/redis-cache-hit-log.png`

### 2. ส่องข้อมูลที่ถูกแคชใน Redis โดยตรง

```powershell
docker exec -it myblogs-redis redis-cli
```

ในนั้นพิมพ์

```
KEYS blogrepository*
GET blogrepository::Getblogs
TTL blogrepository::Getblogs
```

จะเห็น key ที่เก็บ JSON ของบล็อก และ TTL ที่นับถอยหลังจาก 10 → 0

📸 _screenshot:_ `images/redis-cli-keys-ttl.png`

### 3. (เสริม) ดู cache ถูกล้างเมื่อมีการสร้างบล็อกใหม่

ในโค้ด ตอน `CreateBlog` / `CreateUser` จะลบ key `blogrepository*` ทิ้ง
ลองสร้างบล็อกใหม่แล้วเช็คว่า key หายไป (cache invalidation)

📸 _screenshot:_ `images/redis-cache-invalidation.png`

## สิ่งที่หลักฐานนี้บอก

- เข้าใจ pattern **cache-aside** (อ่าน cache ก่อน ไม่เจอค่อยลง DB แล้ว set cache)
- เข้าใจ **TTL** และ **cache invalidation** เมื่อข้อมูลเปลี่ยน
- รู้วิธี debug Redis ผ่าน `redis-cli`
