# Guia de Observabilidade — Go Stack

---

## Os três pilares da Observabilidade

Observabilidade tem **3 pilares**. Cada um responde uma pergunta diferente sobre o sistema:

| Pilar | Pergunta que responde | Ferramenta na stack |
|---|---|---|
| **Logs** | O que aconteceu? | Zap + Loki |
| **Métricas** | Quanto/quantas vezes aconteceu? | Prometheus |
| **Traces** | Por onde passou a requisição? | OpenTelemetry + Tempo |

Você pode usar um, dois ou os três — depende da complexidade do sistema.

---

## Cada peça da stack

### Zap
Produtor de logs dentro da aplicação Go. Gera logs estruturados em JSON, que é o formato que o Loki sabe consumir bem. Você já tem familiaridade com ele.

### Loki
Banco de dados de logs da Grafana Labs. Diferente do Elasticsearch, ele **não indexa o conteúdo** dos logs — indexa apenas os **labels** (metadados). Isso o torna mais leve e barato de operar.

Fluxo:
```
App (Zap) → Alloy (agente coletor) → Loki (armazena) → Grafana (visualiza)
```

- Docs: https://grafana.com/docs/loki/latest/

### Prometheus
Banco de dados de métricas. Funciona no modelo **pull**: ele periodicamente visita um endpoint `/metrics` da sua aplicação e coleta os números.

Tipos de métrica:
- **Counter** — só sobe (ex: total de requisições)
- **Gauge** — sobe e desce (ex: memória em uso)
- **Histogram** — distribui valores em buckets (ex: latência de resposta)

Fluxo:
```
App expõe /metrics ← Prometheus raspa a cada X segundos → Grafana (visualiza)
```

- Docs: https://prometheus.io/docs/introduction/overview/

### Grafana
Painel de visualização. **Não armazena nada** — lê de fontes de dados (Loki, Prometheus, Tempo) e exibe dashboards. É o único lugar onde você vai "ver" tudo junto.

---

## OpenTelemetry e Tempo — o terceiro pilar

### OpenTelemetry (OTel)
Padrão aberto (CNCF) para instrumentação. Em vez de usar SDK proprietário do Jaeger, Zipkin ou Tempo, você usa o SDK do OTel e configura apenas o destino. É o princípio "escreva uma vez, envie para qualquer lugar".

Na prática: você instrumenta seu código Go com `go.opentelemetry.io/otel` e decide depois para onde os traces vão.

Conceitos-chave:
- **Trace** — representa uma requisição completa passando pelo sistema
- **Span** — uma operação individual dentro do trace (ex: chamada ao banco, chamada HTTP)
- **Context propagation** — como o trace atravessa serviços diferentes

- Docs: https://opentelemetry.io/docs/languages/go/

### Tempo
Banco de dados de traces distribuídos da Grafana Labs. Armazena os traces que o OTel envia. Integra nativamente com o Grafana, permitindo navegar de um log direto para o trace correspondente.

Fluxo:
```
App (OTel SDK) → OTel Collector → Tempo (armazena) → Grafana (visualiza)
```

---

## Como tudo se comunica — visão geral

```
┌─────────────────────────────────────────────┐
│                 Sua App Go                  │
│   Zap (logs) + prom client + OTel SDK      │
└────────┬────────────┬────────────┬──────────┘
         │            │            │
         ▼            ▼            ▼
       Alloy      /metrics      OTel Collector
      (agente)   endpoint
         │            │            │
         ▼            ▼            ▼
        Loki      Prometheus     Tempo
         │            │            │
         └────────────┴────────────┘
                      │
                   Grafana
```

---

## Vale usar tudo?

Depende do objetivo do sistema:

| Cenário | O que usar |
|---|---|
| Aprender observabilidade | Comece só com **Prometheus + Grafana** |
| Entender logs estruturados | Adicione **Loki** |
| Sistema distribuído / microserviços | Adicione **OTel + Tempo** |
| Monolito simples | Traces são overkill — logs + métricas bastam |

Para o projeto de aprendizado, **faz sentido usar tudo** — justamente para entender como cada peça se encaixa. Mas a ordem de aprendizado importa.

---

## Ordem sugerida de estudo

### 1. Prometheus + Grafana
- Entenda o modelo pull
- Tipos de métrica: counter, gauge, histogram
- Exponha `/metrics` em Go com `prometheus/client_golang`
- Monte um dashboard básico no Grafana

### 2. Loki + Alloy

#### Modos de implantação do Loki

O Loki tem três modos de deploy. A escolha impacta complexidade e escala:

| Modo | Como funciona | Quando usar |
|---|---|---|
| **Monolítico** | Todos os componentes num processo só | Aprendizado, projetos pequenos |
| **Simples escalável** | Separa leitura e escrita em dois processos | Médio volume, uma instância por papel |
| **Microsserviços** | Cada componente roda separado | Alto volume, escala independente |

> **Nesse projeto usamos o modo Monolítico** — um container, zero config extra. É suficiente para aprendizado e visualização local.

Os outros dois modos serão explorados como exemplos comparativos após dominar o monolítico.

- Docs oficiais dos modos: https://grafana.com/docs/loki/latest/get-started/deployment-modes/

#### O que implementar

- Configure coleta de logs com Alloy
- Aprenda LogQL (linguagem de query do Loki)
- Visualize logs no Grafana

### 3. OpenTelemetry + Tempo
- Instrumente uma rota HTTP com trace
- Veja os spans no Grafana via Tempo
- Entenda a correlação entre log → trace

### 4. Feature real + Banco de dados
> **Introduzir o DB aqui** — junto com a primeira feature que precisar persistir dados. Não antes. Sem uma feature real não há como saber quais tabelas serão necessárias.
- Adicionar PostgreSQL no `docker-compose.yml`
- Conectar a aplicação Go ao banco
- Observar queries com logs e métricas

### 5. cAdvisor + USE em outros sistemas

O método USE (Utilization, Saturation, Errors) mede **recursos de infraestrutura** — CPU, memória, disco, rede. Diferente do RED, ele é coletado **fora da aplicação**, sem tocar no código. Funciona com qualquer linguagem.

> **Introduzir o cAdvisor aqui** — quando quiser métricas USE dos containers. Basta adicionar no `docker-compose.yml`.

**Ferramentas por ambiente:**

| Ambiente | Ferramenta | O que coleta |
|---|---|---|
| **Com Docker** | cAdvisor | CPU, memória, rede por container |
| **Sem Docker (Linux/VM)** | Node Exporter | CPU, memória, disco, rede da máquina |
| **Sem Docker (Windows)** | Windows Exporter | CPU, memória, disco, rede |
| **Kubernetes** | kube-state-metrics + Node Exporter | Estado dos pods + recursos da máquina |
| **Go runtime** | promhttp (automático) | Goroutines, memória heap, GC |

- cAdvisor: https://github.com/google/cadvisor
- Node Exporter: https://github.com/prometheus/node_exporter
- Windows Exporter: https://github.com/prometheus-community/windows_exporter

---

### 6. Nginx
> **Introduzir apenas quando o ambiente pedir** — reverse proxy, SSL termination ou load balancing. Em desenvolvimento local é desnecessário; o Gin já serve HTTP direto.
- Adicionar Nginx no `docker-compose.yml`
- Configurar como reverse proxy para a aplicação

---

## Ordem de desenvolvimento num projeto real

Observabilidade não é uma feature adicionada no final. É **infraestrutura transversal** — como tratamento de erros ou configuração.

### A sequência correta

```
1. Estrutura base mínima
   └── main.go, config, um handler de saúde (/health)

2. Observabilidade base (infraestrutura)
   └── Logger (Zap) configurado globalmente
   └── /metrics endpoint (Prometheus)
   └── Middleware de log de requisições
   └── Middleware de métricas HTTP

3. Features reais (handler + service + repo)
   └── Já nascem com log e métrica plugados via middleware

4. Observabilidade de domínio (contextual)
   └── Logs com contexto específico (ex: "pedido criado", user_id)
   └── Métricas de negócio (ex: total de pedidos, valor médio)
   └── Traces em operações críticas
```

### Por que essa ordem?

O **middleware** é a chave. Quando você configura log e métricas no middleware HTTP antes de criar as rotas de negócio, todo handler que você criar depois já é observado **automaticamente**.

```
Request → [middleware: log] → [middleware: metrics] → Handler → Response
```

Você escreve o middleware uma vez. Cada novo handler ganha observabilidade de graça.

### Observabilidade de infraestrutura vs. de domínio

| Tipo | Exemplo | Quando implementar |
|---|---|---|
| **Infraestrutura** | Latência HTTP, status codes, erros 5xx | No início, via middleware |
| **Domínio** | Pedidos criados, usuários ativos, valor processado | Junto com cada feature |

A de infraestrutura cobre ~80% das necessidades. A de domínio você adiciona conforme entende o negócio.

---

## Etapa 1 — Base do projeto Go

Antes de qualquer observabilidade, o projeto precisa de uma base sólida. Siga as etapas abaixo em ordem.

### 1. Inicializar o módulo Go

Todo projeto Go começa com isso:

```
go mod init github.com/seu-usuario/observabilidade
```

O nome do módulo é convenção — use seu usuário do GitHub ou qualquer identificador reverso. Isso cria o `go.mod`, que é o equivalente ao `package.json` do Node.

- Doc oficial: https://go.dev/doc/modules/gomod-ref

---

### 2. Estrutura de pastas

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

**Por que `internal/`?** Tudo dentro dessa pasta não pode ser importado por outros módulos Go — é uma proteção do compilador. É a convenção para código privado do projeto.

- Referência sobre estrutura: https://github.com/golang-standards/project-layout
  > Leia com calma, mas não siga cegamente — para projetos pequenos, menos pastas é melhor.

---

### 3. Servidor HTTP com Gin

Go tem um servidor HTTP na biblioteca padrão (`net/http`), mas para projetos reais é comum usar um framework. Para esse projeto, use o **Gin** — rápido, maduro e com bastante material de referência:

```
go get github.com/gin-gonic/gin
```

**Por que Gin?** É um dos frameworks mais usados no ecossistema Go. Tem roteamento, middlewares e binding de JSON integrados, sem precisar reinventar a roda.

- Doc do Gin: https://gin-gonic.com/docs/

---

### 4. Endpoint `/health`

O primeiro endpoint de qualquer serviço. Retorna `200 OK` com um JSON simples:

```json
{ "status": "ok" }
```

**Para que serve:** load balancers, orquestradores (Kubernetes) e ferramentas de monitoramento consultam esse endpoint para saber se o serviço está vivo. Prometheus também pode usar.

---

### 5. Configuração via variáveis de ambiente

Nada de valores hardcoded. Porta, URLs, nomes de serviço — tudo via env var. A abordagem mais simples é uma struct `Config` lida no startup:

```
go get github.com/kelseyhightower/envconfig
```

- Doc: https://github.com/kelseyhightower/envconfig

> Alternativa popular (mais poderosa): https://github.com/spf13/viper — mas é overkill para esse projeto.

---

### 6. Dockerfile

O `Dockerfile` empacota o binário Go em uma imagem Docker. Para projetos Go, o padrão é usar **multi-stage build**:

- **Stage 1 (builder)** — compila o binário usando uma imagem com Go instalado
- **Stage 2 (runtime)** — copia apenas o binário para uma imagem mínima (`alpine` ou `scratch`)

Isso resulta em imagens pequenas — sem o compilador Go, sem código-fonte, só o binário.

- Dockerfile reference: https://docs.docker.com/reference/dockerfile/
- Multi-stage builds: https://docs.docker.com/build/building/multi-stage/

---

### 7. docker-compose

O `docker-compose.yml` orquestra todos os serviços juntos. Por enquanto só a app, mas conforme avançamos nas etapas vamos adicionando Prometheus, Grafana, Loki e Tempo.

**Conceitos importantes:**
- `services` — lista de containers que vão subir
- `build` — instrui o compose a buildar a imagem a partir do `Dockerfile`
- `ports` — mapeia porta do container para a máquina host
- `env_file` — carrega o `.env` para dentro do container
- `volumes` — monta arquivos/pastas do host dentro do container

- docker-compose reference: https://docs.docker.com/compose/compose-file/

---

### 8. Zap configurado no startup

Você já conhece o Zap, então aqui é só definir **como** ele vai ser configurado no projeto:

- Modo `development` para rodar local (output legível, colorido)
- Modo `production` para rodar com Loki (output JSON)
- Logger criado uma vez no `main.go` e passado para quem precisar

- Doc oficial do Zap: https://pkg.go.dev/go.uber.org/zap

---

### O que fazer agora

- [ ] Rodar o `go mod init`
- [ ] Criar as pastas conforme a estrutura acima
- [ ] Criar o `cmd/api/main.go` com o servidor Gin subindo na porta `8080`
- [ ] Adicionar a rota `/health` retornando `{ "status": "ok" }`
- [ ] Configurar o envconfig com `.env`
- [ ] Criar o `Dockerfile` com multi-stage build
- [ ] Criar o `docker-compose.yml` e o `.dockerignore`
- [ ] Configurar o Zap no startup e logar `"server started"` com o campo `port`

Quando o container estiver rodando e o `/health` respondendo, a próxima etapa é adicionar o **Prometheus**.
