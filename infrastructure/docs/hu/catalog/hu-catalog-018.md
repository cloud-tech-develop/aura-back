# HU-CATALOG-018: Obtener Unidad por ID

## 📌 Información General
- ID: HU-CATALOG-018
- Epic: EPICA-016 (Gestión de Catálogo)
- Prioridad: Alta
- Estado: ✅ Completado
- Porcentaje: 100%
- Autor: Equipo Backend Aura POS
- Fecha: 2025-01-15

---

## 👤 Historia de Usuario

**Como** administrador  
**Quiero** obtener una unidad por su ID  
**Para** ver los detalles de una unidad específica

---

## 🧠 Descripción Funcional

El sistema debe retornar los datos de una unidad de medida específica por su ID.

---

## ✅ Criterios de Aceptación

### Escenario 1: Obtención exitosa
- Dado que la unidad existe
- Cuando solicito GET a `/catalog/units/:id`
- ThenReturns los datos de la unidad

### Escenario 2: Unidad no encontrada
- Dado que la unidad no existe
- Cuando solicito GET a `/catalog/units/:id`
-Entonces retorna 404

---

## 📡 Requisitos Técnicos

- **Endpoint**: `/catalog/units/:id`
- **Método HTTP**: GET
- **Response**: Unit

---

## 📎 Dependencias

- Servicios: units.Service.GetByID