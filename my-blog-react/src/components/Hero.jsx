export default function Hero() {
  return (
    <section className="wrap pt-16 pb-12 sm:pt-24 sm:pb-16">
      <p className="mb-4 font-mono text-xs font-medium uppercase tracking-[0.2em] text-accent">
        Myblogs
      </p>
      <h1 className="max-w-3xl text-4xl leading-[1.1] sm:text-6xl">
        Writing about backend, Go, and distributed systems.
      </h1>
      <p className="mt-6 max-w-xl text-lg leading-relaxed text-muted">
        บล็อกที่ขับเคลื่อนด้วยสถาปัตยกรรม microservices — Go, Kafka, Redis และ React
        บันทึกสิ่งที่ได้เรียนรู้ระหว่างสร้างระบบจริง
      </p>
    </section>
  );
}
