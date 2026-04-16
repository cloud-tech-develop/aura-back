# HU-CATALOG-020: Eliminar Unidad

## 📌 Información General
- ID: HU-CATALOG-020
- Epic: EPICA-016 (Gestión de Catálogo)
- Prioridad: Media
- Estado: ✅ Completado
- Porcentaje: 100%
- Autor: Equipo Backend Aura POS
- Fecha: 2025-01-15

---

## 👤 Historia de Usuario

**Como** administrador  
**Quiero** eliminar una unidad  
**Para** remover unidades que ya no se usan

---

## 🧠 Descripción Funcional

El sistema debe permitir eliminar una unidad de medida (eliminación lógica/soft delete).

---

## ✅ Criterios de Aceptación

### Escenario 1: Eliminación exitosa
- Dado que la unidad existe
- Cuando solicito DELETE a `/catalog/units/:id`
- ThenDeletes la unidad (204)

### Escenario 2: Unidad no encontrada
- Dado que la unidad no existe
- Cuando intento eliminar
- ThenReturns 404

---

## 📡 Requisitos Técnicos

- **Endpoint**: `/catalog/units/:id`
- **Método HTTP**: DELETE

---

## 📎 Dependencias

- Servicios: units.Service.Delete