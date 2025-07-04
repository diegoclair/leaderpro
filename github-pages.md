# 🌐 GitHub Pages - LeaderPro (Temporário para Demonstração)

> **⚠️ IMPORTANTE**: Esta configuração é **temporária** apenas para demonstração do protótipo. Após apresentação, todos os arquivos específicos do GitHub Pages serão removidos.

## 📋 O que foi Adicionado/Modificado

### 🆕 Arquivos Criados (para remoção posterior)
```
.github/workflows/deploy.yml          # GitHub Actions workflow
.nojekyll                            # Evita processamento Jekyll
frontend/next.config.js              # Config Next.js para export estático  
frontend/src/app/profile/[id]/layout.tsx  # generateStaticParams para rotas dinâmicas
github-pages.md                      # Este arquivo de documentação
```

### 🔧 Arquivos Modificados (para reverter depois)
```
README.md                           # Adicionada seção "Protótipo Online"
frontend/eslint.config.mjs          # Relaxadas regras para warnings
frontend/src/app/profile/[id]/page.tsx  # Removido useSearchParams (substituído por window.location)
frontend/src/components/layout/ThemeProvider.tsx  # Simplificado types para build
```

### 📁 Arquivo Original Mantido
```
frontend/next.config.ts              # Configuração original preservada
```

## 🚀 Como Ativar GitHub Pages

### 1. Deploy Inicial
```bash
git add .
git commit -m "Add temporary GitHub Pages setup for demo"
git push origin main
```

### 2. Configurar no GitHub
1. **Settings** > **Pages**
2. **Source**: GitHub Actions  
3. Aguardar workflow concluir (~2 min)
4. Acessar: `https://diegoclair.github.io/leaderpro/`

## 🗑️ Como Remover Completamente

### 1. Desativar GitHub Pages
1. **Settings** > **Pages** 
2. **Source**: None
3. Confirmar desativação

### 2. Reverter Código (1 comando)
```bash
# Remove todos os arquivos adicionados para GitHub Pages
git rm -rf .github .nojekyll frontend/next.config.js frontend/src/app/profile/[id]/layout.tsx github-pages.md

# Reverte modificações nos arquivos existentes  
git checkout HEAD -- README.md frontend/eslint.config.mjs frontend/src/app/profile/[id]/page.tsx frontend/src/components/layout/ThemeProvider.tsx

# Commit da limpeza
git commit -m "Remove GitHub Pages setup - back to development state"
git push origin main
```

### 3. Verificar Estado Limpo
```bash
git status  # Deve mostrar working directory clean
git diff HEAD~1  # Mostra apenas remoções dos arquivos temporários
```

## 🎯 Estado Final (Pós-Remoção)

Após a remoção, o projeto voltará ao estado original:
- ✅ Frontend Next.js normal (sem export estático)
- ✅ Configurações de desenvolvimento mantidas
- ✅ Nenhum arquivo específico do GitHub Pages
- ✅ README.md sem referência ao protótipo online
- ✅ ESLint com regras originais
- ✅ useSearchParams funcionando normalmente

## 📊 Arquivos do Build (em /frontend/out/)

**Estes arquivos são gerados automaticamente e não precisam ser commitados:**
```
out/
├── index.html               # Dashboard principal  
├── profile/1/index.html     # Maria Santos
├── profile/2/index.html     # João Silva  
├── profile/3/index.html     # Pedro Costa
├── profile/4/index.html     # Ana Lima
└── _next/                   # Assets do Next.js
```

## 🔍 Resumo das Modificações

### Essenciais para GitHub Pages
1. **GitHub Actions**: Build e deploy automático
2. **Static Export**: Next.js configurado para gerar HTML estático
3. **Dynamic Routes**: generateStaticParams para `/profile/[id]`
4. **Base Path**: URLs ajustadas para subdiretório do GitHub Pages

### Compatibilidade
1. **useSearchParams**: Substituído por window.location (client-side)
2. **ThemeProvider**: Simplificado para evitar build errors
3. **ESLint**: Warnings ao invés de errors para build

---

**🎯 Objetivo**: Demonstrar funcionalidades completas do LeaderPro para stakeholders antes da implementação do backend real.