# HU-017-002 - Sincronizacion Offline hacia Online

## 📌 Informacion General
- ID: HU-017-002
- Epic: EPICA-017
- Prioridad: Alta
- Estado: Backlog
- Porcentaje: 0%
- Autor: QA Team
- Fecha: 2026-05-04

---

## 👤 Historia de Usuario

**Como** operador de punto de venta offline  
**Quiero** que los productos creados o actualizados en la version offline (SQLite) se sincronicen automaticamente hacia la base de datos online (PostgreSQL) cuando tenga conexion  
**Para** mantener los datos actualizados en el servidor central sin necesidad de sincronizacion manual

---

## 🧠 Descripcion Funcional

Cuando un cliente offline hace cambios en productos (crear, actualizar, eliminar), el sistema debe enviar estos cambios al servidor online. El servidor online procesara los cambios y actualizara la base de datos PostgreSQL.

Para manejar clientes offline que pueden estar desconectados por periodos prolongados, se implementara un endpoint de sincronizacion que el cliente offline puede llamar cuando tengas conexion.

---

## ✅ Criterios de Aceptacion

### Escenario 1: Sincronizar producto creado en offline
- Dado que el cliente offline ha creado un nuevo producto localmente
- Cuando el cliente offline se conecta al servidor online
- Y envia una solicitud POST a /offline/sync/products
- Entonces el servidor online crea el producto en PostgreSQL
- Y devuelve confirmacion al cliente offline
- Y el producto esta disponible para otros clientes

### Escenario 2: Sincronizar producto actualizado en offline
- Dado que el cliente offline ha modificado un producto existente
- Cuando sincroniza los cambios hacia online
- Entonces el servidor online actualiza el producto en PostgreSQL
- Y se publica el evento "product.updated" para otros clientes offline

### Escenario 3: Sincronizar producto eliminado en offline
- Dado que el cliente offline ha eliminado un producto localmente
- Cuando sincroniza hacia online
- Entonces el servidor online realiza soft delete del producto

### Escenario 4: Cliente offline con multiples cambios pendientes
- Dado que el cliente offline tiene 50 productos modificados desde la ultima sincronizacion
- Cuando sincroniza todos los cambios
- Entinces el servidor procesa todos los productos en una sola peticion
- Y devuelve un resumen de operaciones exitosas y fallidas

### Escenario 5: Reintento de sincronizacion
- Dado que la sincronizacion falla por error de red
- Cuando el cliente offline reintenta la solicitud
- Entonces no se crean duplicados (el servidor usa upsert)
- Y los datos se mantienen consistentes

### Escenario 6: Autenticacion requerida
- Dado que un cliente offline no autenticado intenta sincronizar
- Cuando envia solicitud de sincronizacion
- Entonces el servidor devuelve 401 Unauthorized
- Y rechaza la sincronizacion

---

## ❌ Casos de Error

- Si el producto tiene datos invalidos: El servidor devuelve 400 con los errores de validacion
- Si el tenant_slug no corresponde al token: El servidor devuelve 403 Forbidden
- Si el producto tiene SKU duplicado: Se retorna error 409 Conflict
- Si la conexion a base de datos falla: El servidor devuelve 500 y el cliente reintenta
- Si el servidor recibe un ID que no existe para actualizacion: Se retorna 404 Not Found

---

## 🔐 Reglas de Negocio

- El cliente offline debe enviar JWT valido con el tenant_slug
- El servidor usa upsert (INSERT ON CONFLICT UPDATE) para evitar duplicados
- El timestamp del evento se usa para resolucion de conflictos
- Solo se aceptan productos con active = true
- El servidor valida el SKU segun las reglas existentes
- Se publica evento "product.offline.created/updated/deleted" despues de procesar

---

## 🎨 Consideraciones UI/UX

- No aplica (API backend)

---

## 📡 Requisitos Tecnicos

### Endpoint de Sincronizacion

```
POST /offline/sync/products
Authorization: Bearer {jwt_token}
Content-Type: application/json

Request:
{
  "tenant_slug": "empresa_uno",
  "products": [
    {
      "product_id": 123,
      "sku": "PROD-001",
      "barcode": "7501234567890",
      "name": "Producto de ejemplo",
      "action": "create",
      "timestamp": "2026-05-04T10:30:00Z",
      "source": "offline"
    }
  ]
}

Response - Exito (200):
{
  "success": true,
  "message": "Sincronizacion completada",
  "data": {
    "results": [
      {
        "product_id": 123,
        "sku": "PROD-001",
        "action": "create",
        "status": "success",
        "message": "Product created successfully"
      }
    ],
    "total_synced": 1,
    "success_count": 1,
    "error_count": 0,
    "sync_time": "2026-05-04T10:35:00Z"
  }
}

Response - Error (400):
{
  "success": false,
  "message": "Error de validacion",
  "data": {
    "results": [...],
    "total_synced": 1,
    "success_count": 0,
    "error_count": 1
  }
}
```

### Eventos a publicar (post-procesamiento)

| Evento | Routing Key | Descripcion |
|-------|-------------|-------------|
| product.offline.created | product.offline.created.{slug} | Producto creado desde offline |
| product.offline.updated | product.offline.updated.{slug} | Producto actualizado desde offline |
| product.offline.deleted | product.offline.deleted.{slug} | Producto eliminado desde offline |

### Payload del evento

```json
{
  "event": "product.offline.created",
  "tenant_slug": "empresa_uno",
  "enterprise_id": 123,
  "product_id": 456,
  "source": "offline",
  "timestamp": "2026-05-04T10:35:00Z",
  "data": { ... }
}
```

### Cambios requeridos

- modules/offline/handler.go: Agregar endpoint POST /offline/sync/products
- modules/offline/service.go: Agregar metodo SyncProductsFromOffline()
- Validar JWT y extraer tenant_slug del token
- Iterar sobre productos y procesar cada operacion
- Devolver resumen de operaciones

---

## 🧪 Criterios de Testing

- Unit tests para cada tipo de operacion (create, update, delete)
- Integration tests con base de datos online
- Tests de concurrencia (multiples clientes offline)
- Tests de autenticacion (JWT invalido, token expirado)
- Tests de validacion (SKU duplicado, datos invalidos)

---

## 📎 Dependencias

- modules/offline/handler.go: Endpoint existente
- modules/offline/service.go: Metodos sync existentes
- modules/catalog/products/domain.go: Entidad Product existente
- RabbitMQ ya implementado para publicacion de eventos

---

## 🚫 Fuera de Alcance

- No incluye sincronizacion de otras entidades (HU-017-002 fuera de alcance)
- No incluye resolucion de conflictos avanzada (HU-017-003)
- No incluye sincronizacion bidireccional en tiempo real (solo batch)
- No incluye interface de usuario

---

## 🧠 Generacion de Codigo

Requerir:
- Handler: Endpoint POST en routes.go
- Handler: Request/Response structs en handler.go
- Service: Metodo SyncProductsFromOffline() en service.go
- Repository: UpsertProductoDesdeOffline() en repository.go
- Publicar eventos post-procesamiento
- Tests unitarios y de integracion