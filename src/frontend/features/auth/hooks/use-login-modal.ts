"use client";

import { useCallback, useSyncExternalStore } from "react";

type LoginModalStore = {
  isOpen: boolean;
  listeners: Set<() => void>;
};

const store: LoginModalStore = {
  isOpen: false,
  listeners: new Set(),
};

const notifyListeners = () => {
  store.listeners.forEach((listener) => listener());
};

const subscribe = (listener: () => void) => {
  store.listeners.add(listener);
  return () => store.listeners.delete(listener);
};

const getSnapshot = () => store.isOpen;

const getServerSnapshot = () => false;

const open = () => {
  if (!store.isOpen) {
    store.isOpen = true;
    notifyListeners();
  }
};

const close = () => {
  if (store.isOpen) {
    store.isOpen = false;
    notifyListeners();
  }
};

export const useLoginModal = () => {
  const isOpen = useSyncExternalStore(subscribe, getSnapshot, getServerSnapshot);

  const onOpenChange = useCallback((isOpen: boolean) => {
    if (isOpen) {
      open();
    } else {
      close();
    }
  }, []);

  return {
    close,
    isOpen,
    onOpenChange,
    open,
  };
};
