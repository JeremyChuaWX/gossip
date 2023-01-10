import create from "zustand";
import { User } from "../models/entities";

interface AuthState {
  user: User | undefined;
  setUser: (curr: User) => void;
}

const useAuthStore = create<AuthState>()((set) => ({
  user: undefined,
  setUser: (curr) => set({ user: curr }),
}));

export default useAuthStore;
