# HU-017-001 - Sincronizacion Automatica Online hacia Offline

## 📌 Informacion General
- ID: HU-017-001
- Epic: EPICA-017
- Prioridad: Alta
- Estado: Backlog
- Porcentaje: 0%
- Autor: QA Team
- Fecha: 2026-05-04

---

## 👤 Historia de Usuario

**Como** sistema de Punto de Venta  
**Quiero** que los productos creados, actualizados o eliminados en la version online (PostgreSQL) se reflejen automaticamente en los clientes offline (SQLite) en tiempo real  
**Para** que los operadores de tienda siempre tengan informacion actualizada sin necesidad de sincronizacion manual

---

## 🧠 Descripcion Funcional

Cuando un producto es creado, actualizado o eliminado en la base de datos PostgreSQL (online), el servicio de productos debe publicar un evento en RabbitMQ que sera consumido por todos los clientes offline suscritos al tenant. El cliente offline recibira el evento y actualizara su base de datos SQLite local.

El sistema ya cuenta con los eventos de dominio definidos:
- product.created
- product.updated  
- product.deleted

Solo falta implementar la publicacion a RabbitMQ y el consumo en el cliente offline.

---

## ✅ Criterios de Aceptacion

### Escenario 1: Crear producto en online se sincroniza a offline
- Dado que el usuario ha iniciado sesion en el sistema online
- Cuando crea un nuevo producto con datos validos
- Entonces el producto se guarda en PostgreSQL
- Y se publica el evento "product.created" en RabbitMQ
- Y el cliente offline recibe el evento en menos de 5 segundos
- Y el cliente offline guarda el producto en SQLite

### Escenario 2: Actualizar producto en online se sincroniza a offline
- Dado que existe un producto en la base de datos online
- Cuando el usuario actualiza el producto
- Entonces el producto se actualiza en PostgreSQL
- Y se publica el evento "product.updated" en RabbitMQ
- Y el cliente offline recibe el evento y actualiza su copia local

### Escenario 3: Eliminar producto en online se sincroniza a offline
- Dado que existe un producto en la base de datos online
- Cuando el usuario elimina el producto (soft delete)
- Entonces el producto se marca como eliminado en PostgreSQL
- Y se publica el evento "product.deleted" en RabbitMQ
- Y el cliente offline recibe el evento y marca como eliminado localmente

### Escenario 4: Cliente offline desconectado
- Dado que el cliente offline esta desconectado de RabbitMQ
- Cuando se publica un evento en RabbitMQ
- Entonces el evento no se pierde (RabbitMQ guarda en cola)
- Y cuando el cliente se reconecta recibe los eventos pendientes

### Escenario 5: Multi-tenant isolation
- Dado que existen multiple empresas (tenants) en el sistema
- Cuando se crea un producto para "empresa_uno"
-Entonces solo los clientes de "empresa_uno" reciben el evento
- Y los clientes de otras empresas no lo reciben

---

## ❌ Casos de Error

- Si RabbitMQ no esta disponible: El sistema debe registrar el error y reintentar la publicacion
- Si el payload del evento es invalido: El evento se descarta y se registra el error
- Si el tenant no existe: El evento se ignora silenciosamente
- Si el cliente offline no puede procesar el evento: Se reintenta hasta 3 veces

---

## 🔐 Reglas de Negocio

- El evento debe incluir el tenant_slug para aislamiento multi-tenant
- El timestamp del evento usa formato ISO 8601 UTC
- El nombre del evento sigue el patron: product.{action}
- La cola de RabbitMQ sigue el patron: aura.{env}.product.{tenant}
- Solo se sincronizan productos con deleted_at IS NULL en la consulta inicial
- El cliente offline debe confirmar recepcional del evento (ACK)

---

## 🎨 Consideraciones UI/UX

- No aplica (API backend)

---

## 📡 Requisitos Tecnicos

### Eventos a publicar

| Evento | Routing Key | Payload Keys |
|-------|------------|--------------|
| product.created | product.created.{slug} | id, sku, barcode, name, description, category_id, brand_id, unit_id, active, cost_price, sale_price, created_at, updated_at |
| product.updated | product.updated.{slug} | id, sku, barcode, name, description, category_id, brand_id, unit_id, active, cost_price, sale_price, updated_at |
| product.deleted | product.deleted.{slug} | id, deleted_at |

### Cambios en modulo products

- Modificar service.go para injectar eventBus en constructor
- Agregar llamada a eventBus.Publish() en Create(), Update(), Delete()
- El eventBus ya esta configurado en NewService()

### Cambios en modulo offline

- Suscribirse a los eventos de producto en RabbitMQ
- Crear handler para procesar eventos recibidos
- Actualizar base de datos SQLite con los datos del evento

---

## 🧪 Criterios de Testing

- Unit tests para service.Create() publicando evento exitoso
- Unit tests para service.Update() publicando evento exitoso
- Unit tests para service.Delete() publicando evento exitoso
- Integration tests con RabbitMQ mock
- Tests de concurrencia (multiples productos)
- Tests de red (simular desconexion)

---

## 📎 Dependencias

-modules/catalog/products/domain.go: Ya tiene los eventos definidos
- modules/catalog/products/service.go: Ya tiene el metodo publish()
- infrastructure/messaging/rabbit/: Ya implementado
- RabbitMQ debe estar configurado en ambiente

---

## 🚫 Fuera de Alcance

- No incluye sincronizacion desde offline hacia online (es HU-017-002)
- No incluye resolucion de conflictos (es HU-017-003)
- No incluye otras entidades (terceros, categorias, marcas)

---

## 🧠 Generacion de Codigo

Requerir:
- Modulo products: Agregar eventBus en constructor NewService()
- Modulo products: Llamar publish() en cada operacion CRUD
- Modulo offline: Suscribir a eventos RabbitMQ
- Modulo offline: Implementar handler para procesar eventos
- Tests unitarios para publicacion de eventos
- Tests de integracion con RabbitMQ