import { Else, If } from '@/components/condition';
import Heading from '@/components/heading';
import { Button } from '@/components/shadcn/button';
import { CardContent, CardDescription, CardHeader, CardTitle } from '@/components/shadcn/card';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/shadcn/tabs';
import AppLayout from '@/layouts/app-layout';
import { BreadcrumbItem } from '@/types';
import { Environment } from '@/types/entity';
import { Head, router } from '@inertiajs/react';
import { PlusIcon } from 'lucide-react';

interface Props {
    env: Environment;
}

export default function Index({ env }: Props) {
    const breadcrumbs: BreadcrumbItem[] = [
        {
            title: 'Projects',
            href: '/projects',
        },
        {
            title: env.project.name,
            href: `/projects/${env.project.id}/environments`,
        },
        {
            title: env.name,
            href: `/projects/${env.project.id}/environments/${env.id}`,
        },
    ];

    return (
        <AppLayout breadcrumbs={breadcrumbs}>
            <Head title="Resources" />

            <main className="space-y-8 px-4 py-6">
                <Heading title="Resources" description={`Resources for ${env.name} environment`} />

                <Button className="cursor-pointer" onClick={() => router.visit(`/projects/${env.project.id}/environments/${env.id}/create`)}>
                    <PlusIcon /> Add
                </Button>

                <div className="w-full">
                    <Tabs defaultValue="applications" className="w-full">
                        <TabsList className="grid w-1/2 grid-cols-2">
                            <TabsTrigger value="applications">Applications</TabsTrigger>
                            <TabsTrigger value="databases">Databases</TabsTrigger>
                        </TabsList>

                        <TabsContent value="applications" className="mt-6 space-y-4">
                            <CardHeader>
                                <CardTitle>Applications</CardTitle>
                                <CardDescription>Manage your applications and their configurations</CardDescription>
                            </CardHeader>
                            <If condition={(env.applications?.length ?? 0) === 0}>
                                <div className="ml-6">
                                    <p className="text-sm">No applications in this environment.</p>
                                </div>
                                <Else>
                                    <CardContent className="space-y-4">
                                        <section className="grid grid-cols-1 gap-4 md:grid-cols-2 lg:grid-cols-3">
                                            {env.applications?.map((app) => (
                                                <div className="flex flex-col rounded-lg border p-4">
                                                    <div className="mb-2 flex items-center justify-between">
                                                        <h3 className="font-semibold">{app.name}</h3>
                                                        <span className="rounded-full bg-green-100 px-2 py-1 text-xs">{app.build_tool}</span>
                                                    </div>
                                                    <p className="text-muted-foreground text-sm">{app.description ?? 'No description'}</p>
                                                </div>
                                            ))}
                                        </section>
                                    </CardContent>
                                </Else>
                            </If>
                        </TabsContent>

                        <TabsContent value="databases" className="mt-6 space-y-4">
                            <CardHeader>
                                <CardTitle>Databases</CardTitle>
                                <CardDescription>Monitor and manage your database connections</CardDescription>
                            </CardHeader>
                            <If condition={(env.databases?.length ?? 0) === 0}>
                                <div className="ml-6">
                                    <p className="text-sm">No databases in this environment.</p>
                                </div>

                                <Else>
                                    <CardContent className="space-y-4">
                                        <section className="grid grid-cols-1 gap-4 md:grid-cols-2 lg:grid-cols-3">
                                            {env.databases?.map((db) => (
                                                <div className="flex flex-col rounded-lg border p-4">
                                                    <div className="mb-2 flex items-center justify-between">
                                                        <h3 className="font-semibold">{db.name}</h3>
                                                        <span className="rounded-full bg-green-100 px-2 py-1 text-xs">{db.database_type}</span>
                                                    </div>
                                                    <p className="text-muted-foreground text-sm">{db.description ?? 'No description'}</p>
                                                </div>
                                            ))}
                                        </section>
                                    </CardContent>
                                </Else>
                            </If>
                        </TabsContent>
                    </Tabs>
                </div>
            </main>
        </AppLayout>
    );
}
