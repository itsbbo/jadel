import { Else, If } from '@/components/condition';
import Heading from '@/components/heading';
import PaginationControls from '@/components/pagination-controls';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/shadcn/card';
import AppLayout from '@/layouts/app-layout';
import { BreadcrumbItem, Pagination } from '@/types';
import { Server } from '@/types/entity';
import { Head, router } from '@inertiajs/react';
import { useMemo } from 'react';
import AddServer from './components/add-server';

const breadcrumbs: BreadcrumbItem[] = [
    {
        title: 'Server',
        href: '/servers',
    },
];

interface Props extends Pagination<Server> {}

export default function Index({ items, prevId, nextId }: Props) {
    const prevPageURL = useMemo(() => {
        if (prevId) {
            return `/servers?prevId=${prevId}`;
        }

        return '';
    }, [prevId]);

    const nextPageURL = useMemo(() => {
        if (nextId) {
            return `/servers?nextId=${nextId}`;
        }

        return '';
    }, [nextId]);

    const handleOnClickServer = (id: string) => {
        router.visit(`/servers/${id}`);
    };

    return (
        <AppLayout breadcrumbs={breadcrumbs}>
            <Head title="Servers" />

            <main className="space-y-16 px-4 py-6">
                <Heading title="Servers" description="All your servers are here." />

                <section className="space-y-4">
                    <div>
                        <AddServer />
                    </div>
                    <If condition={(items?.length ?? 0) === 0}>
                        <small className="text-sm leading-none font-medium">No servers yet</small>
                        <Else>
                            <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
                                {items?.map((project) => (
                                    <Card key={project.id} onClick={() => handleOnClickServer(project.id)}>
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
                            <PaginationControls nextPageURL={nextPageURL} prevPageURL={prevPageURL} />
                        </Else>
                    </If>
                </section>
            </main>
        </AppLayout>
    );
}
