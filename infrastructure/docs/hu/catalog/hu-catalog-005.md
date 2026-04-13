# HU-CATALOG-005: Eliminar Producto

## 📌 Información General
- ID: HU-CATALOG-005
- Epic: EPICA-016 (Gestión de Catálogo)
- Prioridad: Alta
- Estado: ✅ Completado
- Porcentaje: 100%
- Autor: Equipo Backend Aura POS
- Fecha: 2025-01-15

---

## 👤 Historia de Usuario

**Como** administrador  
**Quiero** eliminar un producto del catálogo  
**Para** quitar productos que ya no se venden

---

## 🧠 Descripción Funcional

El sistema debe realizar un soft delete del producto, cambiando su estado a DELETED en lugar de eliminarlo físicamente de la base de datos.

---

## ✅ Criterios de Aceptación

### Escenario 1: Eliminación exitosa (soft delete)
- Dado que existe un producto con ID 5
- Cuando solicito DELETE a `/products/5`
- Entonces el sistema cambia el status a "DELETED"
- Y retorna código 204 (sin contenido)

### Escenario 2: Producto no encontrado
- Dado que no existe el producto
- Cuando solicito DELETE
- Entonces retorna error 404

---

## ❌ Casos de Error

- **Producto no encontrado**: Retorna 404

---

## 🔐 Reglas de Negocio

- La eliminación es un soft delete (cambia status a DELETED)
- El producto no aparecerá en listados pero existirá en BD
- Solo usuarios ADMIN/MANAGER pueden eliminar

---

## 📡 Requisitos Técnicos

- **Endpoint**: `/products/:id`
- **Método HTTP**: DELETE
- **Response**: 204 No Content

---

## 📎 Dependencias

- Servicios: products.Service.Delete