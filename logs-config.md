# Guia de integração com a stack de observabilidade

Checklist para integrar qualquer projeto Go com o Loki e o Prometheus centralizados.

---

## 1. Variáveis de ambiente (`.env`)

```env
PROMETHEUS_HOST=monitor.topsoft.inf.br
LOKI_HOST=logs.topsoft.inf.br
```

---

## 2. `docker-compose.yml`

Adicione o serviço `alloy` e certifique-se de que o serviço da API exponha os logs em volume.

```yaml
services:
  api:
    build: .
    env_file:
      - .env
    volumes:
      - ./logs:/app/logs   # a API deve escrever logs em /app/logs/*.log
    # ... demais configurações

  alloy:
    image: grafana/alloy:latest
    ports:
      - "12345:12345"
    env_file:
      - .env
    volumes:
      - ./alloy/config.alloy:/etc/alloy/config.alloy
      - ./logs:/logs
    command: run /etc/alloy/config.alloy
```

---

## 3. `alloy/config.alloy`

Crie o arquivo com o conteúdo abaixo. O Alloy lê os logs locais e envia tanto logs quanto métricas para os servidores remotos.

```hcl
local.file_match "app_logs" {
  path_targets = [{"__path__" = "/logs/*.log"}]
}

loki.source.file "app_logs" {
  targets    = local.file_match.app_logs.targets
  forward_to = [loki.write.remote.receiver]
}

loki.write "remote" {
  endpoint {
    url = "http://" + env("LOKI_HOST") + "/loki/api/v1/push"
  }
}

prometheus.scrape "api_metrics" {
  targets = [{"__address__" = "api:8081"}]
  forward_to = [prometheus.remote_write.remote.receiver]
}

prometheus.remote_write "remote" {
  endpoint {
    url            = "http://" + env("PROMETHEUS_HOST") + "/api/v1/write"
    remote_timeout = "60s"
  }
}
```

> **Requisito:** a API precisa expor um endpoint `/metrics` no formato Prometheus (ex: usando `promhttp.Handler()`).

---

## 4. Verificação

Suba os containers e confira os logs do Alloy:

```bash
docker compose up -d
docker compose logs alloy
```

Procure por estas linhas pra confirmar que está tudo ok:

| Log | Significa |
|---|---|
| `finished node evaluation ... loki.write.remote` | envio pro Loki configurado |
| `finished node evaluation ... prometheus.remote_write.remote` | envio pro Prometheus configurado |
| `{^_^} Alloy is running` | stack completa de pé |
| `start tailing file ... path=/logs/app.log` | monitorando o arquivo de log |
