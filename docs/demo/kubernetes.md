# Kubernetes — หลักฐานว่า deploy บน cluster ได้

manifest ทั้งหมดอยู่ที่ [`blogservice/k8s-files/`](../../blogservice/k8s-files/) ครอบคลุม

- `kafka/` — Kafka + Zookeeper เป็น **StatefulSet** (มี persistent volume)
- `postgresql/` — PostgreSQL + PV / PVC / Secret
- `redis/` — Redis

แนวทางนี้คือ "ยกทุกอย่างเข้า cluster" — ต่างจากตอน deploy บน PaaS ที่ต้องพึ่งบริการข้างนอก

## วิธีพิสูจน์ (ใช้ k8s local ก็ได้: Docker Desktop / minikube)

### 1. เปิด cluster แล้วเช็ค node พร้อม

```powershell
kubectl get nodes
```

node ต้องสถานะ `Ready`

📸 _screenshot:_ `images/k8s-nodes-ready.png`

### 2. apply manifest แล้วดูทุก resource ขึ้นครบ

```powershell
kubectl apply -f blogservice/k8s-files/postgresql/
kubectl apply -f blogservice/k8s-files/kafka/
kubectl apply -f blogservice/k8s-files/redis/
```

```powershell
kubectl get all
```

ต้องเห็น pod ของ kafka / zookeeper / postgres / redis สถานะ `Running`

📸 _screenshot:_ `images/k8s-get-all-running.png`

### 3. ดู log ของ pod ว่าทำงานจริง

```powershell
kubectl logs statefulset/kafka
```

📸 _screenshot:_ `images/k8s-kafka-logs.png`

### 4. (เสริม) โชว์ว่า self-healing ทำงาน

ลองลบ pod ทิ้งแล้วดู k8s สร้างใหม่ให้อัตโนมัติ

```powershell
kubectl delete pod <ชื่อ-pod>
kubectl get pods -w
```

จะเห็น pod ใหม่ถูกสร้างแทนทันที — เป็นจุดเด่นที่ docker-compose ทำไม่ได้

📸 _screenshot:_ `images/k8s-self-healing.png`

## สิ่งที่หลักฐานนี้บอก

- แยกแยะ **Deployment** (stateless) กับ **StatefulSet** (stateful เช่น Kafka/DB) ได้
- เข้าใจ Service / PV / PVC / Secret
- เข้าใจจุดต่างระหว่าง k8s กับ docker-compose (scaling, self-healing)

> หมายเหตุ: manifest ชุดนี้สาธิตการรัน infra (Kafka/Redis/Postgres) บน cluster
> ส่วนการ apply `deploy.yaml` ของ user/blog service ต้อง build image เองก่อน
