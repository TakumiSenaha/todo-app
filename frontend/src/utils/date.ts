// 日付関連のユーティリティ関数

/**
 * 日付文字列をローカル形式でフォーマット
 */
export const formatDate = (dateString: string): string => {
  const date = new Date(dateString);
  return date.toLocaleDateString("ja-JP", {
    year: "numeric",
    month: "short",
    day: "numeric",
  });
};

/**
 * 日付文字列を相対形式でフォーマット（例: "2日前", "明日"）
 */
export const formatRelativeDate = (dateString: string): string => {
  const date = new Date(dateString);
  const now = new Date();
  const diffTime = date.getTime() - now.getTime();
  const diffDays = Math.ceil(diffTime / (1000 * 60 * 60 * 24));

  if (diffDays === 0) {
    return "今日";
  } else if (diffDays === 1) {
    return "明日";
  } else if (diffDays === -1) {
    return "昨日";
  } else if (diffDays > 1) {
    return `${diffDays}日後`;
  } else {
    return `${Math.abs(diffDays)}日前`;
  }
};

/**
 * 今日の日付をYYYY-MM-DD形式で取得
 */
export const getTodayString = (): string => {
  const today = new Date();
  return today.toISOString().split("T")[0];
};

/**
 * 日付が有効かどうかを判定
 */
export const isValidDate = (dateString: string): boolean => {
  const date = new Date(dateString);
  return !isNaN(date.getTime());
};
