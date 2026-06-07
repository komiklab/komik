import {create} from 'zustand';
import { getAdmin } from '../api/komik';

interface AdminStore {
    doesAdminExist: boolean;
    checkIfAdminExists: () => Promise<void>;
    error: string | null;
}

export const useAdminStore = create<AdminStore>((set) => ({
    doesAdminExist: false,
    error: null,
    checkIfAdminExists: async () => {
        try {
            const response = await getAdmin();
            if (response.status === 200) {
                set({doesAdminExist: response.data.exists ?? false});
            } else {
                throw new Error(`Unexpected status: ${response.status}`);
            }
        } catch (error) {
            console.error('Error checking if admin exists:', error);
            set({error: "Failed to check if admin exists. Please try again later."});
        }
    },
}));