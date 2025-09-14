// Validation utilities for forms

export interface ValidationErrors {
  [key: string]: string;
}

// Username validation
export const validateUsername = (username: string): string | null => {
  if (username === "") {
    return "ユーザー名は必須です";
  }
  if (username.length < 3 || username.length > 20) {
    return "ユーザー名は3-20文字で入力してください";
  }
  const usernameRegex = /^[a-zA-Z0-9_]+$/;
  if (!usernameRegex.test(username)) {
    return "ユーザー名は英数字とアンダースコアのみ使用できます";
  }
  return null;
};

// Email validation
export const validateEmail = (email: string): string | null => {
  if (email === "") {
    return "メールアドレスは必須です";
  }
  const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
  if (!emailRegex.test(email)) {
    return "有効なメールアドレスを入力してください";
  }
  return null;
};

// Password validation
export const validatePassword = (password: string): string | null => {
  if (password === "") {
    return "パスワードは必須です";
  }
  if (password.length < 8) {
    return "パスワードは8文字以上で入力してください";
  }
  const hasLetter = /[a-zA-Z]/.test(password);
  const hasNumber = /[0-9]/.test(password);
  if (!hasLetter || !hasNumber) {
    return "パスワードは英数字の両方を含む必要があります";
  }
  return null;
};

// Confirm password validation
export const validateConfirmPassword = (
  password: string,
  confirmPassword: string,
): string | null => {
  if (confirmPassword === "") {
    return "パスワード確認は必須です";
  }
  if (password !== confirmPassword) {
    return "パスワードが一致しません";
  }
  return null;
};

// Registration form validation
export const validateRegistrationForm = (data: {
  username: string;
  email: string;
  password: string;
  confirmPassword: string;
}): ValidationErrors => {
  const errors: ValidationErrors = {};

  const usernameError = validateUsername(data.username);
  if (usernameError) errors.username = usernameError;

  const emailError = validateEmail(data.email);
  if (emailError) errors.email = emailError;

  const passwordError = validatePassword(data.password);
  if (passwordError) errors.password = passwordError;

  const confirmPasswordError = validateConfirmPassword(
    data.password,
    data.confirmPassword,
  );
  if (confirmPasswordError) errors.confirmPassword = confirmPasswordError;

  return errors;
};

// Login form validation
export const validateLoginForm = (data: {
  username: string;
  password: string;
}): ValidationErrors => {
  const errors: ValidationErrors = {};

  if (data.username === "") {
    errors.username = "ユーザー名は必須です";
  }

  if (data.password === "") {
    errors.password = "パスワードは必須です";
  }

  return errors;
};

// Profile update validation
export const validateProfileUpdate = (data: {
  username: string;
  email: string;
  current_password: string;
  new_password: string;
  confirm_password: string;
}): { isValid: boolean; errors: ValidationErrors } => {
  const errors: ValidationErrors = {};

  const usernameError = validateUsername(data.username);
  if (usernameError) errors.username = usernameError;

  const emailError = validateEmail(data.email);
  if (emailError) errors.email = emailError;

  // Only validate new password if it's provided
  if (data.new_password) {
    if (data.current_password === "") {
      errors.current_password = "現在のパスワードを入力してください";
    }

    const newPasswordError = validatePassword(data.new_password);
    if (newPasswordError) errors.new_password = newPasswordError;

    const confirmPasswordError = validateConfirmPassword(
      data.new_password,
      data.confirm_password,
    );
    if (confirmPasswordError) errors.confirm_password = confirmPasswordError;
  }

  return {
    isValid: Object.keys(errors).length === 0,
    errors,
  };
};
