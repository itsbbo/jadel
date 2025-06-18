import { cn } from '@/lib/utils';

export default function Heading({ title, description, className }: { title: string; description?: string; className?: string }) {
    return (
        <div className={cn('mb-8 space-y-0.5', className)}>
            <h2 className="text-xl font-semibold tracking-tight md:text-2xl">{title}</h2>
            {description && <p className="text-muted-foreground md:text-md text-sm">{description}</p>}
        </div>
    );
}
