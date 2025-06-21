import { Database } from 'lucide-react';
import { ResourceCard } from '../../components/resource-card';

const cardProps = {
    title: 'PostgreSQL',
    description: 'PostgreSQL is an object-relational database known for its robustness, advanced features, and strong standards compliance.',
    icon: Database,
    category: 'Database',
    isPopular: true,
};

export default function Postgres() {
    return <ResourceCard {...cardProps} />;
}
