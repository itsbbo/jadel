import { Button } from '@/components/shadcn/button';
import { Card, CardContent, CardFooter, CardHeader } from '@/components/shadcn/card';
import { Link } from '@inertiajs/react';
import { FileQuestion, Home } from 'lucide-react';

export default function NotFound() {
    return (
        <div className="bg-background flex min-h-screen items-center justify-center p-4">
            <Card className="mx-auto max-w-md text-center shadow-lg">
                <CardHeader>
                    <div className="mx-auto mb-4 flex h-20 w-20 items-center justify-center rounded-full bg-blue-50">
                        <FileQuestion className="h-10 w-10 text-blue-500" />
                    </div>
                    <h1 className="text-2xl font-bold tracking-tight">Page Not Found</h1>
                </CardHeader>
                <CardContent>
                    <p className="text-muted-foreground">
                        Sorry, we couldn't find the page you're looking for. The page might have been moved, deleted, or never existed.
                    </p>
                </CardContent>
                <CardFooter className="flex justify-center">
                    <Button asChild variant="default">
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
