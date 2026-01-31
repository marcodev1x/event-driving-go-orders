# Sistema de Pedidos Event-Driven em Go

## Visão Geral da Arquitetura

Este projeto implementa uma arquitetura de microsserviços event-driven para processamento de pedidos utilizando Go, Kafka e MySQL. O sistema segue padrões de comunicação assíncrona para lidar com criação de pedidos, validação de pagamento e processos de negócio subsequentes.

## Implementação Atual

### Serviços

#### Order Service
- **Porta**: 8080
- **Database**: MySQL (banco de dados orders)
- **Cache**: Redis
- **Responsabilidades**:
  - Gerencia criação de pedidos via endpoint `/create-checkout`
  - Persiste pedidos com status "pendente pagamento"
  - Publica eventos `order.created` no Kafka
  - Consome eventos de confirmação de pagamento do Kafka

#### Payment Service
- **Porta**: 8082
- **Cache**: Redis
- **Responsabilidades**:
  - Consome eventos `order.created` do Kafka
  - Realiza validação de pagamento
  - Publica eventos de status de pagamento (`payment.confirmed` ou `payment.failed`) no Kafka

### Fluxo de Eventos

1. **Criação do Pedido**
   - Cliente envia requisição para `order-service/create-checkout`
   - Pedido é salvo com status "pendente pagamento" no MySQL
   - Evento `order.created` é publicado no Kafka

2. **Processamento de Pagamento**
   - `payment-service` consome evento `order.created`
   - Validação de pagamento é executada
   - Com base no resultado:
     - Sucesso: evento `payment.confirmed` publicado
     - Falha: evento `payment.failed` publicado

3. **Atualização de Status do Pedido**
   - `order-service` consome eventos de status de pagamento
   - Status do pedido é atualizado no MySQL
   - Pedidos com falha são marcados como "com falha"

## Componentes de Infraestrutura

### Message Broker
- **Kafka**: Plataforma de streaming de eventos para comunicação assíncrona
- **Tópicos**:
  - `order.created`: Acionado quando novos pedidos são criados
  - `payment.confirmed`: Validação de pagamento bem-sucedida
  - `payment.failed`: Falha na validação de pagamento

### Bancos de Dados
- **MySQL**: Armazenamento principal de informações de pedidos
- **Redis**: Camada de cache para otimização de performance

### Ferramentas de Desenvolvimento
- **Kafka UI**: Interface web para gerenciamento do Kafka (Porta 8081)

## Melhorias Planejadas

### Notify Service
- **Propósito**: Sistema de notificação interna de compras
- **Gatilho**: Consome eventos `payment.confirmed`
- **Funcionalidade**: Enviar alertas internos para compras bem-sucedidas

### Inventory Service
- **Propósito**: Integração com gestão de estoque
- **Gatilho**: Consome eventos `payment.confirmed`
- **Funcionalidade**: 
  - Recuperar informações de produto e estoque
  - Reduzir níveis de inventário para pedidos confirmados
  - Comunicar com microsserviço externo de estoque

### Fluxo Futuro
1. Após confirmação de pagamento:
   - `notify-service` processa notificações internas de compra
   - `inventory-service` gerencia redução de estoque
   - Eventos adicionais podem ser publicados para atualizações de inventário

## Especificações Técnicas

### Stack Tecnológico
- **Linguagem**: Go
- **Web Framework**: Gin
- **Message Queue**: Apache Kafka
- **Database**: MySQL 8.0
- **Cache**: Redis 8.2
- **Containerização**: Docker & Docker Compose

### Estrutura do Projeto
```
services/
├── order-service/
│   ├── internal/
│   │   ├── domain/
│   │   ├── usecases/
│   │   ├── repository/
│   │   └── rest/
│   ├── kafka/
│   │   ├── events/
│   │   ├── producer/
│   │   └── consumer/
│   └── infra/
└── payment-service/
    ├── internal/
    │   ├── domain/
    │   ├── usecases/
    │   └── rest/
    ├── kafka/
    │   ├── events/
    │   ├── producer/
    │   └── consumer/
    └── infra/
```

### Configuração de Ambiente
Cada serviço utiliza variáveis de ambiente para:
- Conexões de banco de dados
- Configuração do Redis
- Configurações do broker Kafka
- Parâmetros específicos do serviço

### Endpoints da API

#### Order Service
- `GET /create-checkout`: Cria novo checkout de pedido

#### Payment Service
- Sem endpoints HTTP diretos (apenas event-driven)

## Configuração de Desenvolvimento

### Pré-requisitos
- Docker & Docker Compose
- Go 1.21+

### Executando a Aplicação
```bash
docker-compose up -d
```

Isso iniciará:
- MySQL (Porta 3306)
- Kafka (Porta 9092)
- Kafka UI (Porta 8081)
- Redis (Porta 6380)
- Order Service (Porta 8080)
- Payment Service (Porta 8082)

## Schema de Eventos

### Evento Order Created
```json
{
  "event_id": "string",
  "event_type": "order.created",
  "timestamp": "datetime",
  "content_id": "int",
  "checkout": {
    "id": "int",
    "price": "float",
    "status": "string"
  }
}
```

### Evento Payment Confirmed
```json
{
  "event_id": "string",
  "event_type": "payment.confirmed",
  "timestamp": "datetime",
  "content_id": "int",
  "order_id": "int"
}
```

### Evento Payment Failed
```json
{
  "event_id": "string",
  "event_type": "payment.failed",
  "timestamp": "datetime",
  "content_id": "int",
  "order_id": "int"
}
```

## Monitoramento e Observabilidade

### Logging
- Logging estruturado com informações contextuais
- Rastreamento de eventos com IDs de correlação
- Logging de erros com stack traces

### Health Checks
- Verificações de conectividade com banco de dados
- Disponibilidade do broker Kafka
- Validação de conexão Redis

## Considerações de Escalabilidade

### Configuração Kafka
- Setup de broker único para desenvolvimento
- Configurável para deployment multi-broker em produção
- Criação automática de tópicos habilitada

### Escalabilidade de Serviços
- Serviços stateless permitem escalonamento horizontal
- Pool de conexões com banco de dados
- Clustering Redis para alta disponibilidade
