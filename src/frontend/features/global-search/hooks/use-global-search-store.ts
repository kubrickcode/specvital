"use client";

import { useSyncExternalStore } from "react";

import type { GlobalSearchState } from "../types";

type Listener = () => void;

// Must be cached to avoid infinite loop in useSyncExternalStore
const SERVER_STATE: GlobalSearchState = { isOpen: false };

const createGlobalSearchStore = () => {
  let state: GlobalSearchState = { isOpen: false };
  const listeners = new Set<Listener>();

  const getState = () => state;
  const getServerState = () => SERVER_STATE;

  const subscribe = (listener: Listener) => {
    listeners.add(listener);
    return () => listeners.delete(listener);
  };

  const notify = () => {
    listeners.forEach((listener) => listener());
  };

  const open = () => {
    state = { ...state, isOpen: true };
    notify();
  };

  const close = () => {
    state = { ...state, isOpen: false };
    notify();
  };

  const toggle = () => {
    state = { ...state, isOpen: !state.isOpen };
    notify();
  };

  return { close, getServerState, getState, open, subscribe, toggle };
};

export const globalSearchStore = createGlobalSearchStore();

export const useGlobalSearchStore = () => {
  const state = useSyncExternalStore(
    globalSearchStore.subscribe,
    globalSearchStore.getState,
    globalSearchStore.getServerState
  );

  return {
    ...state,
    close: globalSearchStore.close,
    open: globalSearchStore.open,
    toggle: globalSearchStore.toggle,
  };
};
