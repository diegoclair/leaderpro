# 🎬 Sistema de Animação Refinado - LeaderPro

Sistema de animações inspirado em **GitHub**, **Linear**, **Stripe** e **Notion**.

## 🎯 Princípios

### ✨ **Sutileza**: Animações funcionais, não decorativas
### ⚡ **Performance**: Apenas `transform`, `opacity` e `shadow`
### 🎛️ **Consistência**: Classes padronizadas reutilizáveis
### 📱 **Acessibilidade**: Estados focus/active bem definidos

## 🏗️ Classes Disponíveis

### **Cards Interativos**

```css
/* Cards gerais - Efeito sutil com lift */
.card-interactive
- Hover: shadow-md + translateY(-2px)
- Focus: ring outline
- Active: translateY(0)
- Timing: 200ms ease-out

/* Cards de pessoas - Estilo GitHub */
.person-card  
- Hover: shadow-lg + translateY(-4px) + border-primary
- Timing: 150ms ease-out
- Border color changes on hover

/* Cards timeline - Estilo Linear */
.timeline-card
- Hover: shadow-md + translateY(-2px) + border fade
- Timing: 150ms ease-out
```

### **Botões e Micro-interações**

```css
/* Botões principais - Estilo Stripe */
.btn-interactive
- Hover: translateY(-2px) + shadow-sm
- Active: translateY(0) + no shadow
- Focus: ring outline
- Timing: 150ms ease-out

/* Micro lift para ícones/pequenos elementos */
.micro-lift
- Hover: translateY(-2px)
- Timing: 100ms ease-out

/* Micro scale para elementos especiais */
.micro-scale  
- Hover: scale(1.05)
- Timing: 100ms ease-out
```

### **Estados e Feedback**

```css
/* Loading states */
.loading-fade     /* Pulse suave */
.loading-slide    /* Transições 300ms */

/* Feedback success */
.success-bounce   /* 2 bounces, 0.6s total */
```

## 📐 Especificações Técnicas

### **Timing Curves**
- **Hover/Focus**: 150-200ms `ease-out`
- **State changes**: 300ms `ease-out`  
- **Micro-interactions**: 100ms `ease-out`
- **Active/Press**: Instantâneo ou 50ms

### **Transform Values**
- **Cards lift**: `-translate-y-0.5` (2px) até `-translate-y-1` (4px)
- **Button lift**: `-translate-y-0.5` (2px)
- **Scale**: `1.05` máximo para elementos pequenos

### **Shadows**
- **Rest**: `shadow-sm`
- **Hover**: `shadow-md` 
- **Elevated**: `shadow-lg`

## 🎨 Implementação

### **Aplicado em:**
✅ **PersonCard**: `.person-card` - Efeito GitHub style  
✅ **Timeline Cards**: `.timeline-card` - Efeito Linear style  
✅ **Botões Header**: `.btn-interactive` - Micro-feedback  
✅ **Filtros**: `.btn-interactive` - Consistência  
✅ **Avatar**: `.micro-lift` - Sutil interação  

### **Performance**
- 🚀 **Hardware accelerated**: Apenas `transform` e `opacity`
- ⚡ **60 FPS**: Transições otimizadas
- 🎯 **Específico**: Evita `transition-all`

## 📚 Referências

### **GitHub**: Hover states sutis, sem scale exagerado
### **Linear**: Lift elegante com translateY
### **Stripe**: Micro-feedback em botões  
### **Notion**: Transições funcionais

## 🎮 Como Usar

```tsx
// Card de pessoa
<Card className="person-card">
  {/* Conteúdo */}
</Card>

// Botão interativo  
<Button className="btn-interactive">
  Ação
</Button>

// Elemento da timeline
<div className="timeline-card">
  {/* Conteúdo */}  
</div>

// Micro-interação
<IconButton className="micro-lift">
  <Icon />
</IconButton>
```

## 🔄 Resultados

✨ **Visual mais refinado** - Inspirado nos melhores sistemas  
⚡ **Performance otimizada** - Apenas propriedades aceleradas por hardware  
🎯 **Consistência total** - Classes padronizadas  
📱 **Acessibilidade** - Estados focus/active bem definidos  
🎨 **Hierarquia clara** - Diferentes níveis de interação  

---

*Sistema implementado em 2025 - Baseado nas melhores práticas da indústria*