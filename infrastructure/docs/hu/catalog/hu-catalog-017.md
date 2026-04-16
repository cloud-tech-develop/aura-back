# HU-CATALOG-017: Listar Unidades

## 📌 Información General
- ID: HU-CATALOG-017
- Epic: EPICA-016 (Gestión de Catálogo)
- Prioridad: Alta
- Estado: ✅ Completado
- Porcentaje: 100%
- Autor: Equipo Backend Aura POS
- Fecha: 2025-01-15

---

## 👤 Historia de Usuario

**Como** administrador  
**Quiero** listar todas las unidades de medida  
**Para** ver las opciones disponibles

---

## 🧠 Descripción Funcional

El sistema debe retornar todas las unidades de medida activas e inactivas de la empresa, ordenadas por nombre.

---

## ✅ Criterios de Aceptación

### Escenario 1: Listado exitoso
- Dado que hay unidades registradas
- Cuando solicito GET a `/catalog/units`
- Entonces retorna todas las unidades

### Escenario 2: Sin unidades
- Dado que no hay unidades registradas
- Cuando solicito GET a `/catalog/units`
- ThenReturns empty array

---

## 📡 Requisitos Técnicos

- **Endpoint**: `/catalog/units`
- **Método HTTP**: GET
- **Response**: Array de Unit

---

## 📎 Dependencias

- Servicios: units.Service.List