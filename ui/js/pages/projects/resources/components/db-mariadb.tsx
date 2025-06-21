import { Database } from 'lucide-react';
import { ResourceCard } from '../../components/resource-card';

const cardProps = {
    title: 'MariaDB',
    description: 'MariaDB is a community-developed, commercially supported fork of the MySQL relational database management system.',
    icon: Database,
    category: 'Database',
};

export default function MariaDB() {
    return <ResourceCard {...cardProps} />;
}
