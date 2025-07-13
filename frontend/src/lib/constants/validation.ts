// Validation rules
export const VALIDATION_RULES = {
  // Password
  MIN_PASSWORD_LENGTH: 8,
  MAX_PASSWORD_LENGTH: 128,
  
  // Name
  MIN_NAME_LENGTH: 2,
  MAX_NAME_LENGTH: 100,
  
  // Email
  EMAIL_REGEX: /^[^\s@]+@[^\s@]+\.[^\s@]+$/,
  
  // Phone (Brazilian format)
  PHONE_REGEX: /^\(?[1-9]{2}\)?\s?9?\d{4}-?\d{4}$/,
  
  // Company name
  MIN_COMPANY_NAME_LENGTH: 2,
  MAX_COMPANY_NAME_LENGTH: 100,
  
  // Role/Position
  MIN_ROLE_LENGTH: 2,
  MAX_ROLE_LENGTH: 50,
  
  // Notes/Content
  MIN_NOTE_LENGTH: 1,
  MAX_NOTE_LENGTH: 5000,
} as const

// Validation functions
export const VALIDATORS = {
  email: (email: string) => VALIDATION_RULES.EMAIL_REGEX.test(email),
  
  password: (password: string) => 
    password.length >= VALIDATION_RULES.MIN_PASSWORD_LENGTH && 
    password.length <= VALIDATION_RULES.MAX_PASSWORD_LENGTH,
  
  phone: (phone: string) => 
    !phone || VALIDATION_RULES.PHONE_REGEX.test(phone.replace(/\s/g, '')),
  
  name: (name: string) => 
    name.length >= VALIDATION_RULES.MIN_NAME_LENGTH && 
    name.length <= VALIDATION_RULES.MAX_NAME_LENGTH,
  
  companyName: (name: string) => 
    name.length >= VALIDATION_RULES.MIN_COMPANY_NAME_LENGTH && 
    name.length <= VALIDATION_RULES.MAX_COMPANY_NAME_LENGTH,
  
  role: (role: string) => 
    role.length >= VALIDATION_RULES.MIN_ROLE_LENGTH && 
    role.length <= VALIDATION_RULES.MAX_ROLE_LENGTH,
  
  note: (content: string) => 
    content.length >= VALIDATION_RULES.MIN_NOTE_LENGTH && 
    content.length <= VALIDATION_RULES.MAX_NOTE_LENGTH,
  
  confirmPassword: (password: string, confirmPassword: string) => 
    password === confirmPassword,
} as const

// Form field validation
export function validateField(field: string, value: string, extraValue?: string) {
  switch (field) {
    case 'email':
      return VALIDATORS.email(value)
    case 'password':
      return VALIDATORS.password(value)
    case 'confirmPassword':
      return VALIDATORS.confirmPassword(extraValue || '', value)
    case 'phone':
      return VALIDATORS.phone(value)
    case 'name':
      return VALIDATORS.name(value)
    case 'companyName':
      return VALIDATORS.companyName(value)
    case 'role':
      return VALIDATORS.role(value)
    case 'note':
      return VALIDATORS.note(value)
    default:
      return true
  }
}