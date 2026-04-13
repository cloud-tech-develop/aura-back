# HU-CATALOG-004: Actualizar Producto

## 📌 Información General
- ID: HU-CATALOG-004
- Epic: EPICA-016 (Gestión de Catálogo)
- Prioridad: Alta
- Estado: ✅ Completado
- Porcentaje: 100%
- Autor: Equipo Backend Aura POS
- Fecha: 2025-01-15

---

## 👤 Historia de Usuario

**Como** administrador/gerente  
**Quiero** actualizar la información de un producto  
**Para** mantener el catálogo actualizado con precios, información correcta

---

## 🧠 Descripción Funcional

El sistema debe permitir actualizar todos los campos de un producto existente, manteniendo la validación de negocio.

---

## ✅ Criterios de Aceptación

### Escenario 1: Actualización exitosa
- Dado que existe un producto con ID 5
- Cuando envío PUT a `/products/5` con nuevos valores
-Entonces el sistema actualiza el producto y retorna el producto actualizado

### Escenario 2: Actualizar solo algunos campos
- Dado que existe un producto
- Cuando envío solo el campo `name`
-Entonces el sistema actualiza solo ese campo, manteniendo los demás

### Escenario 3: Cambiar estado del producto
- Dado que existe un producto activo
- Cuando actualizo el status a "INACTIVE"
-Entonces el producto queda inactivo y no aparece en listados

---

## ❌ Casos de Error

- **SKU duplicado**: Si se intenta usar un SKU que ya existe en la empresa
- **Precio inválido**: Si `sale_price < cost_price`
- **Producto no encontrado**: Retorna 404

---

## 🔐 Reglas de Negocio

- Todos los campos son opcionales en la actualización
- El SKU solo se valida si se incluye y es diferente al actual
- Solo usuarios ADMIN/MANAGER pueden actualizar

---

## 📡 Requisitos Técnicos

- **Endpoint**: `/products/:id`
- **Método HTTP**: PUT
- **Request**: JSON con campos a actualizar
- **Response**: 200 OK con producto actualizado

---

## 📎 Dependencias

- Servicios: products.Service.Update