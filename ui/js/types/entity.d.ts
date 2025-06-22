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
    project: Project;
    applications: Application[];
    databases: Database[];
}

export interface Application {
    id: string;
    environment_id: string;
    name: string;
    description?: string;
    build_tool: string;
    domain: string;
    enable_https: boolean;
    pre_deployment_script?: string;
    post_deployment_script?: string;
    port_mappings?: Record<string, string>;
    metadata?: Record<string, any>;
    created_at: string;
    updated_at: string;
}

export interface Database {
    id: string;
    environment_id: string;
    database_type: string;
    name: string;
    description?: string;
    image: string;
    username: string;
    password?: string;
    port_mappings?: Record<string, string>;
    custom_config?: Record<string, any>;
    metadata?: any;
    created_at: string;
    updated_at: string;
}

export interface PrivateKey {
    id: string;
    name: string;
    description?: string;
    public_key: string;
    private_key: string;
    created_at: string;
    updated_at: string;
}
