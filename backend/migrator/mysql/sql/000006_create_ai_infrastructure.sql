-- ================================================
-- Migration 000006: AI Infrastructure Tables
-- ================================================

-- Table for flexible person attributes
CREATE TABLE IF NOT EXISTS `person_attributes` (
    `id` BIGINT NOT NULL AUTO_INCREMENT,
    `person_id` INT NOT NULL,
    `attribute_key` VARCHAR(100) NOT NULL,
    `attribute_value` TEXT NOT NULL,
    `source` ENUM('manual', 'ai_extracted', 'imported') NOT NULL DEFAULT 'manual',
    `extracted_from_note_id` INT NULL,
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    PRIMARY KEY (`id`),
    UNIQUE KEY `unique_person_attribute` (`person_id`, `attribute_key`),
    INDEX `idx_person` (`person_id`),
    INDEX `idx_source` (`source`),
    INDEX `idx_extracted_from_note` (`extracted_from_note_id`),
    
    CONSTRAINT `fk_person_attributes_person` FOREIGN KEY (`person_id`) REFERENCES `tab_person` (`person_id`) ON DELETE CASCADE,
    CONSTRAINT `fk_person_attributes_note` FOREIGN KEY (`extracted_from_note_id`) REFERENCES `tab_note` (`note_id`) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Table for prompt versioning and history
CREATE TABLE IF NOT EXISTS `ai_prompts` (
    `id` BIGINT NOT NULL AUTO_INCREMENT,
    `type` VARCHAR(50) NOT NULL, -- 'leadership_coach', 'attribute_extraction', 'meeting_suggestions'
    `version` INT NOT NULL,
    `prompt` TEXT NOT NULL,
    `model` VARCHAR(50) NOT NULL, -- 'gpt-4o-mini', 'gpt-4', etc
    `temperature` DECIMAL(3,2) NOT NULL DEFAULT 0.7,
    `max_tokens` INT NOT NULL DEFAULT 2000,
    `is_active` BOOLEAN NOT NULL DEFAULT TRUE,
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `created_by` INT NOT NULL,
    
    PRIMARY KEY (`id`),
    INDEX `idx_type_active` (`type`, `is_active`),
    UNIQUE KEY `unique_type_version` (`type`, `version`),
    
    CONSTRAINT `fk_ai_prompts_user` FOREIGN KEY (`created_by`) REFERENCES `tab_user` (`user_id`) ON DELETE RESTRICT
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Table for AI usage tracking
CREATE TABLE IF NOT EXISTS `ai_usage_tracker` (
    `id` BIGINT NOT NULL AUTO_INCREMENT,
    `user_id` INT NOT NULL,
    `company_id` INT NOT NULL,
    `prompt_id` BIGINT NOT NULL,
    `person_id` INT NULL, -- Person context (if applicable)
    `request_type` VARCHAR(50) NOT NULL, -- 'chat', 'extraction', 'suggestion'
    `tokens_used` INT NOT NULL,
    `cost_usd` DECIMAL(10,6) NOT NULL,
    `response_time_ms` INT NOT NULL,
    `feedback` ENUM('helpful', 'not_helpful', 'neutral') NULL,
    `feedback_comment` TEXT NULL,
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    
    PRIMARY KEY (`id`),
    INDEX `idx_user_date` (`user_id`, `created_at`),
    INDEX `idx_company_date` (`company_id`, `created_at`),
    INDEX `idx_feedback` (`feedback`),
    INDEX `idx_request_type` (`request_type`),
    
    CONSTRAINT `fk_ai_usage_user` FOREIGN KEY (`user_id`) REFERENCES `tab_user` (`user_id`) ON DELETE CASCADE,
    CONSTRAINT `fk_ai_usage_company` FOREIGN KEY (`company_id`) REFERENCES `tab_company` (`company_id`) ON DELETE CASCADE,
    CONSTRAINT `fk_ai_usage_prompt` FOREIGN KEY (`prompt_id`) REFERENCES `ai_prompts` (`id`) ON DELETE RESTRICT,
    CONSTRAINT `fk_ai_usage_person` FOREIGN KEY (`person_id`) REFERENCES `tab_person` (`person_id`) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Table for conversation content (separated to facilitate cleanup)
CREATE TABLE IF NOT EXISTS `ai_conversations` (
    `id` BIGINT NOT NULL AUTO_INCREMENT,
    `usage_id` BIGINT NOT NULL,
    `user_message` TEXT NOT NULL,
    `ai_response` TEXT NOT NULL,
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `expires_at` TIMESTAMP NOT NULL DEFAULT (CURRENT_TIMESTAMP + INTERVAL 180 DAY), -- 180 days default
    
    PRIMARY KEY (`id`),
    INDEX `idx_usage` (`usage_id`),
    INDEX `idx_expires` (`expires_at`),
    
    CONSTRAINT `fk_ai_conversations_usage` FOREIGN KEY (`usage_id`) REFERENCES `ai_usage_tracker` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Insert initial system prompts
INSERT INTO `ai_prompts` (`type`, `version`, `prompt`, `model`, `temperature`, `max_tokens`, `is_active`, `created_by`) VALUES
(
    'leadership_coach',
    1,
    'Você é um coach especialista em gestão de pessoas e liderança tecnológica com 20 anos de experiência. Você ajuda líderes tech a se tornarem melhores gestores através de conselhos práticos e personalizados.

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
- Com exemplos concretos quando relevante
- Em português brasileiro',
    'gpt-4o-mini',
    0.7,
    2000,
    TRUE,
    1
),
(
    'attribute_extraction',
    1,
    'Analise as notas fornecidas e extraia APENAS informações que você tem 100% de certeza sobre a pessoa mencionada.

Retorne um JSON com pares chave-valor simples. Use apenas estas chaves permitidas:
- has_children: "true" ou "false"
- children_names: "Nome1, Nome2"
- hobbies: "hobby1, hobby2"
- communication_style: "direct", "diplomatic", "informal"
- preferred_meeting_time: "morning", "afternoon", "evening"
- feedback_preference: "written", "verbal", "immediate"
- personality_traits: "trait1, trait2"
- technical_interests: "area1, area2"
- career_goals: "objetivo mencionado"
- work_challenges: "desafio mencionado"

Exemplo de resposta:
{"has_children": "true", "children_names": "João, Maria", "hobbies": "corrida, leitura"}

Se não tiver certeza absoluta sobre algo, NÃO inclua no JSON. Prefira não extrair do que extrair informação incorreta.',
    'gpt-4o-mini',
    0.3,
    1000,
    TRUE,
    1
);