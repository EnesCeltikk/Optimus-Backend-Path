# Deployment & DevOps

Bu klasör, projelerin (API Implementation, Project Setup and Basic Structure, Enes_Celtik_CoreImplementation) Docker ile containerize edilmesi, servislerin orkestrasyonu ve monitoring için gerekli dosyaları içerir.

## İçerik
- **Dockerfile**: Go uygulamasını çok aşamalı olarak derler ve çalıştırır.
- **docker-compose.yml**: Uygulama, veritabanı, Redis, Prometheus ve Grafana servislerini başlatır.
- **prometheus.yml**: Prometheus monitoring konfigürasyonu.
- **grafana/**: Grafana dashboard jsonları veya provisioning scriptleri.
- **.env.example**: Ortak environment değişkenleri örneği.

## Kullanım
1. Proje kökünde `.env` dosyanı oluştur (örneği `.env.example`dan kopyala).
2. `docker-compose up --build` komutunu çalıştır.
3. Prometheus ve Grafana arayüzlerinden monitoring ve dashboardlara eriş.

## Notlar
- Dockerfile ve compose dosyası, tüm ana Go projeleriyle uyumlu olacak şekilde parametrik hazırlanmıştır.
- Monitoring ve tracing için Prometheus ve Grafana entegre edilmiştir. 