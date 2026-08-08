import { create } from "zustand";

export type SettingsSection = "agent" | "hooks";

interface SettingsSectionStore {
  currentSection: SettingsSection;
  setCurrentSection: (section: SettingsSection) => void;
}

export const useSettingsSectionStore = create<SettingsSectionStore>((set) => ({
  currentSection: "agent",
  setCurrentSection: (section) => set({ currentSection: section }),
}));
