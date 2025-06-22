export interface User {
    id: string;
    name: string;
    email: string;
    avatar?: string;
    createdAt: string;
    updatedAt: string;
}

export interface Project {
    id: string;
    name: string;
    description?: string;
    createdAt: string;
    updatedAt: string;
}

export interface Server {
    id: string;
    name: string;
    description?: string;
    createdAt: string;
    updatedAt: string;
}

export interface Environment {
    id: string;
    name: string;
    createdAt: string;
    updatedAt: string;
    project: Project;
    applications: Application[];
    databases: Database[];
}

export interface Application {
    id: string;
    environmentId: string;
    name: string;
    description?: string;
    buildTool: string;
    domain: string;
    enableHttps: boolean;
    preDeploymentScript?: string;
    postDeploymentScript?: string;
    portMappings?: Record<string, string>;
    metadata?: Record<string, any>;
    createdAt: string;
    updatedAt: string;
}

export interface Database {
    id: string;
    environmentId: string;
    database_type: string;
    name: string;
    description?: string;
    image: string;
    username: string;
    password?: string;
    portMappings?: Record<string, string>;
    customConfig?: Record<string, any>;
    metadata?: any;
    createdAt: string;
    updatedAt: string;
}

export interface PrivateKey {
    id: string;
    name: string;
    description?: string;
    publicKey: string;
    privateKey: string;
    createdAt: string;
    updatedAt: string;
}
