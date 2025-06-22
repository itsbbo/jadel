import Heading from '@/components/heading';
import AppLayout from '@/layouts/app-layout';
import { BreadcrumbItem } from '@/types';
import { Environment } from '@/types/entity';
import { Head } from '@inertiajs/react';
import GitPublic from './components/app-git-public';
import GithubApp from './components/app-githubapp';
import PrivateGit from './components/app-private-git';
import MariaDB from './components/db-mariadb';
import MySQL from './components/db-mysql';
import Postgres from './components/db-postgres';
import Redis from './components/db-redis';
import DockerCompose from './components/docker-compose';
import Dockerfile from './components/docker-dockerfile';
import DockerImage from './components/docker-image';

const apps = [GitPublic, GithubApp, PrivateGit];

const dockers = [Dockerfile, DockerCompose, DockerImage];

const databases = [Postgres, MySQL, MariaDB, Redis];

interface Props {
    env: Environment;
}

export default function CreateResources({ env }: Props) {
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
        {
            title: 'Create',
            href: `/projects/${env.project.id}/environments/${env.id}/create`,
        },
    ];

    return (
        <AppLayout breadcrumbs={breadcrumbs}>
            <Head title="New Resources" />

            <main className="space-y-16 px-4 py-6">
                <Heading title="Resources" description="Deploy resources, like Applications, Databases, Services..." />

                <section>
                    <Heading className="mb-2" title="Applications" />
                    <h3 className="text-primary mb-4 text-lg font-medium">Git Based</h3>
                    <div className="mb-8 grid grid-cols-1 gap-4 xl:grid-cols-3">
                        {apps.map((App, index) => (
                            <App key={index} env={env} />
                        ))}
                    </div>
                    <h3 className="text-primary mb-4 text-lg font-medium">Docker Based</h3>
                    <div className="mb-8 grid grid-cols-1 gap-4 xl:grid-cols-3">
                        {dockers.map((Docker, index) => (
                            <Docker key={index} />
                        ))}
                    </div>
                </section>

                <section>
                    <Heading className="mb-2" title="Databases" />
                    <div className="mb-8 grid grid-cols-1 gap-4 xl:grid-cols-3">
                        {databases.map((Database, index) => (
                            <Database key={index} />
                        ))}
                    </div>
                </section>
            </main>
        </AppLayout>
    );
}
