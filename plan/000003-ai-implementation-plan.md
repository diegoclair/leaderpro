# Plano de Implementação da IA no LeaderPro

## 1. Visão Geral

### 1.1 Objetivo
Desenvolver um assistente de IA especialista em gestão de pessoas que seja o diferencial competitivo do LeaderPro, capaz de:
- Ser um coach de liderança disponível 24/7
- Ajudar com situações difíceis (conflitos, feedback, demissões)
- Sugerir abordagens personalizadas para cada membro do time
- Analisar dinâmicas de equipe e propor melhorias
- Orientar sobre desenvolvimento de carreira
- Auxiliar em decisões de promoção e reconhecimento
- Sugerir perguntas para 1:1s (entre muitas outras capacidades)

### 1.2 Desafios Principais
1. **Flexibilidade de Dados**: Cada pessoa tem atributos únicos (filhos, pets, hobbies, etc.)
2. **Performance**: Evitar carregar histórico infinito a cada requisição
3. **Relevância Temporal**: Priorizar informações recentes mantendo contexto histórico
4. **Escalabilidade**: Sistema deve funcionar com milhares de usuários e milhões de notas

## 2. Arquitetura de Dados

### 2.1 Tabela de Atributos Dinâmicos

```sql
-- Tabela para atributos flexíveis da pessoa
CREATE TABLE person_attributes (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    person_id BIGINT NOT NULL,
    attribute_key VARCHAR(100) NOT NULL,
    attribute_value TEXT NOT NULL,
    extracted_from_note_id BIGINT NULL, -- Referência à nota que originou o atributo
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (person_id) REFERENCES persons(id),
    FOREIGN KEY (extracted_from_note_id) REFERENCES notes(id),
    UNIQUE KEY unique_person_attribute (person_id, attribute_key),
    INDEX idx_person (person_id)
);

-- Exemplos de uso (tudo como string):
-- ('has_children', 'true', NULL)
-- ('children_names', 'João, Maria', 12345)
-- ('preferred_meeting_time', '14:00', NULL)
-- ('communication_style', 'direct', 12346)
-- ('hobbies', 'corrida, leitura, videogame', NULL)
```

### 2.2 Tabela de Prompts

```sql
-- Versionamento e histórico de prompts
CREATE TABLE ai_prompts (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    type VARCHAR(50) NOT NULL, -- 'leadership_coach', 'attribute_extraction', 'meeting_suggestions'
    version INT NOT NULL,
    prompt TEXT NOT NULL,
    model VARCHAR(50) NOT NULL, -- 'gpt-4', 'gpt-3.5-turbo', etc
    temperature DECIMAL(3,2) DEFAULT 0.7,
    max_tokens INT DEFAULT 2000,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by BIGINT NOT NULL,
    FOREIGN KEY (created_by) REFERENCES users(id),
    INDEX idx_type_active (type, is_active),
    UNIQUE KEY unique_type_version (type, version)
);
```

### 2.3 Tabela de Uso e Feedback

```sql
-- Rastreamento de uso da IA
CREATE TABLE ai_usage_tracker (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    user_id BIGINT NOT NULL,
    company_id BIGINT NOT NULL,
    prompt_id BIGINT NOT NULL,
    person_id BIGINT NULL, -- Contexto da pessoa (se aplicável)
    request_type VARCHAR(50) NOT NULL, -- 'chat', 'extraction', 'suggestion'
    tokens_used INT NOT NULL,
    cost_usd DECIMAL(10,6) NOT NULL,
    response_time_ms INT NOT NULL,
    feedback ENUM('helpful', 'not_helpful', 'neutral') NULL,
    feedback_comment TEXT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (company_id) REFERENCES companies(id),
    FOREIGN KEY (prompt_id) REFERENCES ai_prompts(id),
    FOREIGN KEY (person_id) REFERENCES persons(id),
    INDEX idx_user_date (user_id, created_at),
    INDEX idx_company_date (company_id, created_at),
    INDEX idx_feedback (feedback)
);

-- Conteúdo das conversas (separado para facilitar expurgo)
CREATE TABLE ai_conversations (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    usage_id BIGINT NOT NULL,
    user_message TEXT NOT NULL,
    ai_response TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP DEFAULT (CURRENT_TIMESTAMP + INTERVAL 180 DAY), -- 180 dias padrão
    FOREIGN KEY (usage_id) REFERENCES ai_usage_tracker(id),
    INDEX idx_usage (usage_id),
    INDEX idx_expires (expires_at)
);
```

**Vantagens da estrutura separada:**
- Facilita expurgo de dados antigos (DELETE FROM ai_conversations WHERE expires_at < NOW())
- Mantém métricas e feedback permanentemente em ai_usage_tracker
- Permite análise de custos e performance sem dados sensíveis
- Flexibilidade para diferentes políticas de retenção

### 2.4 Job de Limpeza de Dados

```sql
-- Procedure para expurgo automático
DELIMITER //
CREATE PROCEDURE cleanup_ai_conversations()
BEGIN
    DELETE FROM ai_conversations 
    WHERE expires_at < NOW() 
    LIMIT 1000; -- Processa em lotes para não travar o banco
END//
DELIMITER ;

-- Agendar para rodar diariamente
CREATE EVENT cleanup_old_ai_data
ON SCHEDULE EVERY 1 DAY
DO CALL cleanup_ai_conversations();
```

## 3. Arquitetura Simplificada da IA

### 3.1 Princípio: Uma Pessoa por Vez

A IA sempre trabalha no contexto de UMA pessoa específica:
- Busca as notas recentes dessa pessoa (últimos 3-6 meses)
- Envia para a IA processar
- Extrai insights e atributos
- Salva no banco

**Sem necessidade de:**
- Vector stores
- Embeddings complexos
- Cache em múltiplas camadas
- Processamento assíncrono elaborado

### 3.2 Fluxo Correto

#### Parte 1: Assistente de Liderança (Chat)
```
1. Usuário faz pergunta/pede ajuda (ex: "como dar feedback para João?")
2. Busca contexto da pessoa + notas recentes + perfil (person_attributes)
3. Envia para OpenAI com system prompt de especialista em gestão
4. Recebe resposta contextualizada
5. Responde pro usuário de forma conversacional
```

#### Parte 2: Após o 1:1 (Extração de Perfil)
```
6. Usuário sobe anotações do 1:1
7. IA analisa as notas e extrai informações para o perfil (se for informação que caracteriza o perfil da pessoa)
```

### 3.3 Contexto Temporal Pragmático

```go
type PersonAIContext struct {
    Person struct {
        ID         int64             `json:"id"`
        Name       string            `json:"name"`
        Attributes map[string]string `json:"attributes"` // Da tabela person_attributes
    } `json:"person"`
    
    RecentNotes []entity.Note `json:"recent_notes"` // Últimas 50-100 notas (3-6 meses)
    
    LastMeeting *struct {
        Date  time.Time `json:"date"`
        Notes string    `json:"notes"`
    } `json:"last_meeting,omitempty"`
}
```

## 4. Integração Simples com OpenAI

### 4.1 Serviço de IA Especialista

```go
type AIService struct {
    openaiClient *openai.Client
}

func NewAIService(apiKey string) *AIService {
    return &AIService{
        openaiClient: openai.NewClient(apiKey),
    }
}

// Chat principal - assistente de liderança
func (s *AIService) ChatWithLeadershipCoach(ctx context.Context, userMessage string, person *entity.Person, recentNotes []entity.Note, attributes map[string]string) (string, error) {
    systemPrompt := s.buildSystemPrompt()
    contextPrompt := s.buildContextPrompt(person, recentNotes, attributes)
    
    resp, err := s.openaiClient.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
        Model: openai.GPT4oMini,
        Messages: []openai.ChatCompletionMessage{
            {Role: openai.ChatMessageRoleSystem, Content: systemPrompt},
            {Role: openai.ChatMessageRoleUser, Content: contextPrompt + "\n\nPergunta do líder: " + userMessage},
        },
        Temperature: 0.7,
    })
    
    if err != nil {
        return "", err
    }
    
    return resp.Choices[0].Message.Content, nil
}

// Extração automática de atributos após notas
func (s *AIService) ExtractAttributes(ctx context.Context, person entity.Person, notes []entity.Note) (map[string]string, error) {
    prompt := s.buildAttributePrompt(person, notes)
    // ... implementação similar
}
```

### 4.2 Prompts Especializados

#### System Prompt - Especialista em Gestão
```go
func (s *AIService) buildSystemPrompt() string {
    return `Você é um coach especialista em gestão de pessoas e liderança tecnológica com 20 anos de experiência. 
Você ajuda líderes tech a se tornarem melhores gestores através de conselhos práticos e personalizados.

Suas especialidades incluem:
- Dar e receber feedback construtivo
- Conduzir reuniões 1:1 eficazes
- Gerenciar conflitos e situações difíceis
- Desenvolver e promover talentos
- Criar ambientes psicologicamente seguros
- Balancear demandas técnicas com gestão de pessoas
- Lidar com diferentes perfis comportamentais

Sempre responda de forma:
- Prática e acionável
- Empática mas direta
- Baseada no contexto específico da pessoa
- Com exemplos concretos quando relevante`
}
```

#### Context Prompt - Informações da Pessoa
```go
func (s *AIService) buildContextPrompt(person *entity.Person, notes []entity.Note, attributes map[string]string) string {
    if person == nil {
        return "Contexto geral da equipe."
    }
    
    var profileLines []string
    for key, value := range attributes {
        profileLines = append(profileLines, fmt.Sprintf("%s: %s", key, value))
    }
    
    var recentEvents []string
    for i, note := range notes {
        if i >= 10 { break } // Últimos 10 eventos mais relevantes
        recentEvents = append(recentEvents, fmt.Sprintf("- %s: %s", 
            note.CreatedAt.Format("2006-01-02"), note.Content))
    }
    
    return fmt.Sprintf(`CONTEXTO SOBRE %s:

PERFIL:
%s

EVENTOS RECENTES:
%s`, 
        person.Name, 
        strings.Join(profileLines, "\n"), 
        strings.Join(recentEvents, "\n"))
}
```

#### Prompt para Extração de Atributos
```go
func (s *AIService) buildAttributePrompt(person entity.Person, notes []entity.Note) string {
    var noteContents []string
    for _, note := range notes {
        noteContents = append(noteContents, note.Content)
    }
    
    return fmt.Sprintf(`Analise as notas sobre %s e extraia APENAS informações que você tem 
100%% de certeza. Retorne um JSON com pares chave-valor simples.

NOTAS:
%s

Exemplo de resposta:
{"has_children": "true", "children_names": "João, Maria", "hobbies": "corrida"}

Se não tiver certeza absoluta sobre algo, NÃO inclua.`, 
        person.Name, 
        strings.Join(noteContents, "\n"))
}
```

## 5. Implementação em 3 Passos

### Passo 1: Base de Dados (1-2 dias)

**Migration 000006 - Estrutura da IA:**
- Tabela `person_attributes` (seção 2.1)
- Tabela `ai_prompts` (seção 2.2) 
- Tabela `ai_usage_tracker` (seção 2.3)
- Tabela `ai_conversations` (seção 2.3)

**Migration 000007 - Job de limpeza automática:**
- Procedure `cleanup_ai_conversations()` (seção 2.4)
- Event scheduler para executar diariamente

### Passo 2: Endpoints de IA (3-5 dias)
```go
// POST /companies/:company_id/ai/chat
type ChatRequest struct {
    Message  string  `json:"message"`
    PersonID *int64  `json:"person_id,omitempty"` // Opcional, para contexto específico
}

func ChatWithAI(c echo.Context) error {
    // 1. Validar request
    // 2. Se PersonID fornecido, buscar contexto da pessoa
    // 3. Chamar AIService.ChatWithLeadershipCoach
    // 4. Retornar resposta
}

// Exemplos de perguntas que o usuário pode fazer:
// - "Como dar feedback construtivo para João sobre atrasos?"
// - "Sugestões de perguntas para 1:1 com Maria"
// - "Como lidar com conflito entre dois devs sênior?"
// - "João está desmotivado, o que fazer?"
// - "Como preparar Maria para uma promoção?"

// POST /companies/:company_id/people/:person_id/notes (hook após criar nota)
func OnNoteCreated(note Note) {
    // Se contém informações sobre a pessoa, extrair atributos em background
}
```

### Passo 3: Frontend Integration (2-3 dias)
```typescript
// Chat com IA - pode ser um botão flutuante ou seção dedicada
interface AIChat {
  message: string;
  personId?: number; // Opcional, preenchido automaticamente se estiver no perfil
}

const ChatWithAI = () => {
  const [message, setMessage] = useState("");
  const [conversation, setConversation] = useState<Message[]>([]);
  
  const sendMessage = async () => {
    const response = await apiClient.post('/ai/chat', {
      message,
      personId: currentPersonId // se estiver no contexto de uma pessoa
    });
    
    setConversation([...conversation, 
      { role: 'user', content: message },
      { role: 'assistant', content: response.data }
    ]);
  };
  
  return <ChatInterface />;
};
```

## 6. Custo Estimado

### Cenário Típico
- Usuário faz ~5 pedidos de sugestões/mês
- ~100 notas analisadas por pedido
- Modelo: GPT-4o-mini

**Custo aproximado: $0.10-0.20 por usuário/mês**

Bem abaixo do valor da mensalidade (R$ 49.90), tornando viável.

## 7. Segurança Básica

### Pontos Essenciais
- Dados ficam sempre no nosso banco
- OpenAI API com HTTPS (dados não são armazenados pela OpenAI)
- Rate limiting para evitar abuso
- Logs para auditoria

## 8. MVP em 1 Semana

### Checklist
- [ ] Migration da tabela person_attributes
- [ ] Implementar AIService com ChatWithLeadershipCoach
- [ ] Endpoint POST /ai/chat (assistente de liderança)
- [ ] Hook após criar nota (extrai atributos em background)
- [ ] Interface de chat no frontend
- [ ] System prompt especialista em gestão de pessoas

### Exemplos de Uso Real

**Cenário 1: Feedback Difícil**
- Líder: "João está chegando atrasado há 2 semanas. Como abordar?"
- IA: *Analisa histórico de João, vê que tem filhos pequenos*
- IA: "Considerando que João tem filhos pequenos, sugiro uma abordagem empática..."

**Cenário 2: Preparação para 1:1**
- Líder: "Vou ter 1:1 com Maria amanhã, o que perguntar?"
- IA: *Analisa últimas notas, vê que Maria mencionou interesse em arquitetura*
- IA: "Baseado no interesse de Maria em arquitetura, considere perguntar sobre..."

**Cenário 3: Desenvolvimento de Carreira**
- Líder: "Como ajudar Pedro a crescer para sênior?"
- IA: *Analisa perfil e gaps identificados*
- IA: "Pedro precisa desenvolver estas habilidades: ..."

## 9. Avaliação Futura: MCP (Model Context Protocol)

### 9.1 O que é MCP
MCP (Model Context Protocol) é um protocolo padronizado que permite à IA interagir diretamente com sistemas externos, incluindo bancos de dados. Existem implementações maduras para MySQL em 2025.

### 9.2 Potencial para LeaderPro
```go
// Possível implementação futura
type MCPProvider interface {
    InsertPersonAttribute(ctx context.Context, attr PersonAttribute) error
    GetPersonContext(ctx context.Context, personID int64) (PersonContext, error)
}
```

### 9.3 Benefícios Potenciais
- ✅ IA escreve diretamente no banco (sem passar pelo código)
- ✅ Maior precisão (IA vê schema real, evita alucinações)
- ✅ Redução de código intermediário
- ✅ Protocolo padronizado e seguro

### 9.4 Riscos e Considerações
- ⚠️ **Segurança**: IA com acesso direto ao banco
- ⚠️ **Controle**: Bypass dos controles de validação da aplicação
- ⚠️ **Auditoria**: Mais difícil rastrear mudanças
- ⚠️ **Debugging**: Harder to debug issues
- ⚠️ **Compliance**: Pode violar políticas de acesso a dados

### 9.5 Implementações Disponíveis
- `@benborla29/mcp-server-mysql` (Node.js)
- `designcomputer/mysql_mcp_server` (Python)
- Azure Database for MySQL MCP Server (Microsoft)

### 9.6 Recomendação
**Implementar primeiro sem MCP**. Avaliar MCP apenas após:
1. Sistema básico funcionando em produção
2. Análise de segurança completa
3. Definição de políticas de acesso
4. Testes extensivos em ambiente isolado

## Conclusão Simplificada

**Filosofia**: Começar ultra simples e iterar rapidamente.

- ✅ Uma pessoa por vez
- ✅ Sem cache complexo
- ✅ Sem vector stores
- ✅ Apenas OpenAI
- ✅ Repositórios tradicionais (sem MCP inicialmente)
- ✅ Custo baixo
- ✅ Implementação em 1 semana

O sistema vai aprender e melhorar conforme os usuários usam, sem complexidade desnecessária.