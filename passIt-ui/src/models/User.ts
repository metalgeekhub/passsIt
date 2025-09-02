export type User = {
    id: string;
    keycloack_id: string;
    created_at: string;
    updated_at: string;
    username: string;
    email: string;
    first_name: string;
    last_name: string;
    dob: string;
    phone_number: string; // <-- should be string to match backend
    address: string;
    is_active: boolean;
    is_admin: boolean;
};