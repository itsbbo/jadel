import { Else, If } from '@/components/condition';
import Heading from '@/components/heading';
import HeadingSmall from '@/components/heading-small';
import { Badge } from '@/components/shadcn/badge';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/shadcn/card';
import AppLayout from '@/layouts/app-layout';
import { type BreadcrumbItem } from '@/types';
import { Project, Server as ServerEntity } from '@/types/entity';
import { Head } from '@inertiajs/react';
import { FolderOpen, Server } from 'lucide-react';

const breadcrumbs: BreadcrumbItem[] = [
    {
        title: 'Dashboard',
        href: '/dashboard',
    },
];

interface Props {
    projects: Project[];
    servers: ServerEntity[];
}

export default function Dashboard({ projects, servers }: Props) {
    return (
        <AppLayout breadcrumbs={breadcrumbs}>
            <Head title="Dashboard" />

            <main className="space-y-16 px-4 py-6">
                <Heading title="Dashboard" description="Your self-hosted infrastructure." />
                <section className="space-y-4">
                    <div className="flex items-center gap-2">
                        <FolderOpen className="h-5 w-5" />
                        <HeadingSmall title="Latest Projects" />
                        <Badge variant="secondary" className="ml-2">
                            {projects?.length ?? 0}
                        </Badge>
                    </div>

                    <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
                        <If condition={projects?.length > 0}>
                            {projects?.map((project) => (
                                <Card key={project.id}>
                                    <CardHeader>
                                        <CardTitle className="text-lg">{project.name}</CardTitle>
                                    </CardHeader>
                                    <CardContent>
                                        <CardDescription className="text-sm leading-relaxed">{project.description}</CardDescription>
                                    </CardContent>
                                </Card>
                            ))}

                            <Else>
                                <small className="text-sm leading-none font-medium">No projects yet</small>
                            </Else>
                        </If>
                    </div>
                </section>

                <section className="space-y-4">
                    <div className="flex items-center gap-2">
                        <Server className="h-5 w-5" />
                        <HeadingSmall title="Latest Servers" />
                        <Badge variant="secondary" className="ml-2">
                            {servers?.length ?? 0}
                        </Badge>
                    </div>
                    <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
                        <If condition={servers?.length > 0}>
                            {servers?.map((server) => (
                                <Card key={server.id} className="transition-shadow hover:shadow-md">
                                    <CardHeader>
                                        <CardTitle className="text-lg">{server.name}</CardTitle>
                                    </CardHeader>
                                    <CardContent>
                                        <CardDescription className="text-sm leading-relaxed">{server.description}</CardDescription>
                                    </CardContent>
                                </Card>
                            ))}

                            <Else>
                                <small className="text-sm leading-none font-medium">No servers yet</small>
                            </Else>
                        </If>
                    </div>
                </section>
            </main>
        </AppLayout>
    );
}
