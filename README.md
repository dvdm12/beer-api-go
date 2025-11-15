<!DOCTYPE html>
<html lang="es">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Beer API - Microservicios en Go</title>

    <!-- Iconos CDN -->
    <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/font-awesome/6.5.0/css/all.min.css">

    <style>
        body {
            font-family: Arial, Helvetica, sans-serif;
            margin: 40px;
            line-height: 1.6;
        }
        h1, h2 {
            color: #333;
        }
        code, pre {
            background: #f3f3f3;
            padding: 8px;
            border-radius: 6px;
            display: block;
        }
        ul {
            margin-left: 20px;
        }
    </style>
</head>

<body>

<h1 align="center">
    <i class="fa-solid fa-beer-mug-empty"></i>
    Beer API – Microservicios en Go
</h1>

<p align="center">
    Sistema de microservicios para la gestión de cervezas, desarrollado con Go (Gin), MongoDB y Docker.
</p>

<p align="center">
    <img src="https://img.shields.io/badge/Go-1.x-00ADD8?logo=go&logoColor=white">
    <img src="https://img.shields.io/badge/MongoDB-Database-47A248?logo=mongodb&logoColor=white">
    <img src="https://img.shields.io/badge/Docker-Compose-2496ED?logo=docker&logoColor=white">
    <img src="https://img.shields.io/badge/Microservices-Architecture-blueviolet">
</p>

<hr>

<h2><i class="fa-solid fa-circle-info"></i> Descripción del proyecto</h2>

<p>
    <strong>Beer API</strong> es un sistema diseñado bajo arquitectura de microservicios. Cada servicio cumple una función específica:
</p>

<ul>
    <li><strong>Create Service:</strong> creación de nuevas cervezas.</li>
    <li><strong>Update Service:</strong> actualización de cervezas existentes.</li>
    <li><strong>MongoDB:</strong> base de datos centralizada para todos los servicios.</li>
</ul>

<p>Todo está orquestado mediante Docker Compose.</p>

<hr>

<h2><i class="fa-solid fa-users"></i> Integrantes del proyecto</h2>

<ul>
    <li><strong>DAVID MANTILLA AVILÉS</strong></li>
    <li><strong>SANTIAGO PALACIO VÁSQUEZ</strong></li>
</ul>

<hr>

<h2><i class="fa-solid fa-diagram-project"></i> Arquitectura de microservicios</h2>

<ul>
    <li>
        <strong>MongoDB</strong>
        <ul>
            <li>Base principal del sistema</li>
            <li>Contenedor independiente</li>
        </ul>
    </li>

    <li>
        <strong>Create Service</strong>
        <ul>
            <li>API REST para crear cervezas</li>
            <li>Go + Gin</li>
        </ul>
    </li>

    <li>
        <strong>Update Service</strong>
        <ul>
            <li>API REST para actualizar cervezas</li>
            <li>Conexión al mismo MongoDB</li>
        </ul>
    </li>
</ul>

<hr>

<h2><i class="fa-solid fa-folder-tree"></i> Estructura del proyecto</h2>

<pre>
parcial_three_go/
├── cmd/
├── internal/
│   ├── controllers/
│   ├── models/
│   ├── routes/
│   └── services/
├── create-service/
├── update-service/
├── docker-compose.yml
└── .env
</pre>

<hr>

<h2><i class="fa-solid fa-gears"></i> Requisitos</h2>

<ul>
    <li>Docker</li>
    <li>Docker Compose</li>
    <li>Go 1.x (para desarrollo local)</li>
</ul>

<hr>

<h2><i class="fa-solid fa-file-lines"></i> Variables de entorno</h2>

<p>Archivo <code>.env</code>:</p>

<pre>
MONGO_URI=mongodb://mongodb:27017
DATABASE=beersdb
COLLECTION=beers
</pre>

<hr>

<h2><i class="fa-solid fa-play"></i> Cómo ejecutar</h2>

<ol>
    <li>
        <strong>Clonar el repositorio</strong>
        <pre>git clone git@bitbucket.org:team-development1/beer-api.git</pre>
    </li>

    <li>
        <strong>Crear archivo .env</strong>
        <pre>
MONGO_URI=mongodb://mongodb:27017
DATABASE=beersdb
COLLECTION=beers
        </pre>
    </li>

    <li>
        <strong>Levantar los contenedores</strong>
        <pre>docker compose up -d</pre>
    </li>

    <li>
        <strong>Verificar estado</strong>
        <pre>docker compose ps</pre>
    </li>
</ol>

<hr>

<h2><i class="fa-solid fa-plug"></i> Endpoints principales</h2>

<p><strong>Create Service:</strong></p>

<ul>
    <li><code>POST /beers</code> – crear cerveza</li>
    <li><code>GET /health</code> – healthcheck</li>
</ul>

<p><strong>Update Service:</strong></p>

<ul>
    <li><code>PUT /beers/:id</code> – actualizar cerveza</li>
    <li><code>GET /health</code></li>
</ul>

<hr>

<h2><i class="fa-solid fa-code-branch"></i> Flujo de trabajo Git</h2>

<pre>
git checkout develop
git pull origin develop
git checkout -b feature/nueva-funcionalidad

git add .
git commit -m "mensaje claro"
git push origin feature/nueva-funcionalidad
</pre>

<hr>

<h2><i class="fa-solid fa-handshake-angle"></i> Contribuciones</h2>

<p>
    Aportamos expandiendo los microservicios, mejorando el código en Go y optimizando documentación y arquitectura.
</p>

<p><strong>Proyecto desarrollado por DAVID MANTILLA AVILÉS y SANTIAGO PALACIO VÁSQUEZ.</strong></p>

</body>
</html>
