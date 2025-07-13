# LeaderPro

> **Torne-se um líder mais inteligente**

Uma plataforma de IA que amplifica sua inteligência de liderança, mantendo memória perfeita de cada interação com seu time e sugerindo ações contextuais para transformar você em um líder excepcional.

## 🎯 O Problema

70% dos líderes técnicos são promovidos sem treinamento formal em gestão de pessoas, resultando em:

- **1:1s ineficazes**: Reuniões sem pauta que viram "status reports"
- **Memória limitada**: "O que mesmo o João fez de bom em março?"
- **Performance reviews superficiais**: Viés de recência e falta de dados concretos
- **Insegurança na liderança**: Cada situação nova gera ansiedade

## 💡 Nossa Solução

### IA Contextual Multidimensional
Não apenas lembramos que "João tem filhos" - nossa IA combina:
- **Dados pessoais** + **temporais** + **geográficos** + **histórico**
- Resultado: "João deve estar com as crianças em casa nas férias de julho, pergunte como está sendo conciliar trabalho e família"

### Funcionalidades Core

🤖 **Coach de 1:1s com IA**
- Sugestões de perguntas personalizadas baseadas no contexto completo
- Registro estruturado de respostas
- Continuidade inteligente entre reuniões

📝 **Sistema de Anotações Contínuas**
- Capture feedbacks e conquistas do dia a dia
- Timeline visual de cada pessoa do time
- Conexões automáticas entre eventos

🧠 **Memória Contextual**
- IA que conecta pontos: "Maria melhorou comunicação desde o feedback de março"
- Lembretes inteligentes: "Faz 3 semanas sem 1:1 com Carlos"
- Insights temporais: "Equipe trabalhou extra ontem, reconheça o esforço"

🔗 **Sistema de Menções e Feedback Cruzado**
- Use `@nome` durante 1:1s para referenciar outras pessoas
- Exemplo: "Maria disse que @João é muito atencioso"
- Automaticamente cria feedback no perfil do João
- Constrói rede de percepções e relacionamentos do time

📊 **Performance Reviews Facilitadas**
- Compilação automática dos dados do período
- Exemplos concretos para cada avaliação
- Economize 2+ horas por pessoa

## 🌐 Protótipo Online

**[🔗 Acesse o protótipo funcionando](https://diegoclair.github.io/leaderpro/)**

O protótipo atual demonstra:
- ✅ Interface completa com design moderno
- ✅ Sistema de multi-empresas 
- ✅ Gestão de perfis de pessoas
- ✅ Sistema de @menções com feedback cruzado
- ✅ Histórico de 1:1s e anotações
- ✅ Dados mockados para demonstração

> **Nota**: Este é um protótipo frontend-only com dados simulados. A integração com IA e backend está no roadmap.

## 🚀 Visão de Produto

### MVP (Atual)
Coach de 1:1s com IA contextual

### Roadmap Futuro
- Gestão de projetos com IA
- Analytics de team performance  
- Feedback 360º automatizado
- Integração com sistemas de RH
- Insights de produtividade e bem-estar

## 🏗️ Arquitetura Técnica

### Backend
- **Go (Golang)** - Performance e simplicidade
- **PostgreSQL** - Dados relacionais complexos
- **Redis** - Cache e sessões
- **GORM** - ORM para produtividade

### Frontend
- **Next.js 14** - React com App Router
- **TailwindCSS** - Design system rápido
- **shadcn/ui** - Componentes modernos
- **Zustand** - Estado global simples

### IA
- **OpenAI GPT-4** ou **Claude API** - LLM principal
- **Vector Database** - Memória contextual (Pinecone/Weaviate)
- **APIs externas** - Calendários, feriados, dados contextuais

## 🔧 Melhorias Técnicas Recentes

### ✅ Sistema de Gênero Completo
**Implementação full-stack para texto contextualizado em português:**

**Backend (Go):**
- **Migração** `000003_add_gender_to_person.sql` - Campo ENUM('male', 'female', 'other')
- **Entity** atualizada com campo `Gender *string`
- **ViewModels** com validação `oneof=male female other`
- **Repository** com queries e mappings atualizados

**Frontend (Next.js/TypeScript):**
- **Tipos** atualizados em `/lib/types/index.ts`
- **Formulário** PersonModal com Select de gênero (Masculino/Feminino/Outro)
- **Utilitários** `/lib/utils/gender.ts` para texto contextual português
- **Store** peopleStore.ts com mapeamento do campo gender

**Resultado:** "Sabrina foi **mencionada**" vs "João foi **mencionado**" vs fallback "mencionado(a)"

### ✅ Timeline e Categorização Corrigida
**Problema:** Timeline mostrava "Anotação" para tudo
**Solução:** Correção da estrutura de dados API vs Frontend

- **Interface corrigida** em PersonTimeline.tsx (`type` vs `source_type`)
- **Categorização adequada:** "Reunião 1:1", "Feedback", "Observação"
- **Componentes atualizados** PersonMentions.tsx e PersonTimeline.tsx
- **API mapping** correto entre backend e frontend

### ✅ Constantes Centralizadas (`/lib/constants/`)
**Eliminação de duplicação de código e magic strings:**

- **📝 `/notes.ts`** - Source types, feedback types/categories, labels e cores
- **🏢 `/company.ts`** - Padrões brasileiros SEBRAE/IBGE para tamanho empresa  
- **📱 `/api.ts`** - Endpoints centralizados para consistência
- **🔤 `/messages.ts`** - Mensagens de erro/sucesso padronizadas
- **✅ `/validation.ts`** - Regras de validação compartilhadas

**Helper functions:** `getNoteSourceTypeLabel()`, `getFeedbackTypeColor()`, etc.

### ✅ Componentes Compartilhados (`/components/ui/`)
**Reutilização e consistência visual:**

- **🔄 `LoadingSpinner`** - Eliminando ~10 duplicações diferentes
- **🎨 `AppLogo`** - Logo centralizada (~4 duplicações eliminadas)
- **📱 `PhoneInput`** - Máscara brasileira (+55) compartilhada
- **🔒 `PasswordInput`** - Input com toggle visibilidade olho/olho-riscado
- **📝 `MentionsInputComponent`** - Sistema @mentions com react-mentions
- **🔘 Select components** - Tratamento correto de valores vazios

### ✅ Segurança e Performance
**Padrões obrigatórios para evitar bugs de segurança:**

- **🔐 Storage Manager** - `storageManager.set/get/clearAll()` previne vazamento entre usuários
- **🚀 API Client centralizado** - `apiClient.authGet/authPost()` com token refresh automático
- **📊 Timeline otimizada** - Separação correta Historical vs Mentions
- **🎯 Mentions personalizadas** - Nome + gênero da pessoa vs genérico "Você"
- **🛡️ Validação robusta** - Select constraints, enum validation, type safety

## 🎨 Identidade da Marca

### Posicionamento
**"Seu coach de liderança pessoal com memória perfeita"**

### Tom de Voz
- **Profissional** mas acessível
- **Inteligente** sem ser arrogante  
- **Prático** e orientado a resultados
- **Empático** com as dores do líder

### Diferencial Competitivo
1. **IA contextual multidimensional** vs templates genéricos
2. **Memória contínua** vs reuniões isoladas  
3. **Modelo B2C** = dados são seus, leva ao mudar de empresa
4. **Onboarding 5 minutos** vs meses de implementação
5. **ROI mensurável** = 2h economizadas por review

## 💰 Modelo de Negócio

**B2C - Líder paga individualmente**

- **Básico**: R$ 29,90/mês (limite de tokens)
- **Padrão**: R$ 49,90/mês (uso moderado) 
- **Ilimitado**: R$ 79,90/mês (sem limites)

**Estratégia de lançamento:**
- Trial 30 dias
- Early adopters: 6 meses por R$ 9,90

## 🎯 Mercado

- **TAM**: 4.4M novos tech leads/ano globalmente
- **Concorrentes**: Fellow.app, Culture Amp, 15Five
- **Diferencial**: Único com IA contextual + modelo B2C

## 📁 Documentação do Projeto

### Arquivos Principais
- `/plan/000001-projeto-leaderpro.md` - Plano de negócio completo e análise de mercado
- `/CLAUDE.md` - Diretrizes técnicas gerais e arquitetura do projeto
- `/frontend/README.md` - **Documentação específica do frontend** (arquitetura, componentes, boas práticas)
- `/backend/README.md` - Documentação da API e backend (a ser criado)

### Para Desenvolvedores
**⚠️ Importante**: Sempre leia a documentação específica de cada módulo antes de fazer alterações:
- Frontend: Consulte `/frontend/README.md` para arquitetura, tipos TypeScript, stores Zustand e componentes
- Backend: Consulte `/backend/README.md` para APIs, estrutura Go e banco de dados