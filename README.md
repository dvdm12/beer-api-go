<h1 align="center">
  <i class="fa-solid fa-beer-mug-empty"></i>
  Beer API - Microservicios en Go
</h1>

<p align="center">
  Sistema de gestión de cervezas basado en microservicios, desarrollado en Go (Gin) y MongoDB,
  orquestado con Docker y Docker Compose.
</p>

<p align="center">
  <img src="https://img.shields.io/badge/Go-1.x-00ADD8?logo=go&logoColor=white" alt="Go Badge">
  <img src="https://img.shields.io/badge/MongoDB-Database-47A248?logo=mongodb&logoColor=white" alt="MongoDB Badge">
  <img src="https://img.shields.io/badge/Docker-Compose-2496ED?logo=docker&logoColor=white" alt="Docker Badge">
  <img src="https://img.shields.io/badge/Architecture-Microservices-blueviolet" alt="Microservices Badge">
</p>

<hr/>

<h2>
  <i class="fa-solid fa-circle-info"></i>
  Descripción del proyecto
</h2>

<p>
  <strong>Beer API</strong> es un sistema de microservicios diseñado para gestionar un catálogo de cervezas.
  Cada microservicio se encarga de una responsabilidad específica (crear, actualizar, etc.) y se comunica con
  una base de datos <strong>MongoDB</strong> centralizada.
</p>

<ul>
  <li><i class="fa-solid fa-microchip"></i> Backend implementado en <strong>Go</strong> usando el framework <strong>Gin</strong>.</li>
  <li><i class="fa-solid fa-database"></i> Persistencia de datos en <strong>MongoDB</strong>.</li>
  <li><i class="fa-brands fa-docker"></i> Contenedores gestionados mediante <strong>Docker</strong> y <strong>Docker Compose</strong>.</li>
  <li><i class="fa-solid fa-cubes"></i> Arquitectura basada en <strong>microservicios</strong>.</li>
</ul>

<hr/>

<h2>
  <i class="fa-solid fa-users"></i>
  Integrantes del proyecto
</h2>

<ul>
  <li><strong>DAVID MANTILLA AVILES</strong></li>
  <li><strong>SANTIAGO PALACIO VASQUEZ</strong></li>
</ul>

<hr/>

<h2>
  <i class="fa-solid fa-diagram-project"></i>
  Arquitectura de microservicios
</h2>

<p>Servicios principales (hasta el momento):</p>

<ul>
  <li>
    <strong>MongoDB</strong> (<code>mongodb</code>):
    <ul>
      <li>Base de datos principal del sistema.</li>
      <li>Se levanta como un contenedor independiente.</li>
    </ul>
  </li>
  <li>
    <strong>Create Service</strong> (<code>create-service</code>):
    <ul>
      <li>Microservicio encargado de <strong>crear</strong> nuevas cervezas.</li>
      <li>Expone endpoints REST con Go + Gin.</li>
      <li>Se conecta a MongoDB para almacenar la información.</li>
    </ul>
  </li>
  <li>
    <strong>Update Service</strong> (<code>update-service</code>):
    <ul>
      <li>Microservicio encargado de <strong>actualizar</strong> la información de las cervezas.</li>
      <li>Compartirá la misma base de datos MongoDB.</li>
    </ul>
  </li>
</ul>

<hr/>

<h2>
  <i class="fa-solid fa-folder-tree"></i>
  Estructura del proyecto
</h2>

<pre>
beer-api-go/
├── create-service/          # Microservicio de creación
│   ├── cmd/
│   │   └── main.go          # Punto de entrada del servicio
│   ├── internal/
│   │   ├── controllers/     # Controladores HTTP (Gin handlers)
│   │   ├── db/              # Adaptadores de base de datos
│   │   ├── models/          # Modelos y estructuras de datos
│   │   ├── repository/      # Capa de repositorio
│   │   └── services/        # Lógica de negocio
│   ├── Dockerfile
│   └── go.mod
├── update-service/          # Microservicio de actualización
│   ├── cmd/
│   │   └── main.go          # Punto de entrada del servicio
│   ├── internal/
│   │   ├── controllers/     # Controladores HTTP (Gin handlers)
│   │   ├── db/              # Adaptadores de base de datos
│   │   ├── models/          # Modelos y estructuras de datos
│   │   ├── repository/      # Capa de repositorio
│   │   └── services/        # Lógica de negocio
│   ├── Dockerfile
│   └── go.mod
├── docker-compose.yml       # Orquestación de contenedores
├── .env                     # Variables de entorno
└── README.md
</pre>

<hr/>

<h2>
  <i class="fa-solid fa-gears"></i>
  Requisitos previos
</h2>

<ul>
  <li><i class="fa-brands fa-docker"></i> <strong>Docker</strong> instalado.</li>
  <li><i class="fa-brands fa-docker"></i> <strong>Docker Compose</strong> instalado.</li>
  <li><i class="fa-brands fa-golang"></i> <strong>Go 1.x</strong> (opcional para desarrollo local sin contenedores).</li>
  <li><i class="fa-solid fa-key"></i> Acceso al repositorio (Bitbucket / GitHub) por SSH o HTTPS.</li>
</ul>

<hr/>

<h2>
  <i class="fa-solid fa-file-lines"></i>
  Variables de entorno
</h2>

<p>Archivo <code>.env</code> (ubicado al mismo nivel que <code>docker-compose.yml</code>):</p>

<pre>
MONGO_URI=mongodb://mongodb:27017
DATABASE=beersdb
COLLECTION=beers
</pre>

<p>
  Estas variables son utilizadas por los microservicios para conectarse a la base de datos MongoDB.
  <strong>Nota:</strong> El hostname debe coincidir con el nombre del contenedor de MongoDB definido en <code>docker-compose.yml</code> (<code>mongodb</code>).
</p>

<hr/>

<h2>
  <i class="fa-solid fa-play"></i>
  Cómo ejecutar el proyecto
</h2>

<ol>
  <li>
    <strong>Clonar el repositorio</strong>
    <pre>git clone git@bitbucket.org:team-development1/beer-api.git
cd beer-api</pre>
  </li>
  <li>
    <strong>Crear el archivo <code>.env</code></strong> (si no existe) en la raíz del proyecto:
    <pre>
MONGO_URI=mongodb://mongodb:27017
DATABASE=beers
COLLECTION=beers
    </pre>
  </li>
  <li>
    <strong>Levantar los contenedores con Docker Compose</strong>
    <pre>docker compose up -d</pre>
    <p>Esto iniciará MongoDB y los microservicios (por ejemplo, <code>create-service</code> y <code>update-service</code>).</p>
  </li>
  <li>
    <strong>Verificar contenedores en ejecución</strong>
    <pre>docker compose ps</pre>
    <p>Deberías ver tres contenedores activos:</p>
    <ul>
      <li><code>mongodb</code> - Base de datos MongoDB (Puerto 27017)</li>
      <li><code>create_service</code> - Servicio de creación (Puerto 8080)</li>
      <li><code>update_service</code> - Servicio de actualización (Puerto 8082)</li>
    </ul>
  </li>
  <li>
    <strong>Verificar logs de los servicios</strong>
    <pre>docker logs create_service
docker logs update_service</pre>
  </li>
</ol>

<hr/>

<h2>
  <i class="fa-solid fa-plug"></i>
  Endpoints (ejemplo)
</h2>

<p><strong>Create Service</strong> (Puerto 8080):</p>

<ul>
  <li>
    <code>POST /beers</code> – Crear una nueva cerveza.
    <pre>
curl -X POST http://localhost:8080/beers \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Heineken",
    "brand": "Heineken",
    "alcohol": 5.0,
    "year": 2024
  }'
    </pre>
  </li>
</ul>

<p><strong>Update Service</strong> (Puerto 8082):</p>

<ul>
  <li>
    <code>PUT /beers/:id</code> – Actualizar una cerveza existente.
    <pre>
curl -X PUT http://localhost:8082/beers/BEER_ID \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Heineken Premium",
    "brand": "Heineken International",
    "alcohol": 5.2,
    "year": 2025
  }'
    </pre>
  </li>
</ul>

<p><strong>Modelo de datos Beer:</strong></p>
<pre>
{
  "name": string,      // Nombre de la cerveza
  "brand": string,     // Marca de la cerveza
  "alcohol": float,    // Contenido de alcohol
  "year": int          // Año
}
</pre>

<hr/>

<h2>
  <i class="fa-solid fa-code-branch"></i>
  Flujo de trabajo con Git
</h2>

<ol>
  <li>Crear rama de trabajo desde <code>develop</code>:</li>
</ol>

<pre>git checkout develop
git pull origin develop
git checkout -b feature/nueva-funcionalidad</pre>

<ol start="2">
  <li>Realizar cambios, commitear y pushear:</li>
</ol>

<pre>
git add .
git commit -m "Descripción clara de los cambios"
git push origin feature/nueva-funcionalidad
</pre>

<ol start="3">
  <li>Crear Pull Request hacia <code>develop</code> en Bitbucket.</li>
</ol>

<hr/>

<h2>
  <i class="fa-solid fa-clipboard-check"></i>
  Estado actual
</h2>

<ul>
  <li>✅ Microservicio <strong>create-service</strong> (Puerto 8080).</li>
  <li>✅ Microservicio <strong>update-service</strong> (Puerto 8082).</li>
  <li>✅ <code>docker-compose.yml</code> configurado para orquestar servicios y base de datos.</li>
  <li>✅ Tests unitarios implementados para ambos servicios.</li>
  <li>✅ Arquitectura estandarizada con adaptadores, repositorios, servicios y controladores.</li>
</ul>

<h3>
  <i class="fa-solid fa-flask"></i>
  Ejecutar tests
</h3>

<p><strong>Create Service:</strong></p>
<pre>cd create-service && go test ./... -v</pre>

<p><strong>Update Service:</strong></p>
<pre>cd update-service && go test ./... -v</pre>

<hr/>

<h2>
  <i class="fa-solid fa-handshake-angle"></i>
  Contribuciones
</h2>

<p>
  Las contribuciones están enfocadas a:
</p>
<ul>
  <li>Agregar nuevos microservicios (read, delete, auth, etc.).</li>
  <li>Mejorar la estructura interna del código Go.</li>
  <li>Agregar pruebas unitarias e integración.</li>
  <li>Mejorar la documentación y ejemplos de uso.</li>
</ul>

<p>
  <em>Proyecto desarrollado por <strong>DAVID MANTILLA AVILES</strong> y <strong>SANTIAGO PALACIO VASQUEZ</strong>.</em>
</p>
