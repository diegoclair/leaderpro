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

## 2. Arquitetura de Dados Proposta

### 2.1 Modelo Híbrido: Relacional + Flexível

#### Opção 1: Tabela de Atributos Dinâmicos (Recomendada)
```sql
-- Tabela principal permanece como está
CREATE TABLE persons (
    id BIGINT PRIMARY KEY,
    name VARCHAR(255),
    email VARCHAR(255),
    -- campos estruturados essenciais
);

-- Nova tabela para atributos flexíveis (ULTRA SIMPLIFICADA)
CREATE TABLE person_attributes (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    person_id BIGINT NOT NULL,
    attribute_key VARCHAR(100) NOT NULL,
    attribute_value TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (person_id) REFERENCES persons(id),
    UNIQUE KEY unique_person_attribute (person_id, attribute_key),
    INDEX idx_person (person_id)
);

-- Exemplos de uso (tudo como string):
-- ('has_children', 'true')
-- ('children_names', 'João, Maria')
-- ('preferred_meeting_time', '14:00')
-- ('communication_style', 'direct')
-- ('hobbies', 'corrida, leitura, videogame')
```

**Vantagens:**
- Ultra simples: apenas 5 campos essenciais
- Zero overhead cognitivo
- UNIQUE constraint evita duplicação
- Fácil de implementar, manter e escalar
- Filosofia: menos é mais

**Desvantagens:**
- Sem rastreamento de origem (manual vs IA)
- Parsing necessário para valores complexos (mas isso é OK)

#### Opção 2: Campo JSONB no MySQL 8
```sql
ALTER TABLE persons ADD COLUMN attributes JSON;

-- Exemplo de estrutura:
{
  "personal": {
    "has_children": true,
    "children": ["João", "Maria"],
    "pets": ["Max (dog)"],
    "hobbies": ["running", "reading"]
  },
  "work_preferences": {
    "meeting_time": "afternoon",
    "communication_style": "direct",
    "feedback_preference": "written"
  },
  "ai_insights": {
    "strengths": ["technical", "mentoring"],
    "growth_areas": ["delegation"],
    "last_updated": "2024-01-15T10:30:00Z"
  }
}
```

**Vantagens:**
- Performance para leitura
- Estrutura mais natural
- Suporte nativo MySQL 8

**Desvantagens:**
- Menos flexível para queries
- Histórico de mudanças mais complexo
- Validação mais difícil

### 2.2 Recomendação: Opção 1 (Tabela de Atributos Simplificada)
Recomendo a tabela de atributos simplificada por:
- **Simplicidade**: Tudo como string, sem complexidade desnecessária
- **Filosofia KISS**: IA só adiciona quando tem 100% de certeza
- **Pragmatismo**: Começar simples e evoluir se necessário
- **Manutenção**: Código mais limpo e fácil de debugar

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
```sql
-- Migration simples
CREATE TABLE person_attributes (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    person_id BIGINT NOT NULL,
    attribute_key VARCHAR(100) NOT NULL,
    attribute_value TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (person_id) REFERENCES persons(id),
    UNIQUE KEY unique_person_attribute (person_id, attribute_key),
    INDEX idx_person (person_id)
);
```

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

## Conclusão Simplificada

**Filosofia**: Começar ultra simples e iterar rapidamente.

- ✅ Uma pessoa por vez
- ✅ Sem cache complexo
- ✅ Sem vector stores
- ✅ Apenas OpenAI
- ✅ Custo baixo
- ✅ Implementação em 1 semana

O sistema vai aprender e melhorar conforme os usuários usam, sem complexidade desnecessária.