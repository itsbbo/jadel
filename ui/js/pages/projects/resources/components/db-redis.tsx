import { Database } from 'lucide-react';
import { ResourceCard } from '../../components/resource-card';

const cardProps = {
    title: 'Redis',
    description: 'Redis is a source-available, in-memory storage, used as a distributed, in-memory key-value database, cache and message broker.',
    icon: Database,
    category: 'Database',
};

export default function Redis() {
    return <ResourceCard {...cardProps} />;
}
