import Heading from '@/components/heading';
import AppLayout from '@/layouts/app-layout';
import { BreadcrumbItem } from '@/types';
import { Head } from '@inertiajs/react';
import { Database, FileText, GitBranch, Github, Key } from 'lucide-react';
import { ResourceCard } from '../components/resource-card';

const breadcrumbs: BreadcrumbItem[] = [
    {
        title: 'Projects',
        href: '/projects',
    },
];

const applicationResources = [
    {
        title: 'Public Repository',
        description: 'You can deploy any kind of public repositories from the supported git providers.',
        icon: GitBranch,
        category: 'Git Based',
        isPopular: true,
    },
    {
        title: 'Private Repository (with GitHub App)',
        description: 'You can deploy public & private repositories through your GitHub Apps.',
        icon: Github,
        category: 'Git Based',
    },
    {
        title: 'Private Repository (with Deploy Key)',
        description: 'You can deploy private repositories with a deploy key.',
        icon: Key,
        category: 'Git Based',
    },
];

const dockerResources = [
    {
        title: 'Dockerfile',
        description: 'You can deploy a simple Dockerfile, without Git.',
        icon: FileText,
        category: 'Docker Based',
    },
    {
        title: 'Docker Compose Empty',
        description: 'You can deploy complex application easily with Docker Compose, without Git.',
        icon: FileText,
        category: 'Docker Based',
    },
    {
        title: 'Docker Image',
        description: 'You can deploy an existing Docker Image from any Registry, without Git.',
        icon: FileText,
        category: 'Docker Based',
    },
];

const databaseResources = [
    {
        title: 'PostgreSQL',
        description: 'PostgreSQL is an object-relational database known for its robustness, advanced features, and strong standards compliance.',
        icon: Database,
        category: 'Database',
        isPopular: true,
    },
    {
        title: 'MySQL',
        description: 'MySQL is an open-source relational database management system.',
        icon: Database,
        category: 'Database',
    },
    {
        title: 'MariaDB',
        description: 'MariaDB is a community-developed, commercially supported fork of the MySQL relational database management system.',
        icon: Database,
        category: 'Database',
    },
    {
        title: 'Redis',
        description: 'Redis is a source-available, in-memory storage, used as a distributed, in-memory key-value database, cache and message broker.',
        icon: Database,
        category: 'Database',
    },
];

export default function CreateResources() {
    return (
        <AppLayout breadcrumbs={breadcrumbs}>
            <Head title="New Resources" />

            <main className="space-y-16 px-4 py-6">
                <Heading title="Resources" description="Deploy resources, like Applications, Databases, Services..." />

                <section>
                    <Heading className="mb-2" title="Applications" />
                    <h3 className="text-primary mb-4 text-lg font-medium">Git Based</h3>
                    <div className="mb-8 grid grid-cols-1 gap-4 xl:grid-cols-3">
                        {applicationResources.map((resource, index) => (
                            <ResourceCard key={index} {...resource} />
                        ))}
                    </div>
                    <h3 className="text-primary mb-4 text-lg font-medium">Docker Based</h3>
                    <div className="mb-8 grid grid-cols-1 gap-4 xl:grid-cols-3">
                        {dockerResources.map((resource, index) => (
                            <ResourceCard key={index} {...resource} />
                        ))}
                    </div>
                </section>

                <section>
                    <Heading className="mb-2" title="Databases" />
                    <div className="mb-8 grid grid-cols-1 gap-4 xl:grid-cols-3">
                        {databaseResources.map((resource, index) => (
                            <ResourceCard key={index} {...resource} />
                        ))}
                    </div>
                </section>
            </main>
        </AppLayout>
    );
}
