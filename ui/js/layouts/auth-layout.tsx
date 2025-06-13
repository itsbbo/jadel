import { Head } from '@inertiajs/react';

interface AuthLayoutProps {
    children: React.ReactNode;
    title: string;
    description: string;
}

export default function AuthLayout({ children, title, description }: AuthLayoutProps) {
    return (
        <div className="bg-background flex min-h-screen flex-col items-center justify-center p-4 sm:p-0">
            <Head title={title} />
            <div className="w-full max-w-md space-y-8">
                <div className="text-center">
                    <h1 className="text-2xl font-bold tracking-tight">{title}</h1>
                    <p className="text-muted-foreground mt-2 text-sm">{description}</p>
                </div>
                <div className="bg-card text-card-foreground rounded-lg border p-8 shadow-sm">{children}</div>
            </div>
        </div>
    );
}
