# Kafka + Redis + Kubernetes: ใช้งานร่วมกันใน Microservices จริง ๆ อย่างไร

> บทความนี้เล่าผ่านโปรเจกต์จริง **Myblogs Microservices** — ระบบบล็อกที่แยกเป็น `userservice` และ `blogservice` (เขียนด้วย Go + Fiber + GORM), คุยกันผ่าน **Kafka**, แคชด้วย **Redis**, และ deploy ได้ทั้งบน Docker Compose และ **Kubernetes**

---

## 1. ทำไมต้องใช้ทั้งสามตัวพร้อมกัน

ในระบบ microservices เรามักเจอ "ปัญหาคลาสสิก" สามอย่าง และแต่ละตัวก็แก้คนละเรื่อง:

| ปัญหา | เครื่องมือที่ใช้แก้ | บทบาทในโปรเจกต์นี้ |
|---|---|---|
| service คุยกันแบบไม่ผูกติด (decoupling) และไม่ทำให้กันล่ม | **Kafka** | `userservice` สร้าง user → ส่ง event → `blogservice` รับไป sync ตาราง `blog_users` |
| อ่านข้อมูลซ้ำ ๆ จาก DB ช้า/แพง | **Redis** | แคชผลลัพธ์ `GetBlogs()` ไว้ชั่วคราว ลดภาระ Postgres |
| รันหลาย service + stateful component ให้ scale และ self-heal | **Kubernetes** | แต่ละ service เป็น Deployment, Kafka/Zookeeper เป็น StatefulSet |

หัวใจคือ **แต่ละตัวไม่ได้แข่งกัน แต่อยู่คนละชั้นของระบบ**:

```
            ┌─────────────────────────────────────────────┐
            │                Kubernetes                    │
            │  (orchestration: scaling, healing, service)  │
            │                                              │
   client ──┤  ┌────────────┐   Kafka topic   ┌──────────┐ │
            │  │ userservice │ ──────────────▶ │blogservice│ │
            │  └─────┬──────┘  (user.created)  └────┬─────┘ │
            │        │                              │       │
            │     Postgres                      Postgres    │
            │        │                              │       │
            │     Redis (cache)                 Redis(cache) │
            │                                              │
            └─────────────────────────────────────────────┘
```

- **Kafka** = ชั้น "การสื่อสารแบบ async" (event-driven)
- **Redis** = ชั้น "เร่งความเร็วการอ่าน" (caching)
- **Kubernetes** = ชั้น "รันและดูแลทุกอย่าง" (orchestration)

---

## 2. Kafka: การสื่อสารแบบ event-driven

### แนวคิด

แทนที่ `userservice` จะเรียก HTTP ตรงไปหา `blogservice` ทุกครั้งที่มี user ใหม่ (ซึ่งถ้า blogservice ล่ม การสร้าง user ก็พังตาม) เราใช้ Kafka เป็น "ตัวกลางรับฝากข้อความ":

1. `userservice` (Producer) สร้าง user เสร็จ → ยิง event ลง topic
2. Kafka เก็บ event ไว้อย่างทนทาน (persistent log)
3. `blogservice` (Consumer) ค่อย ๆ ดึง event ไป process ตามจังหวะของตัวเอง

ผลคือ **สอง service ไม่ผูกชะตากัน** — blogservice ล่มชั่วคราวก็ไม่กระทบการสมัครสมาชิก

### Producer ฝั่ง userservice

จากโค้ดจริง `userservice/modules/producer/handler/producer.go`:

```go
func (obj eventProducer) Produce(event entities.Event) error {
    // ถ้าไม่มี broker (Kafka ล่ม) ก็ข้ามการ publish แทนที่จะ panic
    // เพื่อให้การสร้าง user ยังสำเร็จอยู่
    if obj.producer == nil {
        logs.Info("kafka producer not available, skipping event publish")
        return nil
    }

    topic := event.TopicName()
    value, err := json.Marshal(event)
    if err != nil {
        return err
    }

    msg := sarama.ProducerMessage{
        Topic: topic,
        Value: sarama.ByteEncoder(value),
    }
    _, _, err = obj.producer.SendMessage(&msg)
    return err
}
```

> 💡 **จุดที่ออกแบบมาดี:** เช็ค `obj.producer == nil` ก่อน ถ้า Kafka ใช้ไม่ได้ก็ยอมข้าม event ไป ไม่ปล่อยให้ flow หลัก (สร้าง user) ล้มตาม — นี่คือหลัก *graceful degradation*

### Consumer ฝั่ง blogservice

`blogservice/pkg/databases/kafka/kafka.go` สร้าง consumer group:

```go
consumer, err := sarama.NewConsumerGroup(
    []string{serversUrl1}, "myblogs.blogservice", saramaConfig,
)
```

และ `blogservice/modules/consumer/handler/consumer.go` วน loop รับ message:

```go
func (obj *consumerHandler) ConsumeClaim(
    session sarama.ConsumerGroupSession,
    claim sarama.ConsumerGroupClaim,
) error {
    for message := range claim.Messages() {
        err := obj.eventHandler.Handle(message.Topic, message.Value)
        if err != nil {
            logs.Error(fmt.Sprintf("Error handling message: %v", err))
            return err
        }
        session.MarkMessage(message, "")   // commit offset เมื่อ process สำเร็จ
    }
    return nil
}
```

> ⚠️ **ข้อควรระวังเรื่อง offset:** การ `MarkMessage` ทำหลัง `Handle` สำเร็จ = at-least-once delivery. แปลว่า event อาจถูก process ซ้ำได้ (เช่น crash หลัง process แต่ก่อน commit) → consumer ควรเขียนให้ **idempotent** (เช่น upsert แทน insert ตรง ๆ)

### Concept คำศัพท์ที่ต้องรู้

- **Topic** — ช่องทาง/ประเภทของ event (เช่น `user.created`)
- **Partition** — แบ่ง topic เป็นหลายส่วนเพื่อ parallel + scale
- **Consumer Group** — กลุ่ม consumer ที่แบ่ง partition กันกิน; เพิ่ม pod ในกลุ่มเดียวกัน = scale การ consume
- **Offset** — ตำแหน่งที่อ่านถึงแล้วในแต่ละ partition

---

## 3. Redis: แคชเพื่อลดภาระ Database

### Pattern: Cache-Aside (Lazy Loading)

จาก `blogservice/modules/blogs/repositories/blog_repositories.go` ฟังก์ชัน `GetBlogs()`:

```go
key := "blogrepository::Getblogs"

// 1. ลองอ่านจาก Redis ก่อน
productsJson, err := r.client.Get(context.Background(), key).Result()
if err == nil {
    if json.Unmarshal([]byte(productsJson), &blogs) == nil {
        logs.Info("blogs Retrieved From Redis")
        return blogs, nil      // cache hit → คืนเลย ไม่แตะ DB
    }
}

// 2. cache miss → query DB แล้วเก็บลง Redis (TTL 10 วินาที)
data, _ := json.Marshal(blogs)
if err = r.client.Set(context.Background(), key, string(data),
    time.Second*10).Err(); err != nil {
    logs.Error(err)   // แคชล้มเหลวก็ยังคืนข้อมูลจาก DB ได้
}
```

### Cache Invalidation เมื่อมีการเขียน

เมื่อ `CreateBlog()` หรือ `CreateUser()` มีการเขียนข้อมูลใหม่ ต้องล้างแคชเก่าทิ้ง:

```go
keys, err := r.client.Keys(context.Background(), "blogrepository*").Result()
for _, key := range keys {
    r.client.Del(context.Background(), key)
}
```

> 💡 **จุดออกแบบที่ดี:** ทุกการเรียก Redis ในโปรเจกต์นี้เป็น *best-effort* — ถ้า Redis ล่ม จะ log error แล้วไปต่อด้วยข้อมูลจาก Postgres ระบบไม่พังเพราะแคชล่ม
>
> ⚠️ **จุดที่ปรับปรุงได้:** `KEYS blogrepository*` เป็นคำสั่ง O(N) ที่ block Redis ทั้งตัวเมื่อ key เยอะ — production ควรเปลี่ยนไปใช้ `SCAN` หรือออกแบบ key ให้ลบแบบเจาะจงได้

### Kafka คู่กับ Redis: combo ที่ลงตัว

ทั้งสองทำงานเสริมกันได้สวย — **เมื่อ consumer รับ event ที่เปลี่ยนข้อมูล ก็สั่ง invalidate cache ไปด้วย**:

```
userservice เขียน user → Kafka event → blogservice consume
   → เขียน blog_users ลง DB → ล้าง Redis key "blogrepository*"
   → ครั้งหน้าที่อ่าน จะได้ข้อมูลสด
```

นี่คือวิธีรักษา cache ให้ตรงกับ DB ในระบบ distributed โดยไม่ต้องให้ทุก service รู้จักกันตรง ๆ

---

## 4. การ Deploy แบบต่าง ๆ

โปรเจกต์นี้รองรับหลายวิธี deploy ไล่จากง่าย → production:

### 4.1 Docker Compose (Local / Dev)

`docker-compose.yml` รันทุกอย่างในเครื่องเดียวด้วยคำสั่งเดียว — เหมาะกับ dev:

```yaml
services:
  postgres:   { image: postgres:16, ... }
  redis:      { image: redis:7, ports: ["6379:6379"] }
  zookeeper:  { image: confluentinc/cp-zookeeper:7.5.0, ... }
  kafka:
    image: confluentinc/cp-kafka:7.5.0
    depends_on: [zookeeper]
    environment:
      KAFKA_ZOOKEEPER_CONNECT: zookeeper:2181
      KAFKA_AUTO_CREATE_TOPICS_ENABLE: "true"
```

```bash
docker compose up -d
```

**ข้อดี:** ง่าย, เร็ว, เหมาะ dev/CI
**ข้อจำกัด:** ไม่มี auto-scaling, ไม่มี self-healing, เครื่องเดียว = single point of failure

### 4.2 Kubernetes — Stateless services ด้วย Deployment

service ที่ไม่มี state (เช่น `blogservice`) ใช้ **Deployment** ได้เลย เพราะแต่ละ pod เหมือนกันหมด สลับได้อิสระ

`blogservice/deploy.yaml`:

```yaml
kind: Deployment
metadata: { name: blog }
spec:
  selector: { matchLabels: { app: blog } }
  template:
    spec:
      containers:
      - name: blog
        image: ghcr.io/thantam-tumnat/blogservice:latest
        env:
        - { name: KAFKA_SERVERS, value: "kafka" }   # ชี้ไป Service ชื่อ kafka
        - { name: REDIS_HOST,    value: "redis" }   # ชี้ไป Service ชื่อ redis
        - { name: DB_HOST,       value: "postgres" }
        - name: DB_PASSWORD_BLOG
          valueFrom:
            secretKeyRef: { name: postgres-secret, key: password }   # ใช้ Secret
  strategy:
    type: RollingUpdate          # อัปเดตทีละ pod ไม่ให้ downtime
    rollingUpdate: { maxUnavailable: 1, maxSurge: 1 }
```

ประเด็นสำคัญ:
- **Service discovery แบบ DNS** — แค่ใส่ `KAFKA_SERVERS=kafka` pod ก็หา Kafka เจอ เพราะ K8s สร้าง DNS ให้ทุก Service อัตโนมัติ
- **Secret** สำหรับรหัสผ่าน DB แทนการ hardcode
- **RollingUpdate** ทำให้ deploy เวอร์ชันใหม่โดยไม่ล่ม

### 4.3 Kubernetes — Stateful components ด้วย StatefulSet

Kafka, Zookeeper, Postgres มี **state** (ข้อมูลบนดิสก์ + identity คงที่) จึงใช้ **StatefulSet** ไม่ใช่ Deployment

`blogservice/k8s-files/kafka/kafka-stateful-set.yaml`:

```yaml
kind: StatefulSet
metadata: { name: kafka }
spec:
  serviceName: kafka
  replicas: 1
  template:
    spec:
      containers:
      - name: kafka
        image: bitnami/kafka:latest
        env:
        - { name: KAFKA_ZOOKEEPER_CONNECT, value: "zookeeper:2181" }
        volumeMounts:
        - { mountPath: /var/lib/kafka/, name: kafka-data }
  volumeClaimTemplates:           # ★ แต่ละ pod ได้ disk ของตัวเอง
  - metadata: { name: kafka-data }
    spec:
      accessModes: [ReadWriteOnce]
      resources: { requests: { storage: 1Gi } }
```

**ทำไม StatefulSet ไม่ใช่ Deployment?**

| | Deployment | StatefulSet |
|---|---|---|
| ชื่อ pod | สุ่ม (`blog-7d9f...`) | คงที่ มีลำดับ (`kafka-0`, `kafka-1`) |
| storage | แชร์/ไม่มี | แต่ละ pod มี PVC ของตัวเอง |
| การ start/stop | พร้อมกัน | ตามลำดับ |
| เหมาะกับ | service ไม่มี state | DB, Kafka, Zookeeper |

> Redis ในโปรเจกต์นี้ใช้ **Deployment** (`redis-pod.yaml`) เพราะใช้เป็นแค่แคช — ข้อมูลหายได้ ไม่ต้องการ persistence เข้มงวด ถ้าจะใช้ Redis เก็บข้อมูลถาวรหรือทำ cluster ค่อยเปลี่ยนเป็น StatefulSet

### 4.4 เปิดให้เข้าถึงจากภายนอกด้วย Ingress

`blogservice/deploy.yaml` มี Service (ภายใน cluster) + Ingress (เข้าจากภายนอก):

```yaml
kind: Service          # ภายใน cluster: pod อื่นเรียก http://blog:8002
spec:
  ports: [{ port: 8002, targetPort: 8002 }]
  selector: { app: blog }
---
kind: Ingress          # ภายนอก: route ตาม path
spec:
  rules:
  - host: localhost
    http:
      paths:
      - { path: /v1/myblogs/getBlogs,   backend: { service: { name: blog, port: { number: 8002 }}}}
      - { path: /v1/myblogs/createBlog, backend: { service: { name: blog, port: { number: 8002 }}}}
```

ลำดับชั้นการเข้าถึง: **Ingress (L7 routing) → Service (load balance ภายใน) → Pods**

### 4.5 Network Policy — จำกัดการคุยกันระหว่าง pod

โปรเจกต์มี `kafka-network-policy.yaml` เพื่อกำหนดว่าใครคุยกับ Kafka ได้บ้าง — เพิ่ม security layer ไม่ให้ pod แปลกปลอมเข้าถึง broker

### สรุปวิธี deploy เรียงตามความพร้อม production

| วิธี | เหมาะกับ | ความซับซ้อน | scaling/healing |
|---|---|---|---|
| Docker Compose | dev local | ต่ำ | ❌ |
| K8s Deployment | stateless service | กลาง | ✅ |
| K8s StatefulSet | Kafka/DB/Zookeeper | สูง | ✅ + ข้อมูลคงที่ |
| + Ingress + NetworkPolicy + Secret | production | สูง | ✅ ครบ |

---

## 5. เส้นทางของหนึ่ง request — ทุกอย่างทำงานร่วมกัน

ลองดูตอนสมัครสมาชิกใหม่ใน Myblogs ว่าทั้งสามตัวเข้ามาเกี่ยวยังไง:

```
1. client ──HTTP──▶ Ingress ──▶ Service(user) ──▶ pod userservice
2. userservice เขียน user ลง Postgres
3. userservice (Producer) ──▶ Kafka topic "user.created"
   (ถ้า Kafka ล่ม → ข้าม, user ยังสมัครสำเร็จ)
4. blogservice (Consumer group) ดึง event จาก Kafka
5. blogservice เขียน blog_users ลง Postgres + ล้าง Redis cache
6. ครั้งต่อไปที่ client เรียก GET /getBlogs:
     - cache hit  → คืนจาก Redis (~ms)
     - cache miss → query Postgres → เก็บลง Redis (TTL 10s)
```

ทั้งหมดนี้รันอยู่บน Kubernetes ที่คอย restart pod ที่ล่ม, scale ตามโหลด, และ rolling update เวลา deploy เวอร์ชันใหม่

---

## 6. Best Practices & ข้อควรระวัง

**Kafka**
- ทำ consumer ให้ **idempotent** (เพราะเป็น at-least-once)
- ตั้ง `replication.factor > 1` ใน production (โปรเจกต์นี้ตั้ง 1 = dev only)
- พิจารณาเลิกใช้ Zookeeper หันไปใช้ **KRaft mode** (Kafka รุ่นใหม่ไม่ต้องพึ่ง Zookeeper แล้ว)

**Redis**
- เปลี่ยน `KEYS` → `SCAN` เพื่อไม่ block Redis
- ตั้ง TTL ให้เหมาะกับความสด vs ภาระ DB (ตอนนี้ 10s)
- ใส่ password (`requirepass`) — โค้ดตอนนี้ comment ไว้ `// Password:`
- ระวัง **cache stampede** เมื่อ key หมดอายุพร้อมกันใต้โหลดสูง

**Kubernetes**
- ใส่ `resources.requests/limits` (ตอนนี้ `resources: {}` ว่างอยู่)
- ใส่ **liveness/readiness probe** เพื่อให้ K8s รู้ว่า pod พร้อมรับ traffic จริง
- เก็บ secret ทั้งหมดใน `Secret` ไม่ใช่ env value ตรง ๆ
- ตั้ง replicas > 1 สำหรับ stateless service เพื่อ high availability

---

## สรุป

- **Kafka** ทำให้ service คุยกันแบบ async และไม่ล้มตามกัน (decoupling + resilience)
- **Redis** ลดภาระ DB ด้วย cache-aside และเป็น best-effort ที่ไม่ทำให้ระบบพังถ้าแคชล่ม
- **Kubernetes** มัดทุกอย่างเข้าด้วยกัน — Deployment สำหรับ stateless, StatefulSet สำหรับ stateful, Service/Ingress สำหรับ networking, Secret/NetworkPolicy สำหรับ security

จุดแข็งของสถาปัตยกรรมนี้คือ **ทุกชั้นออกแบบให้ "fail gracefully"** — Kafka ล่ม user ยังสมัครได้, Redis ล่มยังอ่าน blog ได้, pod ล่ม K8s restart ให้ นี่คือหัวใจของระบบ distributed ที่ดูแลง่ายและทนทาน
