import {create} from 'zustand';

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
            //const response = await fetch('/api/admin/exists');
            //const data = await response.json();
            // throw new Error("Mock error for testing"); // Mock error for testing
            console.log("calling backend API to check if admin exists")
            var data = {exists: false}; // Mock response for testing
            set({doesAdminExist: data.exists});
        } catch (error) {
            console.error('Error checking if admin exists:', error);
            set({error: "Failed to check if admin exists. Please try again later."});
        }
    },
}));