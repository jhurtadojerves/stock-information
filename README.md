# 📈 Stock Insight – Proyecto Fullstack

Este repositorio contiene una aplicación completa para visualizar recomendaciones de acciones, compuesta por un backend en Go y un frontend en Astro + Vue + Tailwind.

---

## 🧱 Requisitos

- [Docker + Docker Compose](https://www.docker.com/)
- [Node.js](https://nodejs.org/) (recomendado: v18+)
- npm

---

## 🚀 Cómo levantar el entorno local

### 🔹 1. Backend (Go + CockroachDB)

#### 📁 Ir a la carpeta del backend:

```bash
cd backend
```

#### ⚙️ Variables de entorno

Debes renombrar el archivo .env.example a .env y completar las variables de entorno

#### 🐳 Levantar los servicios con Docker:

```bash
make up
```

---

### 🔹 2. Frontend (Astro + Vue)

#### 📁 Ir a la carpeta del frontend:

```bash
cd frontend
```

#### 📦 Instalar dependencias:

```bash
npm install
```

#### 🔥 Levantar el servidor de desarrollo:

```bash
npm run dev
```
