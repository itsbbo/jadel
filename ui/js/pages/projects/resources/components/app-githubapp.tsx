import { Github } from 'lucide-react';
import { ResourceCard } from '../../components/resource-card';

const cardProps = {
    title: 'Private Repository (with GitHub App)',
    description: 'You can deploy public & private repositories through your GitHub Apps.',
    icon: Github,
    category: 'Git Based',
    isPopular: true,
};

export default function GithubApp() {
    return (
        <>
            <ResourceCard {...cardProps} />
        </>
    );
}
