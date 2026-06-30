# Observability

Projeto de observabilidade utilizando **Zap**, **Loki**, **Alloy**, **Prometheus**, **OpenTelemetry**, **Tempo** e **Grafana**.

---

## Ambientes

| Ambiente | Comando | O que sobe |
|---|---|---|
| **Local** | `docker compose watch` | app + prometheus + loki + alloy + grafana |
| **Remoto** | `docker compose -f docker-compose.remote.yml watch` | app + alloy (observabilidade no servidor) |
| **Servidor** | `docker compose -f docker-compose.server.yml up -d` | prometheus + loki + grafana |

Em outro terminal, acompanhe os logs da API:

```bash
docker compose logs -f api
```

---

## Endereços dos painéis

| Serviço | URL |
|---|---|
| Prometheus | http://localhost:9090 |
| Grafana | http://localhost:3000 |
| Métricas da API | http://localhost:8081/metrics |

---

## Métricas

> Referências: [Metodo RED e USE](https://www.opservices.com.br/conceitos-de-red-e-use/)

### Padrão RED

O padrão RED define três categorias de métricas essenciais para monitorar serviços: **R**ate (taxa), **E**rrors (erros) e **D**uration (duração).

> Referência: [O que é o método RED para observabilidade?](https://dev.to/rafaelbonilha/o-que-e-o-metodo-red-para-observabilidade-3l0i)

**Queries PromQL utilizadas:**

```promql
# R - Taxa de requisições por hora/tempo (Time series)
rate(http_requests_total[1h])

# E - Taxa de erros por minuto/tempo (Time series)
rate(http_requests_total{status=~"5.."}[5m])

# D - P95 de latência — tempo de resposta coberto por 95% das requisições (Time series)
histogram_quantile(0.95, sum(rate(http_request_duration_seconds_bucket[5m])) by (le, path))

# GO - Numero de processos Go (Stat)
go_goroutines

# GO - Memória heap (Stat - Unit[bytes-IEC])
go_memstats_alloc_bytes

# GO - Garbage Collector - Duração média de cada pausa do GC (Time series)
rate(go_gc_duration_seconds_sum[5m]) / rate(go_gc_duration_seconds_count[5m])
```

### Padrão USE

O padrão USE define métricas para **recursos de infraestrutura**: **U**tilization (utilização), **S**aturation (saturação) e **E**rrors (erros).

Diferente do RED, o USE é coletado **fora da aplicação** — pelo sistema operacional ou pelo container runtime. A ferramenta varia conforme o ambiente:

| Ambiente | Ferramenta | O que coleta |
|---|---|---|
| **Com Docker** | cAdvisor | CPU, memória, rede por container |
| **Sem Docker (Linux/VM)** | Node Exporter | CPU, memória, disco, rede da máquina |
| **Sem Docker (Windows)** | Windows Exporter | CPU, memória, disco, rede |
| **Kubernetes** | kube-state-metrics + Node Exporter | Estado dos pods + recursos da máquina |
| **Go runtime** | promhttp (automático) | Goroutines, memória heap, GC |

> Nenhuma dessas ferramentas exige mudança no código da aplicação — elas coletam métricas de infraestrutura externamente. Funcionam com qualquer linguagem (PHP, Node, Delphi, Go, etc.)

---

## Loki + Alloy

O middleware de logging registra cada requisição HTTP com nível baseado no status code:

| Status | Nível |
|---|---|
| 2xx | `info` |
| 4xx | `warn` |
| 5xx | `error` |

Para verificar se o Alloy está rodando corretamente:

```bash
docker compose logs alloy
```

**Pontos-chave do log:**

```
finished node evaluation ... local.file_match.app_logs  → encontrou os arquivos de log
finished node evaluation ... loki.source.file.app_logs  → componente de leitura avaliado
finished node evaluation ... loki.write.local           → envio para o Loki configurado
{^_^} Alloy is running                                  → stack completa de pé
start tailing file ... path=/logs/app.log               → monitorando o arquivo em tempo real
```

### Queries LogQL utilizadas

```logql
# Explore — todos os logs
{filename="/logs/app.log"} | json

# Explore — filtrar por nível
{filename="/logs/app.log"} | json | level="error"

# Dashboard — volume de logs por minuto (Time series)
sum(rate({filename="/logs/app.log"}[1m]))

# Dashboard — erros 5xx ao longo do tempo (Time series)
sum(rate({filename="/logs/app.log"} | json | level="error" [5m]))

# Dashboard — warnings 4xx ao longo do tempo (Time series)
sum(rate({filename="/logs/app.log"} | json | level="warn" [5m]))

# Dashboard — distribuição por nível ao longo do tempo (Time series)
sum by(level) (rate({filename="/logs/app.log"} | json [5m]))
```

---

## Configuração dos serviços no Docker Compose

Cada serviço tem uma origem de configuração diferente. Referências para quando precisar customizar:

| Serviço | Config esperada | Referência |
|---|---|---|
| **Prometheus** | `/etc/prometheus/prometheus.yml` — montado via volume | https://hub.docker.com/r/prom/prometheus |
| **Loki** | `/etc/loki/local-config.yaml` — embutida na imagem (modo monolítico) | https://grafana.com/docs/loki/latest/setup/install/docker/ |
| **Alloy** | Arquivo `.alloy` passado via argumento `run <path>` | https://grafana.com/docs/alloy/latest/get-started/run/docker/ |
| **Grafana** | Sem config inicial — datasources e dashboards configurados pela UI | https://hub.docker.com/r/grafana/grafana |

> Para qualquer nova ferramenta da Grafana Labs, procure `<nome> docker install` na doc oficial — quase sempre tem o exemplo de compose pronto.

---

## Referências PromQL

- [Basics](https://prometheus.io/docs/prometheus/latest/querying/basics/)
- [Functions](https://prometheus.io/docs/prometheus/latest/querying/functions/)
- [Cheat Sheet](https://promlabs.com/promql-cheat-sheet/)
