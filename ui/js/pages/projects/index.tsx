import { Else, If } from '@/components/condition';
import Heading from '@/components/heading';
import PaginationControls from '@/components/pagination-controls';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/shadcn/card';
import AppLayout from '@/layouts/app-layout';
import { BreadcrumbItem, Pagination } from '@/types';
import { Project } from '@/types/entity';
import { Head, router } from '@inertiajs/react';
import AddProjects from './components/add-projects';

const breadcrumbs: BreadcrumbItem[] = [
    {
        title: 'Projects',
        href: '/projects',
    },
];

interface Props extends Pagination<Project> {}

export default function Index({ items, prevPageURL, nextPageURL }: Props) {
    const handleOnClickProject = (id: string) => {
        router.visit(`/project/${id}`);
    };

    return (
        <AppLayout breadcrumbs={breadcrumbs}>
            <Head title="Projects" />

            <main className="space-y-16 px-4 py-6">
                <Heading title="Projects" description="All your projects are here." />

                <section className="space-y-4">
                    <div>
                        <AddProjects />
                    </div>
                    <If condition={(items?.length ?? 0) === 0}>
                        <small className="text-sm leading-none font-medium">No projects yet</small>
                        <Else>
                            <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
                                {items?.map((project) => (
                                    <Card
                                        key={project.id}
                                        className="cursor-pointer transition-shadow hover:shadow-md"
                                        onClick={() => handleOnClickProject(project.id)}
                                    >
                                        <CardHeader>
                                            <CardTitle className="text-lg">{project.name}</CardTitle>
                                        </CardHeader>
                                        <CardContent>
                                            <CardDescription className="text-sm leading-relaxed">
                                                {project.description === '' ? 'No description provided' : project.description}
                                            </CardDescription>
                                        </CardContent>
                                    </Card>
                                ))}
                            </div>
                            <PaginationControls prevPageURL={prevPageURL} nextPageURL={nextPageURL} />
                        </Else>
                    </If>
                </section>
            </main>
        </AppLayout>
    );
}
