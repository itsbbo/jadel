import { Database } from 'lucide-react';
import { ResourceCard } from '../../components/resource-card';

const cardProps = {
    title: 'MySQL',
    description: 'MySQL is an open-source relational database management system.',
    icon: Database,
    category: 'Database',
};

export default function MySQL() {
    return <ResourceCard {...cardProps} />;
}
