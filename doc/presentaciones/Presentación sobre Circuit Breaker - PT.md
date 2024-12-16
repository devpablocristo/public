# Apresentação sobre Circuit Breaker

## Introdução ao Circuit Breaker

- **Definição**: O padrão Circuit Breaker (Disjuntor) é um design utilizado em sistemas distribuídos para prevenir que falhas em um serviço ou dependência se propaguem para outros serviços e evitar cascata de erros.
- **Funcionamento**: Semelhante a um disjuntor elétrico, ele se abre para interromper o fluxo quando detecta falhas e se fecha novamente quando o sistema se recupera.

---

## Estados do Circuit Breaker

### 1. Fechado

- **Operação Normal**: Todas as solicitações passam para o serviço de destino.
- **Monitoramento de Erros**: Erros e tempos de resposta são rastreados.
- **Limiar**: Se os erros ultrapassarem um limiar predefinido, o circuito se abre.

### 2. Aberto

- **Bloqueio**: As solicitações são bloqueadas imediatamente e uma resposta de erro é retornada ou um método de fallback é executado.
- **Período de Resfriamento**: Após um tempo, o circuito entra em um estado de teste.

### 3. Meio-Aberto

- **Solicitações de Teste**: Permite a passagem de um número limitado de solicitações para testar se o serviço se recuperou.
- **Recuperação**: Se as solicitações forem bem-sucedidas, o circuito se fecha.
- **Falha**: Se as solicitações falharem, o circuito se abre novamente.

---

## Funcionamento do Circuit Breaker

1. **Monitoramento de Solicitações**: Enquanto o circuito está fechado, as solicitações são monitoradas para detectar falhas e tempos de resposta.
2. **Abertura do Circuito**: Se as falhas ultrapassarem o limiar definido, o circuito se abre. As solicitações subsequentes recebem uma resposta de erro ou um método de fallback é executado.
3. **Período de Resfriamento**: Durante este tempo, não são enviadas solicitações ao serviço com falha.
4. **Reintento de Solicitações**: Após o período de resfriamento, o circuito passa para um estado meio-aberto e permite que um número limitado de solicitações passe.
5. **Fechamento do Circuito**: Se as solicitações de teste forem bem-sucedidas, o circuito se fecha e o fluxo normal de solicitações é retomado. Se falharem, o circuito se abre novamente.

---

## Exemplo de Implementação em Go com Hystrix

1. **Instalação do Hystrix**:
   ```sh
   go get github.com/afex/hystrix-go/hystrix
   ```

2. **Código de Exemplo**:
   ```go
   package main

   import (
       "context"
       "log"

       "github.com/afex/hystrix-go/hystrix"
       "github.com/micro/go-micro/v2"
       "github.com/micro/go-micro/v2/client"
   )

   func main() {
       // Configura o Hystrix
       hystrix.ConfigureCommand("my_command", hystrix.CommandConfig{
           Timeout:               1000,
           MaxConcurrentRequests: 100,
           ErrorPercentThreshold: 25,
       })

       // Cria um novo serviço
       service := micro.NewService(
           micro.Name("example.service"),
       )
       service.Init()

       // Cria um cliente Hystrix
       hystrixClient := client.NewClient(
           client.Wrap(hystrixWrapper),
       )

       // Chama um serviço usando o cliente Hystrix
       req := hystrixClient.NewRequest("example.service", "Example.Endpoint", &Request{})
       rsp := &Response{}

       // Executa a chamada
       err := hystrix.Do("my_command", func() error {
           return hystrixClient.Call(context.Background(), req, rsp)
       }, nil)

       if err != nil {
           log.Fatalf("Erro ao chamar o serviço: %v", err)
       }

       log.Printf("Resposta: %v", rsp)
   }

   func hystrixWrapper(c client.Client) client.Client {
       return &hystrixClient{c}
   }

   type hystrixClient struct {
       client.Client
   }

   type Request struct {
       // define seus campos aqui
   }

   type Response struct {
       // define seus campos aqui
   }
   ```

---

## Exemplo em Go com `goresilience`



1. **Instalação da Biblioteca**:
   ```sh
   go get github.com/slok/goresilience/circuitbreaker
   ```

2. **Código de Exemplo**:
   ```go
   package main

   import (
       "errors"
       "fmt"
       "time"

       "github.com/slok/goresilience"
       "github.com/slok/goresilience/circuitbreaker"
       "github.com/slok/goresilience/retry"
   )

   func main() {
       // Configuração do Circuit Breaker
       cbConfig := circuitbreaker.Config{
           FailureRatio:    0.5,
           MinimumRequests: 10,
           OpenTimeout:     5 * time.Second,
           HalfOpenTimeout: 2 * time.Second,
       }
       cb := circuitbreaker.NewMiddleware(cbConfig)

       // Middleware de reintento
       retryConfig := retry.Config{
           Times: 3,
           WaitBase: 500 * time.Millisecond,
       }
       rt := retry.NewMiddleware(retryConfig)

       // Resiliência composta
       runner := goresilience.RunnerChain(rt, cb)

       // Simula uma função que pode falhar
       err := runner.Run(func() error {
           fmt.Println("Tentando realizar a solicitação...")
           return errors.New("serviço não disponível")
       })

       if err != nil {
           fmt.Println("A solicitação falhou após várias tentativas:", err)
       } else {
           fmt.Println("A solicitação foi bem-sucedida.")
       }
   }
   ```

---

## Comportamento do Circuit Breaker

1. **Configuração**: É configurado com um `FailureRatio` de 0.5, `MinimumRequests` de 10, `OpenTimeout` de 5 segundos e `HalfOpenTimeout` de 2 segundos.
2. **Middleware de Reintento**: É configurado para tentar a operação até 3 vezes com 500 ms de espera entre tentativas.
3. **Runner Compuesto**: Combina o middleware de reintento e o Circuit Breaker para maior resiliência.
4. **Função Simulada**: Simula uma solicitação a um serviço que sempre falha inicialmente.
5. **Execução e Tratamento de Erros**: Mostra mensagens conforme o resultado das solicitações.

---

## Benefícios do Circuit Breaker

1. **Resiliência**: Melhora a resiliência do sistema ao isolar falhas e evitar a propagação de problemas.
2. **Estabilidade**: Mantém a estabilidade da aplicação mesmo quando algumas partes do sistema estão enfrentando problemas.
3. **Desempenho**: Previne que falhas de um serviço degradem o desempenho geral do sistema.
4. **Visibilidade**: Fornece métricas e visibilidade sobre o comportamento dos serviços e seus falhos.

---

## Conclusão

O padrão Circuit Breaker é essencial em arquiteturas de microserviços e sistemas distribuídos para manter a eficiência e estabilidade do sistema. Implementar um Circuit Breaker melhora a resiliência, estabilidade e desempenho da aplicação ao gerenciar de maneira eficaz falhas dos serviços dependentes.

---