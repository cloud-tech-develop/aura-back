# HU-CATALOG-009: Actualizar Categoría

## 📌 Información General
- ID: HU-CATALOG-009
- Epic: EPICA-016 (Gestión de Catálogo)
- Prioridad: Alta
- Estado: ✅ Completado
- Porcentaje: 100%
- Autor: Equipo Backend Aura POS
- Fecha: 2025-01-15

---

## 👤 Historia de Usuario

**Como** administrador/gerente  
**Quiero** actualizar una categoría  
**Para** corregir o modificar información de la categoría

---

## 🧠 Descripción Funcional

El sistema debe permitir actualizar el nombre y descripción de una categoría.

---

## ✅ Criterios de Aceptación

### Escenario 1: Actualización exitosa
- Dado que existe una categoría
- Cuando envío PUT a `/categories/:id`
-Entonces actualiza la categoría

---

## ❌ Casos de Error

- **Nombre duplicado**: Retorna 400
- **No encontrada**: Retorna 404

---

## 📡 Requisitos Técnicos

- **Endpoint**: `/categories/:id`
- **Método HTTP**: PUT

---

## 📎 Dependencias

- Servicios: categories.Service.Update