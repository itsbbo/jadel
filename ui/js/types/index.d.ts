import { LucideIcon } from 'lucide-react';

export interface BreadcrumbItem {
    title: string;
    href: string;
}

export interface NavGroup {
    title: string;
    items: NavItem[];
}

export interface NavItem {
    title: string;
    href: string;
    icon?: LucideIcon | null;
    isActive?: boolean;
}

export interface SharedData {
    name: string;
    user: User;
    sidebarOpen: boolean;
    [key: string]: unknown;
}

export interface Pagination<T> {
    items: T[];
    perPage: number;
    prevId?: string;
    nextId?: string;
}
