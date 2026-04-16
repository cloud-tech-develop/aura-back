# HU-CATALOG-019: Actualizar Unidad

## 📌 Información General
- ID: HU-CATALOG-019
- Epic: EPICA-016 (Gestión de Catálogo)
- Prioridad: Media
- Estado: ✅ Completado
- Porcentaje: 100%
- Autor: Equipo Backend Aura POS
- Fecha: 2025-01-15

---

## 👤 Historia de Usuario

**Como** administrador  
**Quiero** actualizar una unidad  
**Para** corregir o modificar sus datos

---

## 🧠 Descripción Funcional

El sistema debe permitir modificar el nombre, abreviatura, estado activo y permisos de decimales de una unidad existente.

---

## ✅ Criterios de Aceptación

### Escenario 1: Actualización exitosa
- Dado que la unidad existe
- Cuando solicito PUT a `/catalog/units/:id`
- ThenReturns la unidad actualizada

### Escenario 2: Unidad no encontrada
- Dado que la unidad no existe
- Cuando intento actualizar
- ThenReturns 404

---

## 📡 Requisitos Técnicos

- **Endpoint**: `/catalog/units/:id`
- **Método HTTP**: PUT
- **Request**:
  ```json
  {
    "name": "Kilogramo",
    "abbreviation": "kg",
    "active": true,
    "allow_decimals": false
  }
  ```

---

## 📎 Dependencias

- Servicios: units.Service.Update