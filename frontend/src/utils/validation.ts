// Validation utilities for forms

export interface ValidationErrors {
  [key: string]: string;
}

export interface ValidationResult {
  isValid: boolean;
  errors: ValidationErrors;
}

// Username validation
export function validateUsername(username: string): string | null {
  if (!username.trim()) {
    return "Username is required";
  }

  if (username.length < 3 || username.length > 20) {
    return "Username must be between 3 and 20 characters";
  }

  if (!/^[a-zA-Z0-9_]+$/.test(username)) {
    return "Username can only contain letters, numbers, and underscores";
  }

  return null;
}

// Email validation
export function validateEmail(email: string): string | null {
  if (!email.trim()) {
    return "Email is required";
  }

  if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email)) {
    return "Please enter a valid email address";
  }

  return null;
}

// Password validation
export function validatePassword(password: string): string | null {
  if (!password) {
    return "Password is required";
  }

  if (password.length < 8) {
    return "Password must be at least 8 characters";
  }

  if (!/(?=.*[a-zA-Z])(?=.*\d)/.test(password)) {
    return "Password must contain both letters and numbers";
  }

  return null;
}

// Confirm password validation
export function validateConfirmPassword(
  password: string,
  confirmPassword: string,
): string | null {
  if (password !== confirmPassword) {
    return "Passwords do not match";
  }

  return null;
}

// Profile update validation
export function validateProfileUpdate(data: {
  username: string;
  email: string;
  current_password: string;
  new_password: string;
  confirm_password: string;
}): ValidationResult {
  const errors: ValidationErrors = {};

  // Username validation
  const usernameError = validateUsername(data.username);
  if (usernameError) {
    errors.username = usernameError;
  }

  // Email validation
  const emailError = validateEmail(data.email);
  if (emailError) {
    errors.email = emailError;
  }

  // Password validation (only if new password is provided)
  if (data.new_password) {
    if (!data.current_password) {
      errors.current_password =
        "Current password is required to change password";
    }

    const newPasswordError = validatePassword(data.new_password);
    if (newPasswordError) {
      errors.new_password = newPasswordError;
    }

    const confirmPasswordError = validateConfirmPassword(
      data.new_password,
      data.confirm_password,
    );
    if (confirmPasswordError) {
      errors.confirm_password = confirmPasswordError;
    }
  }

  return {
    isValid: Object.keys(errors).length === 0,
    errors,
  };
}

// Login validation
export function validateLogin(data: {
  username: string;
  password: string;
}): ValidationResult {
  const errors: ValidationErrors = {};

  if (!data.username.trim()) {
    errors.username = "Username is required";
  }

  if (!data.password) {
    errors.password = "Password is required";
  }

  return {
    isValid: Object.keys(errors).length === 0,
    errors,
  };
}

// Registration validation
export function validateRegistration(data: {
  username: string;
  email: string;
  password: string;
  confirm_password: string;
}): ValidationResult {
  const errors: ValidationErrors = {};

  // Username validation
  const usernameError = validateUsername(data.username);
  if (usernameError) {
    errors.username = usernameError;
  }

  // Email validation
  const emailError = validateEmail(data.email);
  if (emailError) {
    errors.email = emailError;
  }

  // Password validation
  const passwordError = validatePassword(data.password);
  if (passwordError) {
    errors.password = passwordError;
  }

  // Confirm password validation
  const confirmPasswordError = validateConfirmPassword(
    data.password,
    data.confirm_password,
  );
  if (confirmPasswordError) {
    errors.confirm_password = confirmPasswordError;
  }

  return {
    isValid: Object.keys(errors).length === 0,
    errors,
  };
}
