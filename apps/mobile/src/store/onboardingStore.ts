import { create } from 'zustand';
import type { BudgetTier, Goal } from '@iron-and-spice/shared';

type OnboardingDraft = {
  goal: Goal | null;
  budgetTier: BudgetTier | null;
  region: string;
  restrictions: string[];
  set: (patch: Partial<Omit<OnboardingDraft, 'set' | 'reset'>>) => void;
  reset: () => void;
};

export const useOnboardingStore = create<OnboardingDraft>((set) => ({
  goal: null,
  budgetTier: null,
  region: '',
  restrictions: [],
  set: (patch) => set(patch),
  reset: () => set({ goal: null, budgetTier: null, region: '', restrictions: [] }),
}));
