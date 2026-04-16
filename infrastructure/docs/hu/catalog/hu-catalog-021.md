# HU-CATALOG-021: Listar Unidades Paginadas

## 📌 Información General
- ID: HU-CATALOG-021
- Epic: EPICA-016 (Gestión de Catálogo)
- Prioridad: Alta
- Estado: ✅ Completado
- Porcentaje: 100%
- Autor: Equipo Backend Aura POS
- Fecha: 2025-01-15

---

## 👤 Historia de Usuario

**Como** administrador  
**Quiero** listar las unidades paginadas  
**Para** navegar entre múltiples unidades

---

## 🧠 Descripción Funcional

El sistema debe permitir paginar las unidades de medida con filtros opcionales de búsqueda, ordenamiento y límites.

---

## ✅ Criterios de Aceptación

### Escenario 1: Paginación por defecto
- Dado que tengo más de 10 unidades
- Cuando solicito POST a `/catalog/units/page` sin parámetros
-Entonces retorna las primeras 10 unidades

### Escenario 2: Paginación con búsqueda
- Dado que hay unidades que coinciden con la búsqueda
- Cuando busco por nombre
- ThenReturns las unidades filtradas

### Escenario 3: Paginación vacía
- Dado que no hay unidades
- WhenSolicito paginación
- ThenReturns página vacía

---

## 📡 Requisitos Técnicos

- **Endpoint**: `/catalog/units/page`
- **Método HTTP**: POST
- **Request**:
  ```json
  {
    "page": 1,
    "limit": 10,
    "search": "",
    "sort": "id",
    "order": "asc",
    "params": {}
  }
  ```
- **Response**: PageResult conarraydeUnit
- **Campos filtros**: search, sort, order

---

## 📎 Dependencias

- Servicios: units.Service.Page

---

## 🧪 Criterios de Testing

- Verificar paginación por defecto (page=1, limit=10)
- Verificar filtros de búsqueda
- Verificar ordenamiento (sort, order)
- Verificar página vacía