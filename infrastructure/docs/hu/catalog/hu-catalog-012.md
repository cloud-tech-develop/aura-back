# HU-CATALOG-012: Listar Marcas

## 📌 Información General
- ID: HU-CATALOG-012
- Epic: EPICA-016 (Gestión de Catálogo)
- Prioridad: Alta
- Estado: ✅ Completado
- Porcentaje: 100%
- Autor: Equipo Backend Aura POS
- Fecha: 2025-01-15

---

## 👤 Historia de Usuario

**Como** usuario del sistema  
**Quiero** listar las marcas disponibles  
**Para** seleccionar marcas al crear/filtrar productos

---

## 🧠 Descripción Funcional

El sistema debe retornar todas las marcas activas de la empresa.

---

## 📡 Requisitos Técnicos

- **Endpoint**: `/brands`
- **Método HTTP**: GET

---

## 📎 Dependencias

- Servicios: brands.Service.List