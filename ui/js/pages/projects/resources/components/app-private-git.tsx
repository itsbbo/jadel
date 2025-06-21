import { Key } from 'lucide-react';
import { ResourceCard } from '../../components/resource-card';

const cardProps = {
    title: 'Private Repository (with Deploy Key)',
    description: 'You can deploy private repositories with a deploy key.',
    icon: Key,
    category: 'Git Based',
};

export default function PrivateGit() {
    return (
        <>
            <ResourceCard {...cardProps} />
        </>
    );
}
