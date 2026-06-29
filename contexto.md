# Contexto do Projeto — Observabilidade em Go

## Diretrizes para o assistente

- Este é um projeto **didático e guiado** — o objetivo é aprender, não entregar um produto
- **Não escreva código** a menos que o usuário peça explicitamente
- Sempre aponte **onde encontrar** como fazer (docs, artigos, pkg.go.dev)
- Explique o **porquê** antes do como — o usuário quer entender, não só copiar
- Avance **uma etapa por vez**, confirmando que a anterior está concluída antes de seguir
- Use linguagem direta, sem enrolação
- Quando o usuário concluir uma etapa, marque como `[x]` no checklist abaixo

---

## Stack do projeto

| Ferramenta | Papel |
|---|---|
| **Go** | Linguagem base |
| **Docker** | Containerização da aplicação |
| **docker-compose** | Orquestração local dos serviços |
| **Zap** | Geração de logs estruturados (já tem familiaridade) |
| **Loki** | Armazenamento e captura de logs |
| **Prometheus** | Métricas |
| **Grafana** | Visualização (dashboards) |
| **OpenTelemetry** | Instrumentação de traces (padrão aberto) |
| **Tempo** | Armazenamento de traces distribuídos |

---

## Ordem de estudo acordada

1. ~~Base do projeto Go~~ ✅
2. ~~Prometheus + Grafana~~ ✅
3. ~~Loki + Alloy~~ ✅
4. **OpenTelemetry + Tempo** ← estamos aqui
5. Feature real com banco de dados (PostgreSQL) — **introduzir o DB aqui, junto com a primeira feature que precisar persistir dados**
6. cAdvisor — **introduzir quando quiser métricas USE de containers (CPU, memória, rede por container). Não exige mudança de código — funciona com qualquer linguagem**
7. Nginx — **introduzir apenas quando precisar de reverse proxy ou SSL; desnecessário em desenvolvimento local**

---

## Checklist — Etapa 1: Base do projeto Go ✅

- [x] `go mod init` — inicializar o módulo
- [x] Criar estrutura de pastas (`cmd/`, `internal/`, etc.)
- [x] Servidor HTTP com Gin subindo na porta `8080`
- [x] Endpoint `/health` retornando `{ "status": "ok" }`
- [x] Configuração via variáveis de ambiente (`envconfig` + `godotenv`)
- [x] `Dockerfile` — multi-stage build
- [x] `docker-compose.yml` — orquestrar app + serviços
- [x] `.dockerignore` — excluir arquivos desnecessários do build
- [x] Air — hot reload local
- [x] Zap configurado no startup com log `"server started"` e campo `port`

## Checklist — Etapa 2: Prometheus + Grafana ✅

- [x] Prometheus e Grafana no `docker-compose.yml`
- [x] `prometheus.yml` — configuração de scrape
- [x] Endpoint `/metrics` exposto na app
- [x] Middleware de métricas HTTP (`http_requests_total`, `http_request_duration_seconds`)
- [x] Grafana conectado ao Prometheus como datasource
- [x] Dashboard com painéis RED (taxa, erros 5xx, latência P95)
- [x] Dashboard com painéis USE Go (goroutines, memória heap, GC duration)

## Checklist — Etapa 3: Loki + Alloy ← em andamento

- [x] Loki e Alloy no `docker-compose.yml`
- [x] Configuração do Alloy para coletar logs
- [x] Zap escrevendo logs em arquivo JSON
- [x] Loki como datasource no Grafana
- [x] Visualizar logs no Grafana com LogQL

---

## Estrutura de pastas definida

```
observabilidade/
├── cmd/
│   └── api/
│       └── main.go          ← entrypoint
├── internal/
│   ├── handler/             ← handlers HTTP
│   ├── service/             ← lógica de negócio
│   └── config/              ← leitura de env vars
├── go.mod
└── go.sum
```

---

## Decisões técnicas tomadas

| Decisão | Escolha | Motivo |
|---|---|---|
| Router HTTP | **Gin** | Framework maduro e amplamente usado no ecossistema Go |
| Config | **envconfig** | Simples, sem overhead. Viper é overkill para esse projeto |
| Logger | **Zap** | Usuário já tem familiaridade |
| Logger dev | modo `development` | Output legível/colorido para rodar local |
| Logger prod | modo `production` | Output JSON para integrar com Loki |
| Container | **Docker** | Isola a app e garante paridade entre local e produção |
| Orquestração local | **docker-compose** | Sobe app + Prometheus + Grafana + Loki + Tempo juntos |

---

## Referências por etapa

### Base Go
- Módulos: https://go.dev/doc/modules/gomod-ref
- Estrutura de projeto: https://github.com/golang-standards/project-layout
- Gin: https://gin-gonic.com/docs/
- envconfig: https://github.com/kelseyhightower/envconfig
- Zap: https://pkg.go.dev/go.uber.org/zap

### Docker
- Dockerfile reference: https://docs.docker.com/reference/dockerfile/
- docker-compose reference: https://docs.docker.com/compose/compose-file/
- Multi-stage builds: https://docs.docker.com/build/building/multi-stage/

### Prometheus
- Docs oficiais: https://prometheus.io/docs/introduction/overview/
- Client Go: https://github.com/prometheus/client_golang
- PromQL cheat sheet: https://promlabs.com/promql-cheat-sheet/

### Loki + Alloy ← próxima etapa
- Loki docs: https://grafana.com/docs/loki/latest/
- Alloy docs: https://grafana.com/docs/alloy/latest/
- LogQL basics: https://grafana.com/docs/loki/latest/query/

### OpenTelemetry
- Docs Go: https://opentelemetry.io/docs/languages/go/

---

## Contexto do usuário

- Tem familiaridade com **Zap**
- Nunca usou **Loki, Prometheus, Grafana, OpenTelemetry ou Tempo**
- Quer aprender o ecossistema de observabilidade do zero, de forma estruturada
- Prefere entender o conceito antes de escrever código
