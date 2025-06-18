export interface User {
    id: string;
    name: string;
    email: string;
    avatar?: string;
    created_at: string;
    updated_at: string;
}

export interface Project {
    id: string;
    name: string;
    description?: string;
    created_at: string;
    updated_at: string;
}

export interface Server {
    id: string;
    name: string;
    description?: string;
    created_at: string;
    updated_at: string;
}

export interface Environment {
    id: string;
    name: string;
    created_at: string;
    updated_at: string;
}
