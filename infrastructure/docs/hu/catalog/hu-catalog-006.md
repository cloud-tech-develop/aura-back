# HU-CATALOG-006: Crear Categoría

## 📌 Información General
- ID: HU-CATALOG-006
- Epic: EPICA-016 (Gestión de Catálogo)
- Prioridad: Alta
- Estado: ✅ Completado
- Porcentaje: 100%
- Autor: Equipo Backend Aura POS
- Fecha: 2025-01-15

---

## 👤 Historia de Usuario

**Como** administrador/gerente  
**Quiero** crear categorías para organizar productos  
**Para** facilitar la navegación y búsqueda de productos

---

## 🧠 Descripción Funcional

El sistema debe permitir crear categorías con nombre y descripción, verificando que el nombre no esté duplicado dentro de la empresa.

---

## ✅ Criterios de Aceptación

### Escenario 1: Crear categoría exitosamente
- Dado que no existe una categoría con ese nombre
- Cuando envío POST a `/categories` con nombre válido
-Entonces el sistema crea la categoría con estado ACTIVE

### Escenario 2: Nombre duplicado
- Dado que ya existe una categoría con ese nombre
- Cuando intento crear
-Entonces retorna error 400

---

## ❌ Casos de Error

- **Nombre duplicado**: Retorna 400
- **Nombre requerido**: Si no se proporciona nombre, retorna 400

---

## 📡 Requisitos Técnicos

- **Endpoint**: `/categories`
- **Método HTTP**: POST
- **Request**: `{ "name": "Bebidas", "description": "..." }`
- **Response**: 201 Created

---

## 📎 Dependencias

- Servicios: categories.Service.Create