# HU-CATALOG-022: Verificar Existencia de Producto por SKU

## 📌 Información General
- ID: HU-CATALOG-022
- Epic: EPICA-016 (Gestión de Catálogo)
- Prioridad: Media
- Estado: ✅ Completado
- Porcentaje: 100%
- Autor: Equipo Backend Aura POS
- Fecha: 2026-04-19

---

## 👤 Historia de Usuario

**Como** usuario del sistema (admin, manager, cajero)  
**Quiero** verificar si existe un producto en mi catálogo mediante su SKU  
**Para** poder validar rápidamente si puedo usar ese código antes de crear un nuevo producto

---

## 🧠 Descripción Funcional

El sistema debe permitir verificar la existencia de un producto mediante su código SKU, retornando la información básica del producto si existe o indicando que no existe.

---

## ✅ Criterios de Aceptación

### Escenario 1: SKU existe
- Dado que existe un producto con SKU "154-44" en la empresa
- Cuando solicito GET a `/catalog/products/exist/154-44`
- Entonces el sistema retorna HTTP 200 con `exists: true` y los datos del producto

### Escenario 2: SKU no existe
- Dado que NO existe un producto con SKU "999-99" en la empresa
- Cuando solicito GET a `/catalog/products/exist/999-99`
- Entonces el sistema retorna HTTP 200 con `exists: false`

### Escenario 3: Sin autenticación
- Dado que no proporciono token de autenticación
- Cuando solicito GET a `/catalog/products/exist/154-44`
- Entonces el sistema retorna HTTP 401

---

## 📡 Contrato API

### Endpoint
```
GET /catalog/products/exist/{sku}
```

### Headers
```
Authorization: Bearer <token>
Content-Type: application/json
```

### Respuesta - SKU Existe (200 OK)
```json
{
  "success": true,
  "data": {
    "exists": true,
    "sku": "154-44",
    "product": {
      "id": 1,
      "name": "Producto 1",
      "sku": "154-44",
      "barcode": "123456789"
    }
  },
  "message": "Operacion exitosa"
}
```

### Respuesta - SKU No Existe (200 OK)
```json
{
  "success": true,
  "data": {
    "exists": false,
    "sku": "999-99"
  },
  "message": "Operacion exitosa"
}
```

### Respuesta - No Autorizado (401)
```json
{
  "success": false,
  "data": null,
  "message": "token no proporcionado o inválido"
}
```

---

## 🔄 Flujo Técnico

1. El handler recibe el SKU desde el path parameter
2. Se obtiene el `enterprise_id` del token JWT
3. El servicio busca el producto por SKU y enterprise_id
4. Si existe, retorna los datos del producto
5. Si no existe, retorna `exists: false`

---

## 📋 Reglas de Negocio

- El SKU debe ser único dentro de la empresa
- La búsqueda es sensible a mayúsculas/minúsculas
- Solo retorna productos activos (no eliminados lógicamente)
