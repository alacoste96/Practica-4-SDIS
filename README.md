# Práctica 4 SDIS - Documentación del Proyecto

## 📋 Archivos Generados

### 1. Documento Principal LaTeX
- **`Practica_4_Memoria_Tecnica.tex`**: Memoria técnica completa del proyecto

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

### 3. Análisis de Tests
- **`test_analysis.json`**: Resultados procesados de los tests

---

## 🚀 Pasos para Completar la Documentación

### Paso 1: Editar Diagramas en Draw.io

1. Ve a https://app.diagrams.net (o https://draw.io)
2. Abre cada archivo `.drawio`:
   - File → Open from... → Device
   - Selecciona el archivo correspondiente
3. Edita el diagrama según tus necesidades:
   - Ajusta posiciones de elementos
   - Modifica textos y etiquetas
   - Cambia colores si lo deseas
   - Añade o elimina componentes

### Paso 2: Exportar Diagramas a PDF

Para cada diagrama editado:

1. En Draw.io: **File → Export as → PDF**
2. Configuración recomendada:
   - ✅ Include a copy of my diagram (para poder reabrirlo)
   - ✅ All Pages
   - Calidad: 100%
3. Guarda con los nombres exactos:
   - `diagrama_arquitectura.pdf`
   - `diagrama_secuencia_completo.pdf`
   - `diagrama_secuencia_cliente_servidor.pdf`
   - `diagrama_flujo_prioridades.pdf`
   - `diagrama_flujo_worker.pdf`

### Paso 3: Colocar PDFs en la Carpeta del Proyecto

Coloca todos los PDFs exportados en la **misma carpeta** que el archivo `.tex`

### Paso 4: Compilar el Documento LaTeX

```bash
# Primera compilación (genera referencias)
pdflatex Practica_4_Memoria_Tecnica.tex

# Segunda compilación (actualiza índice y referencias)
pdflatex Practica_4_Memoria_Tecnica.tex
```

**Resultado**: `Practica_4_Memoria_Tecnica.pdf`

---

## 📊 Estructura del Documento LaTeX

La memoria técnica incluye:

### 1. Introducción
- Objetivos del proyecto
- Contexto y requisitos

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
- Race conditions: ✅ Ninguna detectada

### 7. Conclusiones
- Logros principales
- Decisiones de diseño
- Lecciones aprendidas
- Mejoras futuras

### 8. Referencias

---

## 🎯 Detalles Importantes del Código

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

---

## 🛠️ Requisitos del Sistema

### Para Editar Diagramas
- Navegador web moderno (Chrome, Firefox, Edge, Safari)
- Conexión a internet (para https://app.diagrams.net)

### Para Compilar LaTeX
- Distribución LaTeX:
  - **Windows**: MiKTeX o TeX Live
  - **macOS**: MacTeX
  - **Linux**: TeX Live
- Paquetes requeridos (normalmente incluidos):
  - `babel`, `graphicx`, `hyperref`, `listings`
  - `xcolor`, `geometry`, `float`, `longtable`, `booktabs`

### Alternativa Online para LaTeX
Si no tienes LaTeX instalado, usa **Overleaf**:
1. Ve a https://www.overleaf.com
2. Crea una cuenta gratuita
3. Crea un nuevo proyecto
4. Sube el archivo `.tex` y los PDFs de los diagramas
5. Compila directamente en el navegador

---

## 📁 Estructura Recomendada del Proyecto

```
Practica4/
├── doc/
│   ├── 4_practica_ssdd_dist.pdf
│   ├── Practica_4_Memoria_Tecnica.tex      ← Documento LaTeX
│   ├── Practica_4_Memoria_Tecnica.pdf      ← PDF generado
│   ├── diagrama_arquitectura.drawio
│   ├── diagrama_arquitectura.pdf
│   ├── diagrama_secuencia_completo.drawio
│   ├── diagrama_secuencia_completo.pdf
│   ├── diagrama_secuencia_cliente_servidor.drawio
│   ├── diagrama_secuencia_cliente_servidor.pdf
│   ├── diagrama_flujo_prioridades.drawio
│   ├── diagrama_flujo_prioridades.pdf
│   ├── diagrama_flujo_worker.drawio
│   └── diagrama_flujo_worker.pdf
├── README.md
├── src/
│   ├── cliente/
│   │   ├── goroutines.go
│   │   ├── mutex.go
│   │   ├── taller.go
│   │   ├── taller_test.go
│   │   ├── types.go
│   │   └── utility.go
│   ├── go.mod
│   ├── mutua/
│   │   └── mutua.go
│   └── servidor/
│       └── servidor.go
└── tests/
    ├── cover.out
    ├── test_cover.txt
    ├── test_race.txt
    └── test_report.txt
```

---

## 💡 Consejos para los Diagramas

### Diagrama de Arquitectura
- Ajusta el tamaño de las cajas según la cantidad de texto
- Verifica que las flechas de conexión no se superpongan
- Usa colores consistentes para cada tipo de componente

### Diagramas de Secuencia
- Asegúrate de que las lifelines estén verticalmente alineadas
- Las flechas deben ser claras y no cruzarse innecesariamente
- Los mensajes deben estar en orden cronológico de arriba a abajo

### Diagramas de Flujo
- Verifica que todas las decisiones (rombos) tengan 2+ salidas
- Asegúrate de que el flujo sea fácil de seguir
- Usa colores para diferenciar tipos de operaciones

---

## ✅ Checklist Final

Antes de entregar, verifica:

- [ ] Todos los diagramas .drawio están editados y finalizados
- [ ] Todos los diagramas están exportados a PDF
- [ ] Los PDFs tienen los nombres exactos especificados
- [ ] El documento LaTeX compila sin errores
- [ ] Todas las figuras aparecen correctamente en el PDF
- [ ] El índice está completo y correcto
- [ ] Los datos de los tests coinciden con tus resultados
- [ ] El enlace al repositorio GitHub está actualizado
- [ ] Tu nombre de usuario aparece en el documento
- [ ] El documento final está en formato PDF

---

## 📞 Notas Adicionales

### Repositorio GitHub
URL actual: https://github.com/alacoste96/Practica-4-SDIS

### Formato de Entrega
Según el enunciado, el archivo debe llamarse:
- `Practica_4_[tu_nombre_usuario]_SSOO_dist.pdf`

### Contenido Mínimo según Enunciado
✅ Descripción técnica con diagramas UML
✅ Diagramas de secuencia
✅ Código fuente o link a repositorio
✅ Métricas de tests
✅ Comparativas de rendimiento

---

## 🎓 Información del Proyecto

- **Asignatura**: Sistemas Distribuidos
- **Grado**: Ingeniería Telemática
- **Práctica**: 4 - Concurrencia en GO - Servidores
- **Lenguaje**: Go
- **Tema**: Sistema de gestión de taller con concurrencia

---

**¡Mucha suerte con tu entrega! 🚀**

Si encuentras algún problema al compilar o editar, revisa:
1. Que todos los archivos estén en la carpeta correcta
2. Que los nombres de archivo coincidan exactamente
3. Que tu instalación de LaTeX esté completa
4. Los logs de compilación para identificar errores específicos
