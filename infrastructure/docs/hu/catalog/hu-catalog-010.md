# HU-CATALOG-010: Eliminar Categoría

## 📌 Información General
- ID: HU-CATALOG-010
- Epic: EPICA-016 (Gestión de Catálogo)
- Prioridad: Media
- Estado: ✅ Completado
- Porcentaje: 100%
- Autor: Equipo Backend Aura POS
- Fecha: 2025-01-15

---

## 👤 Historia de Usuario

**Como** administrador  
**Quiero** eliminar una categoría  
**Para** remover categorías que ya no se usan

---

## 🧠 Descripción Funcional

El sistema debe permitir eliminar una categoría solo si no hay productos asociados.

---

## ✅ Criterios de Aceptación

### Escenario 1: Eliminación exitosa
- Dado que la categoría no tiene productos asociados
- Cuando solicito DELETE a `/categories/:id`
-Entonces elimina la categoría

### Escenario 2: Productos asociados
- Dado que hay productos usando esa categoría
- Cuando intento eliminar
-Entonces retorna error 400

---

## ❌ Casos de Error

- **Productos asociados**: No se puede eliminar si hay productos
- **No encontrada**: Retorna 404

---

## 📡 Requisitos Técnicos

- **Endpoint**: `/categories/:id`
- **Método HTTP**: DELETE

---

## 📎 Dependencias

- Servicios: categories.Service.Delete