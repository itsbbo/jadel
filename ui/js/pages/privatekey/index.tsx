import { Else, If } from '@/components/condition';
import Heading from '@/components/heading';
import PaginationControls from '@/components/pagination-controls';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/shadcn/card';
import AppLayout from '@/layouts/app-layout';
import { BreadcrumbItem, Pagination } from '@/types';
import { Server } from '@/types/entity';
import { Head, router } from '@inertiajs/react';
import { useMemo } from 'react';
import AddPrivateKey from './components/add-private-key';

const breadcrumbs: BreadcrumbItem[] = [
    {
        title: 'Private Key',
        href: '/private-keys',
    },
];

interface Props extends Pagination<Server> {}

export default function Index({ items, prevId, nextId }: Props) {
    const prevPageURL = useMemo(() => {
        if (prevId) {
            return `/private-keys?prevId=${prevId}`;
        }

        return '';
    }, [prevId]);

    const nextPageURL = useMemo(() => {
        if (nextId) {
            return `/private-keys?nextId=${nextId}`;
        }

        return '';
    }, [nextId]);

    const handleOnClickPrivateKey = (id: string) => {
        router.visit(`/private-keys/${id}`);
    };

    return (
        <AppLayout breadcrumbs={breadcrumbs}>
            <Head title="Private Keys" />

            <main className="space-y-16 px-4 py-6">
                <Heading title="Private Keys" description="All your private keys are here." />

                <section className="space-y-4">
                    <div>
                        <AddPrivateKey />
                    </div>
                    <If condition={(items?.length ?? 0) === 0}>
                        <small className="text-sm leading-none font-medium">No private keys yet</small>
                        <Else>
                            <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
                                {items?.map((privateKey) => (
                                    <Card key={privateKey.id} onClick={() => handleOnClickPrivateKey(privateKey.id)}>
                                        <CardHeader>
                                            <CardTitle className="text-lg">{privateKey.name}</CardTitle>
                                        </CardHeader>
                                        <CardContent>
                                            <CardDescription className="text-sm leading-relaxed">
                                                {privateKey.description === '' ? 'No description provided' : privateKey.description}
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
