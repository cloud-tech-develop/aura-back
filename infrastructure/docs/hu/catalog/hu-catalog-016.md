# HU-CATALOG-016: Crear Unidad

## 📌 Información General
- ID: HU-CATALOG-016
- Epic: EPICA-016 (Gestión de Catálogo)
- Prioridad: Alta
- Estado: ✅ Completado
- Porcentaje: 100%
- Autor: Equipo Backend Aura POS
- Fecha: 2025-01-15

---

## 👤 Historia de Usuario

**Como** administrador  
**Quiero** crear una unidad de medida  
**Para** estandarizar las cantidades en los productos

---

## 🧠 Descripción Funcional

El sistema debe permitir crear una unidad de medida (ej: unidad, kilogramo, litro, caja) con su nombre y/abreviatura. Por defecto, las unidades se crean activas y permitiendo decimales.

---

## ✅ Criterios de Aceptación

### Escenario 1: Creación exitosa
- Dado que proporciono nombre y abreviatura válidos
- Cuando solicito POST a `/catalog/units`
- Entonces crea la unidad y retorna 201

### Escenario 2: Nombre duplicado
- Dado que ya existe una unidad con el mismo nombre
- Cuando intento crear
- Entonces retorna error 400

### Escenario 3: Datos faltantes
- Dado que no proporciono nombre o abreviatura
- Cuando intento crear
- Entonces retorna error 400

---

## ❌ Casos de Error

- Nombre duplicado → 400 Bad Request
- Nombre vacío → 400 Bad Request
- Abreviatura vacía → 400 Bad Request

---

## 📡 Requisitos Técnicos

- **Endpoint**: `/catalog/units`
- **Método HTTP**: POST
- **Request**:
  ```json
  {
    "name": "Kilogramo",
    "abbreviation": "kg",
    "active": true,
    "allow_decimals": true
  }
  ```
- **Response**: Unit (201 Created)

---

## 📎 Dependencias

- Servicios: units.Service.Create