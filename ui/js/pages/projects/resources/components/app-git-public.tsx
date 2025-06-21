import { GitBranch } from 'lucide-react';
import { ResourceCard } from '../../components/resource-card';

const cardProps = {
    title: 'Public Repository',
    description: 'You can deploy any kind of public repositories from the supported git providers.',
    icon: GitBranch,
    category: 'Git Based',
    isPopular: true,
    onClick: () => {
        console.log('foosha');
    },
};

export default function GitPublic() {
    return (
        <>
            <ResourceCard {...cardProps} />
        </>
    );
}
