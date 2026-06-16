import { create } from "zustand";

export type WorkspaceSection = "audit-log" | "inbox";

interface UiStore {
  currentSection: WorkspaceSection;
  setCurrentSection: (section: WorkspaceSection) => void;
}

export const useUiStore = create<UiStore>((set) => ({
  currentSection: "audit-log",
  setCurrentSection: (section) => set({ currentSection: section }),
}));
