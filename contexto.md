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

1. **Base do projeto Go** ← estamos aqui
2. Prometheus + Grafana
3. Loki + Promtail
4. OpenTelemetry + Tempo
5. Feature real com banco de dados (PostgreSQL) — **introduzir o DB aqui, junto com a primeira feature que precisar persistir dados**
6. Nginx — **introduzir apenas quando precisar de reverse proxy ou SSL; desnecessário em desenvolvimento local**

---

## Checklist — Etapa 1: Base do projeto Go

- [ ] `go mod init` — inicializar o módulo
- [ ] Criar estrutura de pastas (`cmd/`, `internal/`, etc.)
- [ ] Servidor HTTP com Gin subindo na porta `8080`
- [ ] Endpoint `/health` retornando `{ "status": "ok" }`
- [ ] Configuração via variáveis de ambiente (`envconfig`)
- [ ] `Dockerfile` — buildar e rodar a app em container
- [ ] `docker-compose.yml` — orquestrar app + serviços futuros
- [ ] `.dockerignore` — excluir arquivos desnecessários do build
- [ ] Zap configurado no startup com log `"server started"` e campo `port`

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

### Prometheus (próxima etapa)
- Docs oficiais: https://prometheus.io/docs/introduction/overview/
- Client Go: https://github.com/prometheus/client_golang

### Loki
- Docs oficiais: https://grafana.com/docs/loki/latest/

### OpenTelemetry
- Docs Go: https://opentelemetry.io/docs/languages/go/

---

## Contexto do usuário

- Tem familiaridade com **Zap**
- Nunca usou **Loki, Prometheus, Grafana, OpenTelemetry ou Tempo**
- Quer aprender o ecossistema de observabilidade do zero, de forma estruturada
- Prefere entender o conceito antes de escrever código
