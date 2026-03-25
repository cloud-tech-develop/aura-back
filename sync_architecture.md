# Arquitectura de Sincronización Offline/Online (Aura POS)

Este documento explica el funcionamiento técnico de la sincronización entre la base de datos local (SQLite) y la versión online desplegada (PostgreSQL).

## 1. Estrategia de Base de Datos Dual

El sistema utiliza una capa de abstracción ([db.Querier](file:///C:/Users/Drako/Desktop/cloud-tecno/aura-back/internal/db/db.go)) que permite que el mismo código de negocio funcione sin cambios tanto en **Postgres** (producción online) como en **SQLite** (entorno local offline).

- **Online**: Base de datos centralizada PostgreSQL para múltiples inquilinos (multi-tenant).
- **Local**: Base de datos ligera SQLite rápida para operaciones offline en el punto de venta.

## 2. Identificadores Globales (GlobalID)

Para evitar colisiones de IDs secuenciales entre diferentes locales offline, cada entidad principal (`Product`, `Sale`, `ThirdParty`, `Invoice`) incluye un campo `global_id` (UUID). 
- Los registros se rastrean y vinculan mediante este `global_id` en lugar del `id` incremental numérico.

## 3. Metadatos de Sincronización

Cada tabla sincronizable cuenta con tres campos clave:
- `global_id`: Identificador único universal.
- `sync_status`: Estado del registro (`PENDING`, `SYNCED`, `FAILED`).
- `last_synced_at`: Marca de tiempo de la última sincronización exitosa.

## 4. Flujo de Sincronización Bidireccional

### A. Pull (Online -> Local)
*Utilizado para: Productos, Precios, Clientes, Configuración.*

1. El cliente local solicita `/sync/pull?last_sync=TIMESTAMP`.
2. El servidor online busca registros creados o modificados después de ese TIMESTAMP.
3. El cliente local recibe el lote e inserta/actualiza los registros en su SQLite local.
4. **Resolución de Conflictos**: La versión online siempre prevalece (Master-Slave).

### B. Push (Local -> Online)
*Utilizado para: Ventas realizadas, Movimientos de inventario, Facturas locales.*

1. El cliente local recopila todos los registros con `sync_status = 'PENDING'`.
2. Envía un lote JSON al endpoint `/sync/push`.
3. El servidor online procesa el lote usando `ON CONFLICT (global_id) DO UPDATE` para asegurar idempotencia (si se reintenta un envío, no se duplican las ventas).
4. El servidor responde con éxito y el cliente local marca los registros como `SYNCED`.

## 5. Endpoints de Sincronización

- **GET `/sync/pull`**: Descarga actualizaciones desde la nube. Soporta filtros por fecha y empresa.
- **POST `/sync/push`**: Sube transacciones locales a la nube. Maneja lotes masivos para eficiencia.

## 6. Generación del Ejecutable Offline

Para generar la versión ejecutable para Windows (offline), utiliza el script `build_offline.bat` incluido en la raíz.

### Uso del Ejecutable
Para ejecutar el POS en modo local/offline:
1. Asegúrate de tener el archivo `.env` configurado o define las variables:
   ```cmd
   set DATABASE_DRIVER=sqlite
   set DATABASE_URL=aura_pos.db
   set PORT=8081
   aura-pos-offline.exe
   ```
### Descarga del Ejecutable
Una vez que el servidor online esté en ejecución, puedes descargar el ejecutable directamente desde el siguiente endpoint:
- **URL**: `GET /download/offline-pos`
- **Descripción**: Inicia la descarga del archivo `aura-pos-offline.exe` generado automáticamente al iniciar el servidor.

---

## 7. Guía de Pruebas (Local vs Online)

Para probar la interacción, puedes simular ambos entornos en tu máquina:

### Paso 1: Levantar Versión Online (PostgreSQL)
1. Asegúrate de que `DATABASE_DRIVER=postgres` en tu `.env`.
2. Ejecuta `go run cmd/api/main.go` (Puerto 8081).
3. Crea un producto o cliente usando los endpoints estándar.
4. Verifica que el registro tenga un `global_id` (se genera automáticamente).

### Paso 2: Levantar Versión Local (SQLite)
1. Abre una nueva terminal.
2. Define variables diferentes para no chocar:
   ```cmd
   set PORT=8082
   set DATABASE_DRIVER=sqlite
   set DATABASE_URL=pos_local.db
   go run cmd/api/main.go
   ```
3. El POS local iniciará en el puerto 8082 usando SQLite.

### Paso 3: Probar Pull (Online -> Local)
Desde la terminal o Insomnia/Postman, llama al endpoint del POS local:
`GET http://localhost:8082/sync/pull`
- El POS descargará los productos creados en el Paso 1.

### Paso 4: Probar Push (Local -> Online)
1. Crea una venta o cliente directamente en el POS local (Puerto 8082).
2. Asegúrate de que el registro local tenga `sync_status = 'PENDING'`.
3. Envía el lote al POS online:
   ```bash
   # Simulación de envío del lote local al online
   POST http://localhost:8081/sync/push
   ```
   (En una implementación real, el cliente offline tiene un botón de "Sincronizar" que empaqueta los datos locales y los envía al servidor online).

---

### Ejemplo de Procesamiento de Lote
```go
// modules/sync/service.go
func (s *service) Push(ctx context.Context, batch *SyncBatch) (*SyncStats, error) {
    // Procesa productos, terceros y ventas de forma atómica
    // Cada inserción en el servidor online valida el GlobalID
}
```
