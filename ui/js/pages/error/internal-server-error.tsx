'use client';

import { Button } from '@/components/shadcn/button';
import { Card, CardContent, CardFooter, CardHeader } from '@/components/shadcn/card';
import { Head, Link, router } from '@inertiajs/react';
import { AlertCircle, Home, RefreshCw } from 'lucide-react';

export default function InternalServerError() {
    const reset = () => {
        router.reload();
    };

    return (
        <div className="bg-background flex min-h-screen items-center justify-center p-4">
            <Head title="Internal Server Error" />
            <Card className="mx-auto max-w-md text-center shadow-lg">
                <CardHeader>
                    <div className="mx-auto mb-4 flex h-20 w-20 items-center justify-center rounded-full bg-red-50">
                        <AlertCircle className="h-10 w-10 text-red-500" />
                    </div>
                    <h1 className="text-2xl font-bold tracking-tight">Something went wrong</h1>
                </CardHeader>
                <CardContent>
                    <p className="text-muted-foreground">
                        We've encountered a critical error. Our team has been notified and is working to resolve the issue.
                    </p>
                </CardContent>
                <CardFooter className="flex flex-col gap-2 sm:flex-row sm:justify-center">
                    <Button onClick={reset} variant="default" className="w-full sm:w-auto">
                        <RefreshCw className="mr-2 h-4 w-4" />
                        Try again
                    </Button>
                    <Button asChild variant="outline" className="w-full sm:w-auto">
                        <Link href="/">
                            <Home className="mr-2 h-4 w-4" />
                            Back to home
                        </Link>
                    </Button>
                </CardFooter>
            </Card>
        </div>
    );
}
