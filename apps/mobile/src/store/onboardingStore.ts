import { create } from 'zustand';
import type { BudgetTier, Goal } from '@fittrack/shared';

type OnboardingDraft = {
  goal: Goal | null;
  weightKg: number | null;
  budgetTier: BudgetTier | null;
  region: string;
  restrictions: string[];
  set: (patch: Partial<Omit<OnboardingDraft, 'set' | 'reset'>>) => void;
  reset: () => void;
};

type OnboardingFields = Omit<OnboardingDraft, 'set' | 'reset'>;

const emptyDraft = (): OnboardingFields => ({
  goal: null,
  weightKg: null,
  budgetTier: null,
  region: '',
  restrictions: [],
});

export const useOnboardingStore = create<OnboardingDraft>((set) => ({
  ...emptyDraft(),
  set: (patch) => set(patch),
  reset: () => set(emptyDraft()),
}));
