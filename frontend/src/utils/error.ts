// エラーハンドリング関連のユーティリティ関数

import { ApiError } from "@/services/api";

/**
 * APIエラーを適切なユーザー向けメッセージに変換
 */
export const getErrorMessage = (
  error: unknown,
  defaultMessage: string,
): string => {
  if (error instanceof ApiError) {
    return error.message;
  }
  if (error instanceof Error) {
    return error.message;
  }
  return defaultMessage;
};

/**
 * エラーがリトライ可能かどうかを判定
 */
export const isRetryableError = (error: unknown): boolean => {
  return ApiError.isTemporaryError(error);
};

/**
 * エラーが認証エラーかどうかを判定
 */
export const isAuthError = (error: unknown): boolean => {
  return ApiError.isAuthError(error);
};

/**
 * エラーがバリデーションエラーかどうかを判定
 */
export const isValidationError = (error: unknown): boolean => {
  return ApiError.isValidationError(error);
};

/**
 * APIエラーからフィールドエラーを抽出
 */
export const extractFieldErrors = (error: unknown): Record<string, string> => {
  if (error instanceof ApiError && error.hasFieldErrors()) {
    return error.fieldErrors || {};
  }
  return {};
};

/**
 * 共通エラーハンドラー
 */
export const handleApiError = (
  error: unknown,
  setError: (message: string) => void,
  fallbackMessage: string = "エラーが発生しました",
): void => {
  console.error("API Error:", error);
  const message = getErrorMessage(error, fallbackMessage);
  setError(message);
};
