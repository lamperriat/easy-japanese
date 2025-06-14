import { GetYoudao } from "./Youdao";

const DICT_CHOSEN = "youdao";

const CACHE_PREFIX = "easyjp_hover_preview_cache_"
const MAX_CACHE_SIZE = 5 * 1024 * 1024; 
const DEFAULT_TTL = 3 * 60 * 60 * 1000; // 3 hours

export const getCache = (key) => {
  if (MAX_CACHE_SIZE === 0) return null; // Cache disabled
  const fullKey = CACHE_PREFIX + key;
  const itemStr = localStorage.getItem(fullKey);

  if (!itemStr) return null;

  const item = JSON.parse(itemStr);
  const isExpired = Date.now() > item.timestamp + DEFAULT_TTL;

  if (isExpired) {
    localStorage.removeItem(fullKey);
    return null;
  }

  return item.data;
};

const clearOldestCacheItem = () => {
  let oldestKey = null;
  let oldestTimestamp = Infinity;

  for (let i = 0; i < localStorage.length; i++) {
    const key = localStorage.key(i);
    if (key?.startsWith(CACHE_PREFIX)) {
      const itemStr = localStorage.getItem(key);
      if (itemStr) {
        const item = JSON.parse(itemStr);
        if (item.timestamp < oldestTimestamp) {
          oldestTimestamp = item.timestamp;
          oldestKey = key;
        }
      }
    }
  }

  if (oldestKey) localStorage.removeItem(oldestKey);
};

export const setCache = (key, data) => {
  const fullKey = CACHE_PREFIX + key;
  const item = {
    data,
    timestamp: Date.now(),
  };

  // Check cache size and evict oldest if needed
  const currentSize = JSON.stringify(localStorage).length;
  if (currentSize + JSON.stringify(item).length > MAX_CACHE_SIZE) {
    clearOldestCacheItem();
  }

  localStorage.setItem(fullKey, JSON.stringify(item));
};

export default async function DictSearch(word, token) {
  const cachedResult = getCache(word);
  if (cachedResult) {
    return Promise.resolve(cachedResult);
  }
  if (DICT_CHOSEN === "youdao") {
    const result = await GetYoudao(word, token);
    if (result) {
      setCache(word, result);
      return Promise.resolve(result);
    } else {
      console.error("Failed to fetch data from Youdao for word:", word);
      return Promise.reject(new Error("Failed to fetch data from Youdao"));
    }

  } else {
    console.error("No dictionary chosen or unsupported dictionary:", DICT_CHOSEN);
    return Promise.reject(new Error("Unsupported dictionary"));
  }
}