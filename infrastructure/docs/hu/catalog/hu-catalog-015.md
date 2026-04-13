# HU-CATALOG-015: Eliminar Marca

## 📌 Información General
- ID: HU-CATALOG-015
- Epic: EPICA-016 (Gestión de Catálogo)
- Prioridad: Media
- Estado: ✅ Completado
- Porcentaje: 100%
- Autor: Equipo Backend Aura POS
- Fecha: 2025-01-15

---

## 👤 Historia de Usuario

**Como** administrador  
**Quiero** eliminar una marca  
**Para** remover marcas que ya no se usan

---

## 🧠 Descripción Funcional

El sistema debe permitir eliminar una marca solo si no hay productos asociados.

---

## ✅ Criterios de Aceptación

### Escenario 1: Eliminación exitosa
- Dado que la marca no tiene productos asociados
- Cuando solicito DELETE a `/brands/:id`
-Entonces elimina la marca

### Escenario 2: Productos asociados
- Dado que hay productos usando esa marca
- Cuando intento eliminar
-Entonces retorna error 400

---

## 📡 Requisitos Técnicos

- **Endpoint**: `/brands/:id`
- **Método HTTP**: DELETE

---

## 📎 Dependencias

- Servicios: brands.Service.Delete