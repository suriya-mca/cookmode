import { Platform } from "react-native";
import * as SecureStore from "expo-secure-store";

// expo-secure-store has no web implementation (its web build exports an empty
// object), so web falls back to localStorage and, if that is unavailable
// (private mode, SSR, static rendering), to an in-memory map. Native uses the
// real keychain/keystore via SecureStore, with the same in-memory fallback so
// a storage failure degrades to a logged-in session instead of a crash.
const memory = new Map<string, string>();

function localStorageOrNull(): Storage | null {
  try {
    if (typeof window !== "undefined" && window.localStorage) {
      return window.localStorage;
    }
  } catch {
    // Access to localStorage can throw when cookies/storage are blocked.
  }
  return null;
}

export const secureStorage = {
  async getItem(key: string): Promise<string | null> {
    if (Platform.OS === "web") {
      return localStorageOrNull()?.getItem(key) ?? memory.get(key) ?? null;
    }
    try {
      return await SecureStore.getItemAsync(key);
    } catch {
      return memory.get(key) ?? null;
    }
  },

  async setItem(key: string, value: string): Promise<void> {
    if (Platform.OS === "web") {
      const store = localStorageOrNull();
      if (store) {
        try {
          store.setItem(key, value);
          return;
        } catch {
          // Quota or private mode — fall through to memory.
        }
      }
      memory.set(key, value);
      return;
    }
    try {
      await SecureStore.setItemAsync(key, value);
    } catch {
      memory.set(key, value);
    }
  },

  async removeItem(key: string): Promise<void> {
    memory.delete(key);
    if (Platform.OS === "web") {
      try {
        localStorageOrNull()?.removeItem(key);
      } catch {
        // Nothing else to do; the memory entry is already gone.
      }
      return;
    }
    try {
      await SecureStore.deleteItemAsync(key);
    } catch {
      // Already absent.
    }
  },
};
