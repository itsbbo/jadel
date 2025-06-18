import { Badge } from '@/components/shadcn/badge';
import { Card, CardContent } from '@/components/shadcn/card';
import { LucideIcon } from 'lucide-react';

interface ResourceCardProps {
    title: string;
    description: string;
    icon: LucideIcon;
    category?: string;
    isPopular?: boolean;
    onClick?: () => void;
}

export function ResourceCard({ title, description, icon: Icon, category, isPopular = false, onClick }: ResourceCardProps) {
    return (
        <Card onClick={onClick}>
            <CardContent>
                <div className="flex items-start gap-4">
                    <div className="flex-shrink-0">
                        <div className="bg-muted flex h-12 w-12 items-center justify-center rounded-lg">
                            <Icon className="text-muted-foreground h-6 w-6" />
                        </div>
                    </div>
                    <div className="min-w-0 flex-1">
                        <div className="mb-2 flex items-center gap-2">
                            <h3 className="font-medium">{title}</h3>
                            {isPopular && <Badge className="bg-primary/10 text-primary border-primary/20 text-xs">Popular</Badge>}
                        </div>
                        <p className="text-muted-foreground line-clamp-2 text-sm">{description}</p>
                        {category && (
                            <div className="mt-3">
                                <Badge variant="outline" className="text-xs">
                                    {category}
                                </Badge>
                            </div>
                        )}
                    </div>
                </div>
            </CardContent>
        </Card>
    );
}
