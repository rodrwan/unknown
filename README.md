# Crypto Dashboard

Este es un proyecto de dashboard para monitoreo de criptomonedas desarrollado en Go. El proyecto utiliza una arquitectura moderna y sigue las mejores prácticas de desarrollo.

## 🚀 Características

- Dashboard para monitoreo de criptomonedas
- Desarrollado en Go
- Interfaz web moderna usando Templ
- Integración con Buda.com para datos de criptomonedas
- Base de datos SQL con SQLx
- Contenedorización con Docker

## 🛠️ Tecnologías Principales

- Go 1.24.0
- Templ (Framework de templates)
- SQLx para manejo de base de datos
- Buda.com API para datos de criptomonedas
- Docker y Docker Compose

## 📋 Prerrequisitos

- Go 1.24.0 o superior
- Docker y Docker Compose (opcional)

## 🔧 Instalación

1. Clonar el repositorio:
```bash
git clone https://github.com/rodrwan/unknown.git
cd unknown
```

2. Copiar el archivo de variables de entorno:
```bash
cp .env.example .env
```

3. Configurar las variables de entorno en el archivo `.env`

4. Ejecutar con Docker:
```bash
docker-compose up
```

O ejecutar localmente:
```bash
make run
```

## 🏗️ Estructura del Proyecto

```
.
├── cmd/server/     # Punto de entrada de la aplicación
├── internal/       # Código interno de la aplicación
├── static/         # Archivos estáticos
├── Dockerfile      # Configuración de Docker
├── docker-compose.yml
├── go.mod          # Dependencias de Go
└── Makefile        # Comandos de utilidad
```

## 📝 Licencia

Este proyecto está bajo la Licencia MIT - ver el archivo LICENSE para más detalles. 