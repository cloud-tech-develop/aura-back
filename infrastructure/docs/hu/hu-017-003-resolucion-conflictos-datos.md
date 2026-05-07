# HU-017-003 - Resolucion de Conflictos de Datos

## 📌 Informacion General
- ID: HU-017-003
- Epic: EPICA-017
- Prioridad: Media
- Estado: Backlog
- Porcentaje: 0%
- Autor: QA Team
- Fecha: 2026-05-04

---

## 👤 Historia de Usuario

**Como** sistema de sincronizacion  
**Quiero** detectar y resolver conflictos cuando el mismo producto es modificado tanto en online como en offline  
**Para** garantizar la integridad de datos y evitar perdida de cambios realizados por operadores en cualquier version

---

## 🧠 Descripcion Funcional

Cuando un cliente offline intenta sincronizar cambios hacia online, pero el producto ya fue modificado en linea desde la ultima sincronizacion, se genera un conflicto. El sistema debe detectarlo y aplicar una estrategia de resolucion automatica basada en timestamps, pudiendo marcar para resolución manual en casos ambiguos.

Los conflictos pueden ocurrir en:
1. Actualizacion offline de producto ya modificado en online
2. Eliminacion offline de producto ya eliminado en online
3. Eliminacion offline de producto modificado en online

---

## ✅ Criterios de Aceptacion

### Escenario 1: Conflicto de actualizacion - online gana
- Dado que el producto fue modificado en offline con timestamp T1
- Y el producto fue modificado en online con timestamp T2 (T2 > T1)
- Cuando el cliente offline sincroniza los cambios
- Entonces el sistema detecta el conflicto
- Y aplica la estrategia "online wins" (el cambio online se preserva)
- Y el cambio offline se descarta silenciosamente
- Y se registra en logs para auditoria

### Escenario 2: Conflicto de actualizacion con offline mas reciente
- Dado que el producto fue modificado en offline con timestamp T2
- Y el producto fue modificado en online con timestamp T1 (T1 > T2)
- Cuando el cliente offline sincroniza los cambios
- Entonces el sistema aplica los cambios offline (sobreescribe online)
- Y se registra en logs para auditoria

### Escenario 3: Producto eliminado en online pero modificado en offline
- Dado que el producto fue eliminado en online (soft delete)
- Y el cliente offline intenta sincronizar cambios
- Entonces el sistema restaura el producto en online con los datos de offline
- Y se registra el evento como restauracion

### Escenario 4: Producto modificado en online pero eliminado en offline
- Dado que el producto fue modificado en online
- Y el cliente offline intenta eliminarlo
-Entonces el sistema marca el producto como eliminado en online
- Y se registra el evento como eliminacion

### Escenario 5: Conflicto ambiguous requiere resolucion manual
- Dado que los cambios offline y online afectan campos diferentes
- Cuando se detecta el conflicto
- Entonces el sistema aplica merge automatique (combina campos)
- No requiere intervencion manual

### Escenario 6: Registro de conflictos para auditoria
- Dado que ocurre un conflicto que se resolvio automaticamente
- Cuando se completa la sincronizacion
- Entonces se registra en una tabla de log de conflictos
- El administrador puede consultar el historial

---

## ❌ Casos de Error

- Si el timestamp es invalido o ausente: Usar created_at como fallback
- Si no se puede determinar el winner: Por defecto gana online
- Si la base de datos no responde durante la verificacion: Sincronizar sin verificacion (modo unsafe)

---

## 🔐 Reglas de Negocio

- La estrategia de resolucion se configura por tenant (por defecto: online wins)
- El timestamp de referencia es updated_at del producto
- Si updated_at es NULL, usar created_at
- Se permite configurar "offline wins" o "manual" por empresa
- El historial de conflictos se mantiene por 30 dias
- Solo se registran conflictos reales (no merge automatico)

---

## 🎨 Consideraciones UI/UX

No aplica - es logica de backend.

---

## 📡 Requisitos Tecnicos

### Estructura de Conflictos

```go
type SyncConflict struct {
    ProductID    int64     `json:"product_id"`
    SKU         string    `json:"sku"`
    OfflineAt   time.Time `json:"offline_at"`
    OnlineAt    time.Time `json:"online_at"`
    Strategy    string    `json:"strategy"` // "online_wins", "offline_wins", "manual"
    Resolution  string    `json:"resolution"`
    Details     string    `json:"details"`
}
```

### Tabla de Auditoria

```sql
CREATE TABLE sync_conflicts (
    id SERIAL PRIMARY KEY,
    tenant_slug VARCHAR(100) NOT NULL,
    product_id BIGINT NOT NULL,
    offline_data JSONB,
    online_data JSONB,
    strategy VARCHAR(20) NOT NULL,
    resolution VARCHAR(20) NOT NULL,
    resolved_at TIMESTAMP,
    resolved_by VARCHAR(50),
    created_at TIMESTAMP DEFAULT NOW()
);
```

### Logica de Resolucion

```
function resolveConflict(offlineProduct, onlineProduct, strategy):
    if strategy == "offline_wins":
        return offlineProduct // offline sobreescribe online
    else if strategy == "online_wins" or default:
        if onlineProduct.updated_at > offlineProduct.updated_at:
            return onlineProduct // online wins
        else:
            return offlineProduct // offline es mas reciente
```

### Metadata en eventos

```json
{
  "event": "product.updated",
  "conflict_resolution": {
    "detected": true,
    "strategy": "online_wins",
    "winner": "online",
    "offline_at": "2026-05-04T10:30:00Z",
    "online_at": "2026-05-04T10:35:00Z"
  }
}
```

---

## 🧪 Criterios de Testing

- Tests unitarios para cada escenario de conflicto
- Tests de integracion con timestamps frontera
- Tests de configuracion de estrategias por tenant
- Tests de tabla de auditoria

---

## 📎 Dependencias

- modules/catalog/products/domain.go: Entidad Product existente
- modules/catalog/products/repository.go: Metodos CRUD existentes
- Tabla sync_conflicts: Nueva migracion

---

## 🚫 Fuera de Alcance

- Resolucion manual via interfaz de usuario
- Notificaciones en tiempo real al admin
- sincronizacion de otras entidades
- Configuracion de estrategias por campo

---

## 🧠 Generacion de Codigo

Requerir:
- Nueva migracion para tabla sync_conflicts
- Metodo en repository para verificar timestamps
- Logica de resolucion en service de offline
- Registro de conflictos en tabla
- Tests unitarios por escenario