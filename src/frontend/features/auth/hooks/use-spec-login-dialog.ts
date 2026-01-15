"use client";

import { useSyncExternalStore } from "react";

type SpecLoginDialogStore = {
  isOpen: boolean;
  listeners: Set<() => void>;
};

const store: SpecLoginDialogStore = {
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

export const useSpecLoginDialog = () => {
  const isOpen = useSyncExternalStore(subscribe, getSnapshot, getServerSnapshot);

  const onOpenChange = (isOpen: boolean) => {
    if (isOpen) {
      open();
    } else {
      close();
    }
  };

  return {
    close,
    isOpen,
    onOpenChange,
    open,
  };
};
