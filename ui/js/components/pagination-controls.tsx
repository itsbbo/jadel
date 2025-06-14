import { Button } from '@/components/shadcn/button';
import { router } from '@inertiajs/react';
import { ChevronLeft, ChevronRight } from 'lucide-react';

interface Props {
    prevPageURL?: string;
    nextPageURL?: string;
}

function isNullishString(value: unknown): boolean {
    return value === null || value === undefined || value === '' || value === 'unknown';
}

export default function PaginationControls({ prevPageURL, nextPageURL }: Props) {
    const handlePrevClick = () => {
        router.get(prevPageURL ?? '/dashboard');
    };

    const handleNextClick = () => {
        router.get(nextPageURL ?? '/dashboard');
    };

    return (
        <div className="flex items-center justify-between">
            <div className="flex items-center space-x-2">
                <Button variant="outline" size="sm" onClick={handlePrevClick} disabled={isNullishString(prevPageURL)}>
                    <ChevronLeft className="h-4 w-4" />
                    Previous
                </Button>
                <Button variant="outline" size="sm" onClick={handleNextClick} disabled={isNullishString(nextPageURL)}>
                    Next
                    <ChevronRight className="h-4 w-4" />
                </Button>
            </div>
        </div>
    );
}
