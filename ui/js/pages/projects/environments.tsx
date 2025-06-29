import { Else, If } from '@/components/condition';
import Heading from '@/components/heading';
import { Card, CardHeader, CardTitle } from '@/components/shadcn/card';
import AppLayout from '@/layouts/app-layout';
import { BreadcrumbItem } from '@/types';
import { Environment, Project } from '@/types/entity';
import { Head, router } from '@inertiajs/react';

interface Props {
    project: Project;
    envs: Environment[];
}

export default function Index({ project, envs }: Props) {
    const breadcrumbs: BreadcrumbItem[] = [
        {
            title: 'Projects',
            href: '/projects',
        },
        {
            title: project.name,
            href: `/projects/${project.id}`,
        },
    ];

    const handleOnClickProject = (id: string) => {
        router.visit(`/projects/${project.id}/environments/${id}`);
    };

    return (
        <AppLayout breadcrumbs={breadcrumbs}>
            <Head title="Environments" />

            <main className="space-y-16 px-4 py-6">
                <Heading title="Environments" description={`Available environments for ${project.name}`} />

                <section className="space-y-4">
                    <If condition={(envs?.length ?? 0) === 0}>
                        <small className="text-sm leading-none font-medium">No environments yet</small>
                        <Else>
                            <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
                                {envs?.map((env) => (
                                    <Card key={env.id} onClick={() => handleOnClickProject(env.id)}>
                                        <CardHeader>
                                            <CardTitle className="text-lg">{env.name}</CardTitle>
                                        </CardHeader>
                                    </Card>
                                ))}
                            </div>
                        </Else>
                    </If>
                </section>
            </main>
        </AppLayout>
    );
}
