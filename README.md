# Distributed Container Infrastructure | Infraestructura distribuida de contenedores

**Repository name suggestion:** `container-infrastructure-go`  
**Description:** Three virtual machines running Go REST APIs across containerd, Podman and Docker, with a Zot OCI registry.

## English

### Overview

An operating-systems lab project that deploys three Go HTTP APIs across three Linux virtual machines. Each VM uses a different container runtime, and the services check the health of their peer APIs over the virtual network. A Zot registry on the third VM stores the container images.

### Architecture

| VM | Address used in the lab | Runtime | Workload |
|---|---|---|---|
| VM1 | `192.168.122.210` | containerd + nerdctl | API1 (`8081`), API2 (`8082`) |
| VM2 | `192.168.122.56` | Podman | API3 (`8083`) |
| VM3 | `192.168.122.126` | Docker | Zot OCI registry (`5000`) |

The API services provide `/health` endpoints and routes to check connectivity to the other APIs. Responses use JSON. The Go services include Dockerfiles with multi-stage builds.

### Repository branches

- `VM1`: API1 and API2 source code.
- `VM2`: API3 source code.
- `Documentacion`: installation guide, technical manual and deployment screenshots.
- `Proyecto2`: a separate kernel monitoring and container management extension with a Go daemon, Linux kernel module, eBPF program, Docker Compose services and scripts.

The default branch is `VM1`; it does not contain the complete three-VM project by itself. Review the relevant branches to see each part.

### Technology

- Go 1.22
- Linux virtual machines with KVM/QEMU networking
- containerd and nerdctl, Podman, Docker
- Zot OCI container registry
- REST over HTTP, JSON
- Dockerfiles and container image workflows
- In the `Proyecto2` branch: C kernel module, Go, eBPF, Docker Compose, Valkey and Grafana

### Documentation

- [Installation and configuration guide](https://github.com/Deltai-cod12/202404856_LAB_SO1_2S2026/blob/Documentacion/Guia_Instalacion_202404856.md)
- [Technical manual (PDF)](https://github.com/Deltai-cod12/202404856_LAB_SO1_2S2026/tree/Documentacion)

### API routes

Each API exposes `GET /health`. Cross-service checks include routes such as:

```text
GET /api1/202404856/call-api2
GET /api1/202404856/call-api3
GET /api2/202404856/call-api1
GET /api2/202404856/call-api3
```

API3 has corresponding routes for checking API1 and API2. These addresses are configured for the lab network; update them before running the services in another environment.

### Run locally

The repository is organized by VM branch, so first check out the branch containing the API you want to run. For example, from the API directory:

```bash
go run .
```

Use the supplied Dockerfile to build a container image. A full deployment also requires the three VMs, their network connectivity, the matching runtime on each VM, and the Zot registry. Follow the installation guide for the lab setup.

### Academic context

Developed for the Operating Systems 1 course at Universidad de San Carlos de Guatemala.

## Español

### Descripción

Proyecto de laboratorio de sistemas operativos que despliega tres API HTTP desarrolladas en Go en tres máquinas virtuales Linux. Cada VM utiliza un runtime de contenedores distinto. Las API comprueban la disponibilidad de sus servicios pares mediante la red virtual y un registro Zot en la tercera VM almacena las imágenes.

### Arquitectura

| VM | Dirección usada en el laboratorio | Runtime | Servicios |
|---|---|---|---|
| VM1 | `192.168.122.210` | containerd + nerdctl | API1 (`8081`), API2 (`8082`) |
| VM2 | `192.168.122.56` | Podman | API3 (`8083`) |
| VM3 | `192.168.122.126` | Docker | Registro OCI Zot (`5000`) |

Las API exponen rutas `/health` y rutas para comprobar la conectividad con las otras API. Las respuestas se envían en JSON. Los servicios Go incluyen Dockerfiles con compilación multi-etapa.

### Ramas del repositorio

- `VM1`: código de API1 y API2.
- `VM2`: código de API3.
- `Documentacion`: guía de instalación, manual técnico y capturas del despliegue.
- `Proyecto2`: extensión separada para supervisión del kernel y gestión de contenedores, con daemon en Go, módulo de kernel Linux, programa eBPF, servicios Docker Compose y scripts.

La rama predeterminada es `VM1`; por sí sola no contiene el proyecto completo de las tres VM. Revisa las ramas correspondientes para consultar cada componente.

### Tecnologías

- Go 1.22
- Máquinas virtuales Linux y red KVM/QEMU
- containerd y nerdctl, Podman, Docker
- Registro OCI Zot
- REST sobre HTTP y JSON
- Dockerfiles y construcción de imágenes
- En la rama `Proyecto2`: módulo de kernel en C, Go, eBPF, Docker Compose, Valkey y Grafana

### Documentación

- [Guía de instalación y configuración](https://github.com/Deltai-cod12/202404856_LAB_SO1_2S2026/blob/Documentacion/Guia_Instalacion_202404856.md)
- [Manual técnico y capturas](https://github.com/Deltai-cod12/202404856_LAB_SO1_2S2026/tree/Documentacion)

### Rutas de las API

Cada API expone `GET /health`. Algunas rutas de comprobación entre servicios son:

```text
GET /api1/202404856/call-api2
GET /api1/202404856/call-api3
GET /api2/202404856/call-api1
GET /api2/202404856/call-api3
```

API3 cuenta con rutas equivalentes para comprobar API1 y API2. Las direcciones están configuradas para la red del laboratorio; cámbialas antes de ejecutar los servicios en otro entorno.

### Ejecución local

El código está separado por ramas según la VM. Cambia primero a la rama que contiene la API que quieras ejecutar. Desde el directorio de esa API:

```bash
go run .
```

También puedes construir una imagen con el Dockerfile incluido. El despliegue completo requiere las tres VM, conectividad entre ellas, el runtime correspondiente en cada VM y el registro Zot. Consulta la guía de instalación para reproducir el laboratorio.

### Contexto académico

Desarrollado para el curso de Sistemas Operativos 1 de la Universidad de San Carlos de Guatemala.

---

**Topics:** `go`, `golang`, `linux`, `containers`, `containerd`, `nerdctl`, `podman`, `docker`, `oci-registry`, `zot`, `rest-api`, `virtual-machines`
