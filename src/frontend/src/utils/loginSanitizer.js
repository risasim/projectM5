export class loginSanitizer {
  // allows only letters, numbers, dashes etc
  static USERNAME_PATTERN = /^[a-zA-Z0-9_-]{3,20}$/;
  static PASSWORD_MIN = 8;  
  static PASSWORD_MAX = 50;
  // potentially dangerous characters
  static DANGEROUS = /[<>&"`'$(){}[\]]/;

  // escapes html special characters
  // converts <, >, "" etc into safe html entities
  static escapeHtml(text) {
    if (typeof text !== 'string') return '';
    return text
      .replace(/&/g, '&amp;')
      .replace(/</g, '&lt;')
      .replace(/>/g, '&gt;')
      .replace(/"/g, '&quot;')
      .replace(/'/g, '&#39;');
  }

  // sanitizes username input by trimming spaces, checking valid chars, rejecting unsafe chars
  static sanitizeUsername(input) {
    if (!input) return { valid: false, error: 'Username required' };
    const trimmed = input.trim();
    if (!this.USERNAME_PATTERN.test(trimmed))
      return { valid: false, error: 'Invalid characters in username' };
    if (this.DANGEROUS.test(trimmed))
      return { valid: false, error: 'Unsafe characters detected' };
    return { valid: true, value: trimmed, safeHtml: this.escapeHtml(trimmed) };
  }

  // sanitizes password input by ensuring a presence of the password and enforces the length limits
  static sanitizePassword(input) {
    if (!input) return { valid: false, error: 'Password required' };
    if (input.length < this.PASSWORD_MIN)
      return { valid: false, error: 'Too short' };
    if (input.length > this.PASSWORD_MAX)
      return { valid: false, error: 'Too long' };
    return { valid: true, value: input };
  }
}

// functions for using them inside the frontend, vue
export const sanitizeUsername = (i) => loginSanitizer.sanitizeUsername(i);
export const sanitizePassword = (i) => loginSanitizer.sanitizePassword(i);
