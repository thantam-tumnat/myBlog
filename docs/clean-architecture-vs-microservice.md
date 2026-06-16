# Clean Architecture กับ Microservice: คนละชั้น ไม่ได้ทับซ้อนกัน

> อ้างอิงจากโปรเจกต์จริง `Myblogs-Microservices` — ReactJS (Vite) + Golang (gofiber) + PostgreSQL + Redis + Kafka, deploy ด้วย Docker / Kubernetes

เวลาเริ่มทำระบบ backend ใหม่ ๆ คำถามที่เจอบ่อยคือ "เราจะใช้ Clean Architecture หรือจะทำ Microservice?" — ซึ่งเป็นคำถามที่ตั้งผิดตั้งแต่แรก เพราะสองอย่างนี้ **ไม่ใช่ทางเลือกที่ต้องเลือกอย่างใดอย่างหนึ่ง** มันอยู่กันคนละระดับและทำงานเสริมกัน

บทความนี้จะอธิบาย 3 เรื่อง:
1. ทำไมสองคอนเซปต์นี้ถึงไม่ทับซ้อนกัน
2. การ deploy ของแต่ละแบบต่างกันอย่างไร
3. มุมมองเรื่อง repository (โครงสร้างโค้ดในที่เก็บ)

---

## 1. คนละชั้น ไม่ได้ทับซ้อนกัน

วิธีจำง่าย ๆ:

- **Microservice เป็นเรื่องระดับ "ระบบ" (system / runtime)** — ตอบคำถามว่า *ระบบของเราถูกแบ่งออกเป็นกี่กระบวนการ (process) ที่รันแยกกัน, deploy แยกกัน, คุยกันผ่านอะไร*
- **Clean Architecture เป็นเรื่องระดับ "ภายในแต่ละ service" (code / module)** — ตอบคำถามว่า *โค้ดภายใน service หนึ่งตัว ถูกจัดวางเป็นชั้น ๆ อย่างไร ใครพึ่งพาใคร*

พูดอีกแบบ: Microservice บอกว่าระบบมี "กี่กล่อง" ส่วน Clean Architecture บอกว่า "ในแต่ละกล่องจัดของยังไง"

### ในโปรเจกต์นี้เห็นทั้งสองอย่างพร้อมกัน

ระดับ Microservice — ระบบถูกแบ่งเป็น 3 ส่วนที่รันแยกกัน:

```
my-blog-react/   → frontend (React + Vite)
userservice/     → จัดการ user, produce event ไปยัง Kafka
blogservice/     → จัดการ blog, consume event จาก userservice
```

`userservice` และ `blogservice` ต่างมี `go.mod`, `Dockerfile`, `deploy.yaml` เป็นของตัวเอง และคุยกันแบบ asynchronous ผ่าน Kafka (ฝั่ง user เป็น producer, ฝั่ง blog เป็น consumer) — นี่คือลักษณะ microservice ชัดเจน: deploy แยก, scale แยก, ฐานข้อมูล/แคชแยกได้

ระดับ Clean Architecture — *ภายใน* `blogservice` ตัวเดียว โค้ดยังถูกจัดเป็นชั้น:

```
modules/blogs/
  controllers/    → blogHandler.go      (ชั้นนอกสุด รับ HTTP, แปลง request/response)
  usecases/       → blog_usecase.go     (business logic: BlogCreated, BlogGets)
  repositories/   → blog_repositories.go (เข้าถึง DB/Redis จริง)
modules/entities/ → domain.go           (interface + struct ที่เป็นแกนกลาง)
```

หัวใจของ Clean Architecture คือ **ทิศทางการพึ่งพา (dependency direction) ชี้เข้าด้านใน** สังเกตที่ `entities/domain.go`:

```go
type BlogRepository interface {
    CreateBlog(blog *Blog) (*Blog, error)
    GetBlogs() ([]Blog, error)
    // ...
}

type BlogService interface {
    BlogCreated(userId int, blogReq *BlogRequest) (*BlogRes, error)
    BlogGets() (*[]BlogRes, error)
}
```

`usecase` ไม่ได้ผูกกับ Redis หรือ Postgres โดยตรง แต่พึ่งพา **interface** `entities.BlogRepository`:

```go
type blogService struct {
    blogRepo entities.BlogRepository  // interface ไม่ใช่ concrete type
    cfg      *configs.Configs
}
```

ส่วนการเลือกว่าจะใช้ implementation ตัวไหน (Redis + Postgres) ถูกประกอบเข้าด้วยกันที่ชั้นนอกใน `servers/handler.go`:

```go
blogRepository := _blogRepo.NewRepositoryRedis(s.Redis, s.Db)
blogUsecases   := _blogUse.NewUserService(blogRepository, s.Cfg)
```

ผลคือ business logic ใน usecase ทดสอบได้โดยไม่ต้องมี DB จริง (mock interface แทน) และเปลี่ยน Redis เป็นอย่างอื่นได้โดยไม่แตะ usecase เลย — นี่คือประโยชน์ของ Clean Architecture ที่เกิด *ภายในกล่องเดียว* ไม่เกี่ยวกับว่าระบบจะมีกี่ service

> **สรุปข้อ 1:** คุณจะมี Monolith ที่ใช้ Clean Architecture ก็ได้, มี Microservice ที่ข้างในแต่ละตัวเขียนแบบ spaghetti ก็ได้ (แต่ไม่ควร), หรือมี Microservice ที่ข้างในเป็น Clean Architecture แบบโปรเจกต์นี้ก็ได้ — เพราะมันอยู่คนละแกน

---

## 2. การ deploy ของแต่ละแบบ

จุดที่ความต่างชัดที่สุดคือ "ตอน deploy" เพราะ Clean Architecture **ไม่มีผลต่อหน่วยการ deploy เลย** ในขณะที่ Microservice **คือเรื่องของหน่วยการ deploy โดยตรง**

### Clean Architecture: ทุกชั้น compile รวมเป็น artifact เดียว

ชั้น controller / usecase / repository / entity ของ `blogservice` ทั้งหมดถูก build รวมเป็น binary เดียวผ่าน `Dockerfile` ตัวเดียว → ได้ container image หนึ่งตัว ตอนรันมันคือ process เดียว เส้นแบ่งชั้นเป็นเพียง **เส้นแบ่งทางตรรกะในโค้ด** ไม่ใช่เส้นแบ่ง network คุณ deploy "ชั้น usecase" แยกจาก "ชั้น repository" ไม่ได้ และไม่มีเหตุผลต้องทำ

### Microservice: แต่ละ service = หน่วย deploy อิสระ

ตรงนี้คือของจริง สังเกตว่าทั้ง `userservice` และ `blogservice` ต่างมีไฟล์ deploy ของตัวเอง:

```
userservice/Dockerfile   + userservice/deploy.yaml
blogservice/Dockerfile   + blogservice/deploy.yaml
my-blog-react/Dockerfile + my-blog-react/deploy.yaml
```

ใน `blogservice/deploy.yaml` มี **Deployment + Service + Ingress** ครบเป็นชุดเฉพาะของ blog:

```yaml
kind: Deployment
metadata:
  name: blog
spec:
  template:
    spec:
      containers:
      - name: blog
        image: ghcr.io/thantam-tumnat/blogservice:latest
        ports:
        - containerPort: 8002
  strategy:
    type: RollingUpdate     # อัปเดต blog ทีละ pod โดยไม่กระทบ user service
    rollingUpdate:
      maxUnavailable: 1
      maxSurge: 1
```

นี่คือสิ่งที่ทำให้ deploy แบบ microservice ต่างจาก monolith:

| ประเด็น | Clean Architecture (ภายใน service) | Microservice (ระหว่าง service) |
|---|---|---|
| หน่วย deploy | 1 artifact ต่อ service | หลาย artifact อิสระ |
| Scale | scale ทั้งก้อน | scale เฉพาะตัวที่โหลดหนัก (เช่น เพิ่ม replica ของ blog อย่างเดียว) |
| Rolling update | อัปทั้ง process | อัป `blog` โดยไม่แตะ `user` ได้ |
| การสื่อสาร | function call ในหน่วยความจำ | ข้ามเครือข่าย (HTTP / Kafka) |
| ขอบเขตพัง | พังทั้ง process | พังเฉพาะ service (ถ้าออกแบบ resilient พอ) |

### Local vs Production: สองสไตล์ในโปรเจกต์เดียว

โปรเจกต์นี้รองรับ deploy 2 ระดับ:

- **Local / dev → `docker-compose.yml`** รวม dependency พื้นฐาน (postgres, redis, zookeeper, kafka) ขึ้นมาด้วยคำสั่งเดียว เหมาะกับนักพัฒนาที่อยากรันทั้งระบบบนเครื่องตัวเองเร็ว ๆ
- **Production → ไฟล์ใน `k8s-files/` + `deploy.yaml`** แตกย่อยเป็น Deployment, Service, StatefulSet (kafka/zookeeper), PV/PVC (postgres), NetworkPolicy — ออกแบบให้ scale และจัดการ lifecycle แยกชิ้นได้บน Kubernetes

จุดที่น่าสนใจในโค้ดคือ service ถูกเขียนให้ **ทนต่อ environment ที่ dependency ไม่ครบ** — ใน `blogservice/modules/servers/servers.go`:

```go
// รัน Kafka consumer เฉพาะเมื่อมี broker เชื่อมต่อ
if s.ConsumerGroup != nil {
    go func() { /* consume loop */ }()
} else {
    logs.Info("Kafka consumer disabled (no broker connection)")
}
```

นี่คือ mindset แบบ microservice: service ต้องยังให้บริการ HTTP ได้แม้ dependency บางตัว (Kafka) จะไม่พร้อม — เพื่อให้ deploy ในสภาพแวดล้อมที่ตัด Kafka ออก (เช่น Render) ยังทำงานได้

> **สรุปข้อ 2:** Clean Architecture จัดระเบียบโค้ด *ก่อน* build — ไม่เปลี่ยนจำนวน artifact. Microservice จัดระเบียบ *หน่วย deploy* — เปลี่ยนทั้งวิธี build, scale, อัปเดต, และการสื่อสาร

---

## 3. มุมมองเรื่อง Repository

เมื่อระบบเป็น microservice คำถามต่อมาคือ "แล้วโค้ดจะเก็บใน repo อย่างไร?" มี 2 แนวหลัก

### Monorepo — repo เดียว หลาย service (โปรเจกต์นี้ใช้แบบนี้)

โปรเจกต์ `Myblogs-Microservices` เก็บทุกอย่างไว้ใน repo เดียว:

```
Myblogs-Microservices/
├── docker-compose.yml      ← orchestrate ทั้งระบบจากจุดเดียว
├── init-db.sql
├── userservice/            ← go.mod ของตัวเอง
├── blogservice/            ← go.mod ของตัวเอง
└── my-blog-react/          ← package.json ของตัวเอง
```

ข้อสังเกตสำคัญ: **เป็น monorepo แต่ละ service ยังแยก dependency กันชัดเจน** — แต่ละ Go service มี `go.mod`/`go.sum` ของตัวเอง, frontend มี `package.json` แยก นั่นคือ "หนึ่ง repo หลายโมดูลอิสระ" ไม่ใช่ "ทุกอย่างผูกเป็นก้อนเดียว"

**ข้อดีของ monorepo สำหรับโปรเจกต์ขนาดนี้:**
- เปิด repo เดียวเห็นทั้งระบบ เข้าใจ flow `user → kafka → blog` ได้ในที่เดียว
- แก้ contract ของ event (เช่น struct ใน `entities/event.go` ทั้งสองฝั่ง) ใน commit/PR เดียว เห็น diff พร้อมกัน ลดปัญหา version ไม่ตรง
- ตั้ง `docker-compose.yml` ที่ระดับ root เพื่อยกทั้งระบบขึ้นมาทดสอบได้ทันที
- เหมาะกับทีมเล็ก / โปรเจกต์ portfolio ที่อยากให้คนอ่านเข้าใจภาพรวมเร็ว

**ข้อควรระวัง:**
- ต้องมีวินัยไม่ให้ service หนึ่ง `import` โค้ดข้ามไปอีก service โดยตรง (ต้องคุยผ่าน API/Kafka เท่านั้น) ไม่งั้นจะกลายเป็น monolith แอบแฝง
- เมื่อทีม/บริการโตขึ้นมาก ๆ CI/CD ต้องฉลาดพอที่จะ build เฉพาะ service ที่เปลี่ยน

### Polyrepo (multi-repo) — หนึ่ง service หนึ่ง repo

อีกแนวคือแยก `userservice`, `blogservice`, `my-blog-react` เป็นคนละ repo

**ข้อดี:** ขอบเขตความเป็นเจ้าของ (ownership) ชัด, สิทธิ์เข้าถึงแยกได้, CI/CD ของแต่ละ repo เป็นอิสระจริง ๆ, เหมาะกับหลายทีมที่ดูแลคนละ service

**ข้อเสีย:** การแก้ที่กระทบหลาย service ต้องเปิดหลาย PR, sync contract ยากขึ้น, setup ระบบทั้งหมดบนเครื่อง dev ใหม่ ๆ ยุ่งกว่า

### เลือกอย่างไร

| สถานการณ์ | แนะนำ |
|---|---|
| ทีมเล็ก, อยากเห็นภาพรวมง่าย, contract เปลี่ยนบ่อย | **Monorepo** (แบบโปรเจกต์นี้) |
| หลายทีม, แต่ละ service มี lifecycle/release cycle ต่างกันชัด | **Polyrepo** |
| อยากได้ข้อดี monorepo แต่ build แยก | Monorepo + เครื่องมือ (Nx, Turborepo, Bazel, หรือ path-filter ใน CI) |

> **จุดสำคัญที่สุด:** การเลือก mono/poly repo เป็นเรื่อง *การจัดการโค้ดและทีม* — **ไม่ได้ตัดสินว่าระบบเป็น microservice หรือไม่** ระบบนี้เป็น microservice เต็มตัว (deploy แยก, คุยผ่าน Kafka, DB แยกได้) ทั้ง ๆ ที่อยู่ใน monorepo — เพราะความเป็น microservice วัดที่ *runtime* ไม่ใช่ที่ *จำนวน repo*

---

## บทสรุป

| มิติ | Clean Architecture | Microservice | Repo strategy |
|---|---|---|---|
| ระดับ | ภายใน 1 service (code) | ทั้งระบบ (runtime) | การเก็บโค้ด (org) |
| ตอบคำถาม | "จัดชั้นโค้ดยังไง" | "แบ่งเป็นกี่ process" | "เก็บใน repo ยังไง" |
| มีผลตอน deploy? | ไม่ (รวมเป็น 1 artifact) | ใช่ (หลายหน่วยอิสระ) | มีผลกับ CI/CD ไม่ใช่ runtime |
| ตัวอย่างในโปรเจกต์ | layering ใน `blogservice/modules` | `userservice` ⇄ Kafka ⇄ `blogservice` | monorepo ที่มี 3 โมดูลอิสระ |

สามแกนนี้ตั้งฉากกัน (orthogonal) เลือกได้อิสระต่อกัน โปรเจกต์ `Myblogs-Microservices` คือตัวอย่างที่ใช้ทั้งสามพร้อมกัน: **Clean Architecture ข้างในแต่ละ service, รันแบบ Microservice ที่คุยกันผ่าน Kafka, และเก็บทั้งหมดไว้ใน Monorepo** — ไม่มีอันไหนทับซ้อนหรือขัดกันเลย
