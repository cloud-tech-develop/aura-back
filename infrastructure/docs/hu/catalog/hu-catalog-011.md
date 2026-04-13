# HU-CATALOG-011: Crear Marca

## 📌 Información General
- ID: HU-CATALOG-011
- Epic: EPICA-016 (Gestión de Catálogo)
- Prioridad: Alta
- Estado: ✅ Completado
- Porcentaje: 100%
- Autor: Equipo Backend Aura POS
- Fecha: 2025-01-15

---

## 👤 Historia de Usuario

**Como** administrador/gerente  
**Quiero** crear marcas para identificar fabricantes  
**Para** clasificar productos por marca

---

## 🧠 Descripción Funcional

El sistema debe permitir crear marcas verificando nombre único por empresa.

---

## ✅ Criterios de Aceptación

### Escenario 1: Crear marca exitosamente
- Dado que no existe una marca con ese nombre
- Cuando envío POST a `/brands`
-Entonces crea la marca con estado ACTIVE

---

## ❌ Casos de Error

- **Nombre duplicado**: Retorna 400

---

## 📡 Requisitos Técnicos

- **Endpoint**: `/brands`
- **Método HTTP**: POST

---

## 📎 Dependencias

- Servicios: brands.Service.Create