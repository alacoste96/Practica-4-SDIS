# Práctica 4 SDIS - Documentación del Proyecto

### 1. Documento Principal PDF
- **`Practica_4_alacoste_SSOO_dist.pdf`**: Memoria técnica completa del proyecto

### 2. Diagramas UML en formato Draw.io
Todos los diagramas están en formato `.drawio` editable:

1. **`diagrama_arquitectura.drawio`**
   - Arquitectura general del sistema
   - Muestra componentes: Servidor, Mutua, Taller
   - Visualiza conexiones TCP y flujo de datos

2. **`diagrama_secuencia_completo.drawio`**
   - Diagrama de secuencia UML completo
   - Flujo de un coche por las 4 fases del taller
   - Interacciones: Productor → Workers → Logger

3. **`diagrama_secuencia_cliente_servidor.drawio`**
   - Comunicación cliente-servidor
   - Protocolo TCP entre Mutua, Servidor y Taller
   - Broadcasting de comandos

4. **`diagrama_flujo_prioridades.drawio`**
   - Algoritmo de gestión de prioridades
   - Función `getCarFromQ()`
   - Manejo de modos restrictivos (1-3) y de prioridad (4-6)

5. **`diagrama_flujo_worker.drawio`**
   - Comportamiento de un worker (goroutine)
   - Ciclo de vida: recibir → procesar → enviar
   - Gestión de fase de entrega


## Estructura del Documento PDF

La memoria técnica incluye:

### 1. Introducción

### 2. Arquitectura del Sistema
- Estructura general
- Componentes principales
- Diagrama de arquitectura

### 3. Descripción Detallada de Componentes
- Servidor (servidor.go)
- Cliente Taller (taller.go)
- Gestión de concurrencia
- Fases del proceso (4 fases)
- Sistema de prioridades (categorías A/B/C)
- Modos de operación (0-9)

### 4. Diagramas UML
- Diagrama de secuencia: Flujo completo
- Diagrama de secuencia: Cliente-Servidor
- Diagrama de flujo: Prioridades
- Diagrama de flujo: Worker

### 5. Implementación de Concurrencia
- Goroutines del sistema (37 concurrentes)
- Sincronización con RWMutex
- Atomic Int32 para estado
- Canal bufferizado como semáforo
- Patrón Select para priorización

### 6. Tests y Resultados
- Estrategia de testing
- 6 configuraciones de test
- Resultados de tiempos de ejecución
- Análisis del impacto de categorías
- Análisis de plazas vs mecánicos
- Cobertura de código: 47.4%
- Race conditions: Ninguna detectada

### 7. Conclusiones
- Decisiones de diseño
- Lecciones aprendidas

### 8. Referencias

---

## Detalles Importantes del Código

### Modos de Operación Explicados

**Modos 1-3 (Restrictivos)**:
- Solo afectan a la **entrada** de coches nuevos
- Los coches ya dentro **siguen siendo procesados normalmente**
- Ejemplo Modo 1: Solo entran coches categoría A

**Modos 4-6 (Prioridades)**:
- Afectan a **todos los coches**, incluidos los que ya están dentro
- Cambian dinámicamente el **orden de atención en todas las fases**
- Ejemplo Modo 4: Se priorizan coches A en todas las fases

**Modos 0 y 9**:
- **No entran** coches nuevos
- Los coches dentro **siguen siendo gestionados**

**Modos 7 y 8**:
- Mantienen el **modo anterior**
- No modifican el estado

### Tests

- **Alcance**: Solo prueba el módulo `cliente` (taller.go y archivos relacionados)
- **No incluyen**: `mutua.go` ni `servidor.go`
- **Modo evaluado**: Todos los tests usan modo 4 (Prioridad categoría A)
- **Métricas**: Robustez, tiempos de ejecución, race conditions, cobertura
- **Reescalado temporal**: 1 segundo real = 100ms en test (multiplicar × 10)

### Resultados de Tests (Tiempos Reales)

| Test | Config | Tiempo Medio | Por Coche |
|------|--------|--------------|-----------|
| T1 | A10 B10 C10, P6 M3 | 107.31s | 3.58s |
| T2 | A20 B5 C5, P6 M3 | 133.55s | 4.45s |
| T3 | A5 B5 C20, P6 M3 | 82.97s | 2.77s |
| T4 | A10 B10 C10, P4 M4 | 153.47s | 5.12s |
| T5 | A20 B5 C5, P4 M4 | 191.88s | 6.40s |
| T6 | A5 B5 C20, P4 M4 | 116.05s | 3.87s |

**Observaciones clave**:
- Mayor proporción de coches A → mayor tiempo total
- 4 plazas + 4 mecánicos es más lento que 6 plazas + 3 mecánicos
- El cuello de botella varía según la configuración

## Información del Proyecto

- **Asignatura**: Sistemas Distribuidos
- **Grado**: Ingeniería Telemática
- **Práctica**: 4 - Concurrencia en GO - Servidores
- **Lenguaje**: Go
- **Tema**: Sistema de gestión de taller con concurrencia

---
