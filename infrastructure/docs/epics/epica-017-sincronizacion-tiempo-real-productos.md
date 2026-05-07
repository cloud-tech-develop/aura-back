# EPICA-017 - Sincronizacion en Tiempo Real de Productos

## 📌 Informacion General
- ID: EPICA-017
- Estado: Backlog
- Prioridad: Alta
- Fecha inicio: 2026-05-04
- Fecha objetivo: 2026-06-30
- Owner: QA Team
- Porcentaje: 0%

---

## 🎯 Objetivo de Negocio

Permitir la sincronizacion bidireccional de productos entre las versiones online (PostgreSQL) y offline (SQLite) del sistema Aura POS en tiempo real, eliminando la necesidad de sincronizacion manual y garantizando la integridad de datos entre ambas versiones.

**Problema actual**: Los usuarios deben ejecutar sincronizacion manual para obtener los ultimos productos en la version offline, y los cambios realizados offline no se reflejan en produccion hasta una sincronizacion manual.

**Valor generado**:
- Actualizacion automatica de productos en clientes offline cuando se modifica en online
- Sincronizacion de cambios offline hacia online cuando hay conexion
- Reduced time de inactividad por falta de datos actualizados
- Melhor experiencia de usuario para puntos de venta distribuidos

---

## 👥 Stakeholders

- Usuario final: Operadores de POS en tiendas con conexion inestable
- Equipo tecnico: Backend developers, mobile developers
- Producto: Product managers de Aura POS

---

## 🧠 Descripcion Funcional General

Implementar un sistema de sincronizacion bidireccional basado en eventos utilizando RabbitMQ:

1. **Online → Offline**: Cuando se crea, actualiza o elimina un producto en la base de datos PostgreSQL, el sistema publicara un evento que sera consumido por los clientes offline conectados.

2. **Offline → Online**:Cuando un cliente offline hace cambios en productos, estos se sincronizaran hacia la base de datos online cuando tenga conexion disponible.

3. **Resolucion de conflictos**: Sistema para detectar y resolver conflictos de datos cuando ambas versiones tienen cambios en el mismo registro.

---

## 📦 Alcance

Incluye:
- Eventos de dominio para operaciones CRUD de productos (ya existen)
- Publicacion de eventos RabbitMQ desde el servicio de productos
- Endpoint para recibir sincronizacion desde offline
- Cola de sincronizacion pendiente para clientes offline desconectados
- Logica de resolucion de conflictos basada en timestamps
- Receptor de eventos en modulo offline

No incluye:
- Sincronizacion de otras entidades (terceros, usuarios) - es fuera de alcance
- Replicacion completa de base de datos
- Sincronizacion de archivos/media

---

## 🧩 Historias de Usuario Asociadas

- [ ] HU-017-001 - Sincronizacion automatica online hacia offline
- [ ] HU-017-002 - Sincronizacion offline hacia online
- [ ] HU-017-003 - Resolucion de conflictos de datos

---

## 🐞 Bugs Asociados

No hay bugs asociados inicialmente.

---

## 🔐 Reglas de Negocio Globales

- Todos los eventos de sincronizacion deben incluir el tenant_slug para aislamiento multi-tenant
- Los clientes offline deben authenticarse mediante token JWT
- Los timestamps de actualizacion se usan para resolucion de conflictos
- La cola de RabbitMQ sigue el patron: aura.{env}.{event}.{tenant}
- Solo se sincronizan productos activos y no eliminados logicamente

---

## 🧱 Arquitectura Relacionada

Frontend: N/A
Backend:
  - modules/catalog/products/ - Event publishing en Create/Update/Delete
  - modules/offline/ - Receptor de eventos y sync endpoint
  - infrastructure/messaging/rabbit/ - Event bus
Base de datos:
  - PostgreSQL (online) - Schema per-tenant
  - SQLite (offline) - Base de datos local
Autenticacion: JWT con tenant_slug

---

## 📊 Metricas de Exito

- Latencia de sincronizacion online→offline < 5 segundos
- Tasas de exito de sincronizacion > 95%
- Tiempo de resolucion de conflictos < 2 segundos

---

## 🚧 Riesgos

- Perdida de eventos por desconexion de RabbitMQ
- Conflictos de datos por ediciones simultaneas
- Latencia en redes lentas

---

## 📎 Eventos de Sincronizacion

### Eventos de Producto

| Evento | Descripcion | Direction |
|--------|------------|-----------|
| product.created | Producto nuevo creado en online | Online → Offline |
| product.updated | Producto modificado en online | Online → Offline |
| product.deleted | Producto eliminado en online | Online → Offline |
| product.offline.created | Producto creado en offline | Offline → Online |
| product.offline.updated | Producto modificado en offline | Offline → Online |
| product.offline.deleted | Producto eliminado en offline | Offline → Online |
| product.sync.conflict | Conflicto detectado | Bidireccional |

### Formato de Payload

```json
{
  "event": "product.created",
  "tenant_slug": "empresa_uno",
  "enterprise_id": 123,
  "product_id": 456,
  "timestamp": "2026-05-04T10:30:00Z",
  "data": { ... }
}
```