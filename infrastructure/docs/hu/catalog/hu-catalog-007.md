# HU-CATALOG-007: Listar Categorías

## 📌 Información General
- ID: HU-CATALOG-007
- Epic: EPICA-016 (Gestión de Catálogo)
- Prioridad: Alta
- Estado: ✅ Completado
- Porcentaje: 100%
- Autor: Equipo Backend Aura POS
- Fecha: 2025-01-15

---

## 👤 Historia de Usuario

**Como** usuario del sistema  
**Quiero** listar las categorías disponibles  
**Para** seleccionar categorías al crear/filtrar productos

---

## 🧠 Descripción Funcional

El sistema debe retornar todas las categorías activas de la empresa ordenadas por nombre.

---

## ✅ Criterios de Aceptación

### Escenario 1: Listar categorías
- Dado que la empresa tiene categorías
- Cuando solicito GET a `/categories`
-Entonces retorna lista de categorías activas

---

## 📡 Requisitos Técnicos

- **Endpoint**: `/categories`
- **Método HTTP**: GET
- **Response**: Array de categorías

---

## 📎 Dependencias

- Servicios: categories.Service.List